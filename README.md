# cpp-module-terraform-azurerm-key-vault

Terraform module for provisioning an [Azure Key Vault](https://learn.microsoft.com/en-us/azure/key-vault/general/overview).

## Example

```hcl
module "key_vault" {
  source = "git@github.com:hmcts/cpp-module-terraform-azurerm-key-vault.git?ref=v1.0.0"

  # ... see variables.tf for the full list of inputs
}
```

## RBAC authorisation

Azure is planning to eventually phase out Access Policies and recommends using to Role Based Access Control as it is a more robust and secure authorisation mechanism.
This module and its current consumers have been fully migrated to use RBAC authorisation instead of the Access Policies.
While Access Policies are still supported by Microsoft, the ability to pass them to this module has been removed and the RBAC is enabled by default.
When deploying new Key Vault using this module you should use RBAC authorisation with an appropriate role assignment and principal id.

## Default role assignments

By default this module will create:
- `Key Vault Administrator` role assignment on the Service Principal that belongs to the ADO pipeline running it
- `Key Vault Secrets User`, `Key Vault Certificate User` and `Key Vault Crypto User` role assignments on every principal id passed to this module within the `rbac_policy` variable.

When passing `rbac_policy` entries where `principal_id` is not yet known because the dependent resource does not exist yet, set the `skip_default_roles = true`.
This is because one `rbac_policy` can contain several roles for the same principal, so to prevent TF from creating default roles for each time the same principal_id occurs
they are deduplicated - this requires knowing the principal_id value at plan time.

## Versioning

This module uses [Semantic Versioning](https://semver.org/#summary) (`MAJOR.MINOR.PATCH`),
published as Git tags of the form `vMAJOR.MINOR.PATCH` (for example `v1.0.0`).

| Segment | Increments when | Action for consumers |
|---------|-----------------|----------------------|
| MAJOR   | A backwards-incompatible change is made (renamed/removed inputs or outputs, or behaviour that forces resource replacement). | Review before upgrading; may need state changes. |
| MINOR   | New backwards-compatible functionality is added (e.g. a new optional input). | Safe to adopt. |
| PATCH   | A backwards-compatible bug fix or maintenance change is made. | Safe to adopt. |

Always pin the `?ref=` in your `source` to a released tag (e.g. `?ref=v1.0.0`)
rather than `main`, so that upgrades are deliberate.

### How releases are produced

Release automation lives in
[`.github/workflows/release-drafter.yaml`](.github/workflows/release-drafter.yaml).
Every merge to `main` keeps a **draft** GitHub Release up to date (drafted by
`release-drafter/release-drafter@v6`). A maintainer **publishes** the draft, which
creates the immutable `vX.Y.Z` Git tag that consumers pin to, and
[`CHANGELOG.md`](./CHANGELOG.md) is updated automatically via the shared
`update-changelog` action from
[`hmcts/cnp-githubactions-library`](https://github.com/hmcts/cnp-githubactions-library).

By default each release is a **patch** bump. Label-driven `MAJOR`/`MINOR` bumps
(a `.github/release-drafter.yml` config plus a standard label set) will be added
separately, as part of the HMCTS module label guidance.

### Release labels and PR titles

PR titles and labels are used by the release drafter to decide how to bump the version and are verified by the [pr-label-check workflow](./github/workflows/pr-label-check.yaml).
More detailed documentation on [how this works is available here](https://github.com/hmcts/cnp-githubactions-library/blob/main/.github/workflows/label-check.md).

Quick reference (double check above for changes/updates):

| PR title prefix | Label applied | Version bump |
| --- | --- | --- |
| `feat:` / `feat(scope):` | `enhancement` | minor `1.0.0 -> 1.1.0` |
| `fix:` / `fix(scope):` | `bug` | patch `1.0.0 -> 1.0.1` |
| `chore:` / `chore(scope):` | `chore` | patch `1.0.0 -> 1.0.1` |
| `docs:` / `docs(scope):` | `documentation` | patch `1.0.0 -> 1.0.1` |
| `deps:` / `build:` (and scoped variants) | `dependencies` / `build` | patch `1.0.0 -> 1.0.1` |
| `<anything>` | `breaking-change` * | major `1.0.0 -> 2.0.0` |

  ⚠️ * Breaking changes are not auto-detected from the PR title. Add the breaking-change label manually to any PR that introduces a breaking change.

<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
| ---- | ------- |
| <a name="provider_azurerm"></a> [azurerm](#provider\_azurerm) | 4.79.0 |
| <a name="provider_random"></a> [random](#provider\_random) | 3.9.0 |

## Resources

| Name | Type |
| ---- | ---- |
| [azurerm_key_vault.key-vault](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault) | resource |
| [azurerm_key_vault_secret.keys](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_secret) | resource |
| [azurerm_private_endpoint.endpoint-vault](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/private_endpoint) | resource |
| [azurerm_private_endpoint.external_endpoint_vault](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/private_endpoint) | resource |
| [azurerm_role_assignment.keyvault_ado_key_vault_admin_role_assignment](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/role_assignment) | resource |
| [azurerm_role_assignment.keyvault_group_role_assignment](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/role_assignment) | resource |
| [azurerm_role_assignment.keyvault_rbac_default_role_assignment](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/role_assignment) | resource |
| [random_password.passwd](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [azurerm_client_config.current](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/data-sources/client_config) | data source |
| [azurerm_private_dns_zone.dns_external](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/data-sources/private_dns_zone) | data source |
| [azurerm_private_dns_zone.dns_kv](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/data-sources/private_dns_zone) | data source |
| [azurerm_subnet.external_subnet](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/data-sources/subnet) | data source |

## Inputs

| Name | Description | Type | Default | Required |
| ---- | ----------- | ---- | ------- | :------: |
| <a name="input_application"></a> [application](#input\_application) | Application to which the s3 bucket relates | `string` | `""` | no |
| <a name="input_attribute"></a> [attribute](#input\_attribute) | An attribute of the s3 bucket that makes it unique | `string` | `""` | no |
| <a name="input_azure_ad_service_principal_names"></a> [azure\_ad\_service\_principal\_names](#input\_azure\_ad\_service\_principal\_names) | Name of theDWP PRJ number (obtained from the project portfolio in TechNow) | `string` | `"ado_nonlive_service_principal_lab"` | no |
| <a name="input_certificate_permissions"></a> [certificate\_permissions](#input\_certificate\_permissions) | n/a | `list(string)` | `[]` | no |
| <a name="input_costcode"></a> [costcode](#input\_costcode) | Name of theDWP PRJ number (obtained from the project portfolio in TechNow) | `string` | `""` | no |
| <a name="input_enable_data_lookup"></a> [enable\_data\_lookup](#input\_enable\_data\_lookup) | n/a | `bool` | `false` | no |
| <a name="input_enable_rbac_authorization"></a> [enable\_rbac\_authorization](#input\_enable\_rbac\_authorization) | Specify whether Azure Key Vault uses Role Based Access Control for authorization | `bool` | `true` | no |
| <a name="input_enabled_for_deployment"></a> [enabled\_for\_deployment](#input\_enabled\_for\_deployment) | n/a | `bool` | `true` | no |
| <a name="input_enabled_for_disk_encryption"></a> [enabled\_for\_disk\_encryption](#input\_enabled\_for\_disk\_encryption) | n/a | `bool` | `true` | no |
| <a name="input_enabled_for_template_deployment"></a> [enabled\_for\_template\_deployment](#input\_enabled\_for\_template\_deployment) | n/a | `bool` | `true` | no |
| <a name="input_environment"></a> [environment](#input\_environment) | Environment into which resource is deployed | `string` | `""` | no |
| <a name="input_external_private_endpoint_map"></a> [external\_private\_endpoint\_map](#input\_external\_private\_endpoint\_map) | Map of external private endpoints to VNet details | <pre>map(object({<br/>    vnet_resource_group_name        = string<br/>    private_dns_resource_group_name = string<br/>    subnet_name                     = string<br/>    external_subnet_id              = optional(string)<br/>  }))</pre> | `{}` | no |
| <a name="input_key_permissions"></a> [key\_permissions](#input\_key\_permissions) | n/a | `list(string)` | `[]` | no |
| <a name="input_location"></a> [location](#input\_location) | n/a | `string` | `"uksouth"` | no |
| <a name="input_name"></a> [name](#input\_name) | n/a | `string` | `""` | no |
| <a name="input_namespace"></a> [namespace](#input\_namespace) | Namespace, which could be an organization name or abbreviation, e.g. 'eg' or 'cp' | `string` | `""` | no |
| <a name="input_network_acls"></a> [network\_acls](#input\_network\_acls) | Network rules to apply to key vault. | <pre>object({<br/>    bypass                     = string<br/>    default_action             = string<br/>    ip_rules                   = list(string)<br/>    virtual_network_subnet_ids = list(string)<br/>  })</pre> | `null` | no |
| <a name="input_object_id"></a> [object\_id](#input\_object\_id) | n/a | `string` | `""` | no |
| <a name="input_owner"></a> [owner](#input\_owner) | Name of the project or sqaud within the PDU which manages the resource. May be a persons name or email also | `string` | `""` | no |
| <a name="input_private_dns_resource_group_name"></a> [private\_dns\_resource\_group\_name](#input\_private\_dns\_resource\_group\_name) | Resource group for private dns | `string` | `null` | no |
| <a name="input_public_network_access_enabled"></a> [public\_network\_access\_enabled](#input\_public\_network\_access\_enabled) | n/a | `bool` | `true` | no |
| <a name="input_purge_protection_enabled"></a> [purge\_protection\_enabled](#input\_purge\_protection\_enabled) | n/a | `bool` | `true` | no |
| <a name="input_random_password_length"></a> [random\_password\_length](#input\_random\_password\_length) | The desired length of random password created by this module | `number` | `32` | no |
| <a name="input_rbac_policy"></a> [rbac\_policy](#input\_rbac\_policy) | List of rbac policies for the Key Vault. | <pre>list(object({<br/>    principal_id         = string<br/>    role_definition_name = string<br/>  }))</pre> | `[]` | no |
| <a name="input_resource_group_name"></a> [resource\_group\_name](#input\_resource\_group\_name) | n/a | `string` | `""` | no |
| <a name="input_secret_permissions"></a> [secret\_permissions](#input\_secret\_permissions) | n/a | `list(string)` | <pre>[<br/>  "Set",<br/>  "List",<br/>  "Get",<br/>  "Delete",<br/>  "Recover",<br/>  "Purge"<br/>]</pre> | no |
| <a name="input_secrets"></a> [secrets](#input\_secrets) | A map of secrets for the Key Vault. | `map(string)` | `{}` | no |
| <a name="input_skip_default_roles"></a> [skip\_default\_roles](#input\_skip\_default\_roles) | Whether to skip the creation of default roles for the Key Vault. | `bool` | `false` | no |
| <a name="input_sku_name"></a> [sku\_name](#input\_sku\_name) | n/a | `string` | `""` | no |
| <a name="input_soft_delete_retention_days"></a> [soft\_delete\_retention\_days](#input\_soft\_delete\_retention\_days) | n/a | `number` | `1` | no |
| <a name="input_subnet_kv"></a> [subnet\_kv](#input\_subnet\_kv) | Subnet for key vault | `string` | `""` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | A map of tags to add to all resources | `map(string)` | `{}` | no |
| <a name="input_tenant_id"></a> [tenant\_id](#input\_tenant\_id) | n/a | `string` | `""` | no |

## Outputs

| Name | Description |
| ---- | ----------- |
| <a name="output_id"></a> [id](#output\_id) | The ID of the Key Vault. |
| <a name="output_name"></a> [name](#output\_name) | Name of key vault created. |
| <a name="output_vault_uri"></a> [vault\_uri](#output\_vault\_uri) | The URI of the Key Vault, used for performing operations on keys and secrets. |
<!-- END_TF_DOCS -->

## Running Terratest Locally

This repository includes a root Makefile with helper targets for local Terratest runs.

- `make az-login`: runs `az login` so Terraform can authenticate with Azure.
- `make terratest-all`: runs all Terratest tests with `go test` (depends on `az-login`).
- `make terratest-all-gotestsum`: runs all Terratest tests with `gotestsum` (depends on `az-login`).
- `make terratest-role-assignments`: runs only role-assignment tests (depends on `az-login`).
- `make terratest-role-assignments-gotestsum`: runs only role-assignment tests with `gotestsum` (depends on `az-login`).
- `make terratest-test TEST_NAME='<regex>'`: universal target to run any test subset with `go test` (depends on `az-login`).
- `make terratest-test-gotestsum TEST_NAME='<regex>'`: universal target to run any test subset with `gotestsum` (depends on `az-login`).

Optional variables:

- `TERRATEST_TIMEOUT` (default `30`) to control test timeout in minutes.

Examples:

```shell
make terratest-all
make terratest-all-gotestsum
make terratest-role-assignments
make terratest-role-assignments-gotestsum
make terratest-test TEST_NAME='^TestRoleAssignments_EmptyRbacPolicy$'
make terratest-test-gotestsum TEST_NAME='^TestRoleAssignments_'
```

Important: tests that perform real `terraform apply`/resource creation and then `destroy` may fail for users without sufficient Azure permissions in the target subscription/tenant.

## Contributing

We use pre-commit hooks for validating the terraform format and maintaining the documentation automatically.
Install it with:

```shell
$ brew install pre-commit terraform-docs
$ pre-commit install
```

If you add a new hook make sure to run it against all files:
```shell
$ pre-commit run --all-files
```
