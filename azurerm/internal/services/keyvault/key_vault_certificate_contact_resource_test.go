package keyvault_test

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KeyVaultCertificateContactResource struct {
}

func TestAccKeyVaultCertificateContact_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contact", "contact1")
	r := KeyVaultCertificateContactResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("email").HasValue("contact1@test.com"),
				check.That(data.ResourceName).Key("name").HasValue("contact1"),
				check.That(data.ResourceName).Key("phone").HasValue("0123456789"),
			),
		},
		data.ImportStep(),
	})
}

func TTestAccKeyVaultCertificateContact_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contact", "contact1")
	r := KeyVaultAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("email").HasValue("contact1@test.com"),
				check.That(data.ResourceName).Key("name").HasValue("contact1"),
				check.That(data.ResourceName).Key("phone").HasValue("0123456789"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_key_vault_certificate_contact"),
		},
	})
}

func TestAccKeyVaultAccessCertificateContact_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contact", "contact1")
	r := KeyVaultCertificateContactResource{}
	resourceName2 := "azurerm_key_vault_certificate_contact.contact2"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multiple(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("email").HasValue("contact1@test.com"),
				check.That(data.ResourceName).Key("name").HasValue("contact1"),
				check.That(data.ResourceName).Key("phone").HasValue("0123456789"),
				resource.TestCheckResourceAttr(resourceName2, "email", "contact2@test.com"),
				resource.TestCheckResourceAttr(resourceName2, "name", "contact2"),
				resource.TestCheckResourceAttr(resourceName2, "phone", "0123456789"),
			),
		},
		data.ImportStep(),
		{
			ResourceName:      resourceName2,
			ImportState:       true,
			ImportStateVerify: true,
		},
	})
}

func TestAccKeyVaultAccessCertificateContact_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contact", "contact1")
	r := KeyVaultCertificateContactResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("email").HasValue("contact1@test.com"),
				check.That(data.ResourceName).Key("name").HasValue("contact1"),
				check.That(data.ResourceName).Key("phone").HasValue("0123456789"),
			),
		},
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("email").HasValue("contact1@test.com"),
				check.That(data.ResourceName).Key("name").HasValue("contact1updated"),
				check.That(data.ResourceName).Key("phone").HasValue("0123456789"),
			),
		},
	})
}

func TestAccKeyVaultAccessCertificateContact_nonExistentVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_certificate_contact", "contact1")
	r := KeyVaultCertificateContactResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:             r.nonExistentVault(data),
			ExpectNonEmptyPlan: true,
			ExpectError:        regexp.MustCompile(`retrieving Key Vault`),
		},
	})
}

func (t KeyVaultCertificateContactResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	keyVaultId := parse.VaultId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
		Name:           id.Path["vaults"],
	}

	ok, err := clients.KeyVault.Exists(ctx, keyVaultId)
	if err != nil || !ok {
		return utils.Bool(false), fmt.Errorf("Error checking if key vault %q exists: %v", keyVaultId, err)
	}

	log.Printf("keyVaultId = %v", keyVaultId)

	keyVaultBaseUrl, err := clients.KeyVault.BaseUriForKeyVault(ctx, keyVaultId)
	if err != nil {
		return utils.Bool(false), fmt.Errorf("Error retieving base url of key vault %q: %v", keyVaultId, err)
	}

	log.Printf("keyVaultBaseUrl = %v", keyVaultBaseUrl)

	contacts, err := clients.KeyVault.ManagementClient.GetCertificateContacts(ctx, *keyVaultBaseUrl)
	if err != nil {
		return utils.Bool(false), fmt.Errorf("Error making Read request on Azure KeyVault Certificate Contacts: %+v", err)
	}

	log.Printf("contacts = %v", contacts)
	log.Printf("email = %v", id.Path["contacts"])

	index := -1
	for i, c := range *contacts.ContactList {
		if id.Path["contacts"] == *c.EmailAddress {
			index = i
			break
		}
	}

	return utils.Bool(index != -1), nil
}

func (r KeyVaultCertificateContactResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate_contact" "contact1" {
  key_vault_id = azurerm_key_vault.test.id
  name         = "contact1"
  email        = "contact1@test.com"
	phone				 = "0123456789"

  depends_on = [
    azurerm_key_vault_access_policy.sp
  ]
}
`, template)
}

func (r KeyVaultCertificateContactResource) multiple(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate_contact" "contact1" {
  key_vault_id = azurerm_key_vault.test.id
  name         = "contact1"
  email        = "contact1@test.com"
	phone				 = "0123456789"

  depends_on = [
    azurerm_key_vault_access_policy.sp
  ]
}

resource "azurerm_key_vault_certificate_contact" "contact2" {
  key_vault_id = azurerm_key_vault.test.id
  name         = "contact2"
  email        = "contact2@test.com"
	phone				 = "0123456789"

  depends_on = [
    azurerm_key_vault_access_policy.sp
  ]
}
`, template)
}

func (r KeyVaultCertificateContactResource) nonExistentVault(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate_contact" "contact1" {
  # Must appear to be URL, but not actually exist - appending a string works
  key_vault_id = "${azurerm_key_vault.test.id}NOPE"
  name         = "contact1"
  email        = "contact1@test.com"
	phone				 = "0123456789"

  depends_on = [
    azurerm_key_vault_access_policy.sp
  ]
}
`, template)
}

func (r KeyVaultCertificateContactResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_certificate_contact" "contact1" {
  key_vault_id = azurerm_key_vault.test.id
  name         = "contact1updated"
  email        = "contact1@test.com"
	phone				 = "0123456789"

  depends_on = [
    azurerm_key_vault_access_policy.sp
  ]
}
`, template)
}

func (KeyVaultCertificateContactResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_key_vault_access_policy" "sp" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  certificate_permissions = [
    "ManageContacts"
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
