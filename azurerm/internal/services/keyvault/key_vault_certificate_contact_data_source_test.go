package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type KeyVaultCertificateContactDataSource struct {
}

func TestAccDataSourceKeyVaultCertificateContact_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_certificate_contact", "test")
	r := KeyVaultCertificateContactDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("email").HasValue("contact1@test.com"),
				check.That(data.ResourceName).Key("name").HasValue("contact1"),
				check.That(data.ResourceName).Key("phone").HasValue("0123456789"),
			),
		},
	})
}

func (KeyVaultCertificateContactDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_key_vault_certificate_contact" "test" {
  key_vault_id = azurerm_key_vault.test.id
  email        = azurerm_key_vault_certificate_contact.contact1.email
}
`, KeyVaultCertificateContactResource{}.basic(data))
}
