---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_key_vault_certificate_contact"
description: |-
  Gets information about an existing Key Vault Certificate Contact.
---

# Data Source: azurerm_key_vault_certificate_contact

Use this data source to access information about an existing Key Vault Certificate Contact.

## Example Usage

```hcl
data "azurerm_key_vault" "example" {
  name                = "mykeyvault"
  resource_group_name = "some-resource-group"
}

data "azurerm_key_vault_certificate_contact" "example" {
  key_vault_id = data.azurerm_key_vault.example.id
  email         = "example@contact.com"
}

output "id" {
  value = data.azurerm_key_vault_certificate_contact.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `key_vault_id` - (Required) The ID of the Key Vault in which to locate the Certificate Contact.

* `email` - (Required) The email of the Key Vault Certificate Contact.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Key Vault Certificate Contact.

* `key_vault_id` - The Key Vault ID of the Key Vault Certificate Contact.

* `email` - The email of the Key Vault Certificate Contact.

* `name` - The name of the Key Vault Certificate Contact.

* `phone` - The phone of the Key Vault Certificate Contact.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Certificate Contact.
