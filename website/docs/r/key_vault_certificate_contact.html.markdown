---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_certificate_contact"
description: |-
  Manages a Key Vault Certificate Contact.
---

# azurerm_key_vault_certificate_contact

Manages a Key Vault Certificate Contact.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                = "examplekeyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "standard"
  tenant_id           = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_key_vault_certificate_contact" "example" {
  key_vault_id  = data.azurerm_key_vault.example.id
  email         = "example@contact.com"
  name          = "example-contact"
}
```

## Arguments Reference

The following arguments are supported:

* `key_vault_id` - (Required) The ID of the Key Vault in which to create the Certificate Contact.

* `email` - (Required) The email address of the Key Vault Certificate Contact. Changing this forces a new Key Vault Certificate Contact to be created.

* `name` - (Optional) The name of the Key Vault Certificate Contact.

* `phone` - (Optional) The phone number of the Key Vault Certificate Contact.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Key Vault Certificate Contact.

* `key_vault_id` - The Key Vault ID of the Key Vault Certificate Contact.

* `email` - The email of the Key Vault Certificate Contact.

* `name` - The name of the Key Vault Certificate Contact.

* `phone` - The phone of the Key Vault Certificate Contact.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Key Vault Certificate Contact.
* `update` - (Defaults to 5 minutes) Used when updating the Key Vault Certificate Contact.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Certificate Contact.
* `delete` - (Defaults to 5 minutes) Used when deleting the Key Vault Certificate Contact.

## Import

Key Vault Certificate Contact can be imported using the using the Resource ID of the Key Vault, plus the email of the Certificate Contact:

```shell
terraform import azurerm_key_vault_certificate_contact.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.KeyVault/vaults/example/contacts/example@contact.com
```
