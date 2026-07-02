package test

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// Testing the secure-file-transfer Module
func TestTerraformAzureKeyVault(t *testing.T) {
	t.Parallel()

	//subscriptionID := "e6b5053b-4c38-4475-a835-a025aeb3d8c7"
	// Terraform plan.out File Path
	exampleFolder := test_structure.CopyTerraformFolderToTemp(t, "../..", "example/")
	planFilePath := filepath.Join(exampleFolder, "plan.out")

	terraformPlanOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../example/",
		Upgrade:      true,

		// Variables to pass to our Terraform code using -var options
		VarFiles: []string{"kv_terratest.tfvars"},

		//Environment variables to set when running Terraform

		// Configure a plan file path so we can introspect the plan and make assertions about it.
		PlanFilePath: planFilePath,
	})

	// Run terraform init plan and show and fail the test if there are any errors
	terraform.InitAndPlanAndShowWithStruct(t, terraformPlanOptions)

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformPlanOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformPlanOptions)

	// Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformPlanOptions, "resource_group_name")
	kv_name := terraform.Output(t, terraformPlanOptions, "name")
	kv_id := terraform.Output(t, terraformPlanOptions, "id")
 	subscriptionID := terraform.Output(t, terraformPlanOptions, "subscription_id")

	keyVault, _ := azure.GetKeyVaultE(t, resourceGroupName, kv_name, subscriptionID)

	assert.Equal(t, kv_id, *keyVault.ID)
}

type roleChange struct {
	Address            string
	ResourceName       string
	PrincipalID        string
	RoleDefinitionName string
}

type rolePlanSummary struct {
	Explicit []roleChange
	Admin    []roleChange
	Default  []roleChange
}

func fixturePath(t *testing.T, fixtureName string) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to resolve current test file path")
	}

	terratestDir := filepath.Dir(currentFile)
	return filepath.Join(terratestDir, "..", "fixtures", fixtureName)
}

func planRoleAssignments(t *testing.T, fixtureName string) rolePlanSummary {
	t.Helper()

	fixtureDir := fixturePath(t, fixtureName)
	planFilePath := filepath.Join(fixtureDir, "plan.out")
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: fixtureDir,
		PlanFilePath: planFilePath,
		NoColor:      true,
	})

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)
	summary := rolePlanSummary{}

	for _, change := range plan.ResourceChangesMap {
		if change.Type != "azurerm_role_assignment" {
			continue
		}

		assert.True(t, strings.Contains(change.Address, "module.key_vault"), "unexpected role assignment address: %s", change.Address)

		rc := roleChange{
			Address:      change.Address,
			ResourceName: change.Name,
		}

		if after, ok := change.Change.After.(map[string]interface{}); ok {
			if value, exists := after["role_definition_name"]; exists && value != nil {
				if roleValue, ok := value.(string); ok {
					rc.RoleDefinitionName = roleValue
				}
			}
			if value, exists := after["principal_id"]; exists && value != nil {
				if principalValue, ok := value.(string); ok {
					rc.PrincipalID = principalValue
				}
			}
		}

		switch change.Name {
		case "keyvault_group_role_assignment":
			summary.Explicit = append(summary.Explicit, rc)
		case "keyvault_ado_key_vault_admin_role_assignment":
			summary.Admin = append(summary.Admin, rc)
		case "keyvault_rbac_default_role_assignment":
			summary.Default = append(summary.Default, rc)
		}
	}

	return summary
}

func roleNames(changes []roleChange) []string {
	result := make([]string, 0, len(changes))
	for _, change := range changes {
		result = append(result, change.RoleDefinitionName)
	}
	return result
}

func principalIDs(changes []roleChange) []string {
	result := make([]string, 0, len(changes))
	for _, change := range changes {
		result = append(result, change.PrincipalID)
	}
	return result
}

func TestRoleAssignments_EmptyRbacPolicy(t *testing.T) {
	t.Parallel()

	summary := planRoleAssignments(t, "rbac-empty")

	assert.Len(t, summary.Explicit, 0)
	assert.Len(t, summary.Admin, 1)
	assert.Len(t, summary.Default, 0)
}

func TestRoleAssignments_RbacDefaultDedupe(t *testing.T) {
	t.Parallel()

	summary := planRoleAssignments(t, "rbac-default-dedupe")

	assert.Len(t, summary.Explicit, 1)
	assert.Equal(t, []string{"Key Vault Reader"}, roleNames(summary.Explicit))

	assert.Len(t, summary.Admin, 1)
	assert.Equal(t, []string{"Key Vault Administrator"}, roleNames(summary.Admin))

	assert.Len(t, summary.Default, 3)
	assert.Contains(t, roleNames(summary.Default), "Key Vault Secrets User")
	assert.Contains(t, roleNames(summary.Default), "Key Vault Certificate User")
	assert.Contains(t, roleNames(summary.Default), "Key Vault Crypto User")
}

func TestRoleAssignments_MultiplePrincipals(t *testing.T) {
	t.Parallel()

	summary := planRoleAssignments(t, "rbac-multiple-principals")

	assert.Len(t, summary.Explicit, 2)
	assert.Equal(t, []string{"Key Vault Reader", "Key Vault Reader"}, roleNames(summary.Explicit))
	assert.Contains(t, principalIDs(summary.Explicit), "00000000-0000-0000-0000-000000000001")
	assert.Len(t, summary.Admin, 1)
	assert.Len(t, summary.Default, 6)
}

func TestRoleAssignments_CaseInsensitiveRoleMatching(t *testing.T) {
	t.Parallel()

	summary := planRoleAssignments(t, "rbac-case-insensitive")

	assert.Len(t, summary.Explicit, 1)
	assert.Equal(t, []string{"Key Vault Reader"}, roleNames(summary.Explicit))
	assert.Len(t, summary.Admin, 1)
	assert.Equal(t, []string{"Key Vault Administrator"}, roleNames(summary.Admin))
	assert.Len(t, summary.Default, 3)
	assert.Contains(t, roleNames(summary.Default), "Key Vault Secrets User")
	assert.Contains(t, roleNames(summary.Default), "Key Vault Certificate User")
	assert.Contains(t, roleNames(summary.Default), "Key Vault Crypto User")
}

func TestRoleAssignments_SkipDefaultRoles(t *testing.T) {
	t.Parallel()

	summary := planRoleAssignments(t, "rbac-skip-default-roles")

	assert.Len(t, summary.Explicit, 3)
	assert.Contains(t, roleNames(summary.Explicit), "Key Vault Secrets User")
	assert.Contains(t, roleNames(summary.Explicit), "Key Vault Reader")
	assert.Contains(t, roleNames(summary.Explicit), "Key Vault Certificates Officer")

	assert.Len(t, summary.Admin, 1)
	assert.Equal(t, []string{"Key Vault Administrator"}, roleNames(summary.Admin))
	assert.Len(t, summary.Default, 0)
}

func TestRoleAssignments_SkipCiKvAdminRole(t *testing.T) {
	t.Parallel()

	summary := planRoleAssignments(t, "rbac-skip-ci-kv-admin-role")

	assert.Len(t, summary.Explicit, 1)
	assert.Equal(t, []string{"Key Vault Reader"}, roleNames(summary.Explicit))

	assert.Len(t, summary.Admin, 0)

	assert.Len(t, summary.Default, 3)
	assert.Contains(t, roleNames(summary.Default), "Key Vault Secrets User")
	assert.Contains(t, roleNames(summary.Default), "Key Vault Certificate User")
	assert.Contains(t, roleNames(summary.Default), "Key Vault Crypto User")
}
