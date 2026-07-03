provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "rg-kv-tt-rbac-multi"
  location = "uksouth"
}

module "key_vault" {
  source = "../../../"

  name                            = "kvttmulti001"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  tenant_id                       = data.azurerm_client_config.current.tenant_id
  sku_name                        = "standard"
  public_network_access_enabled   = true
  enabled_for_template_deployment = true
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  purge_protection_enabled        = false
  soft_delete_retention_days      = 7
  enable_rbac_authorization       = true

  rbac_policy = [
    {
      principal_id         = data.azurerm_client_config.current.object_id
      role_definition_name = "Key Vault Reader"
    },
    {
      principal_id         = "00000000-0000-0000-0000-000000000001"
      role_definition_name = "Key Vault Reader"
    }
  ]
}
