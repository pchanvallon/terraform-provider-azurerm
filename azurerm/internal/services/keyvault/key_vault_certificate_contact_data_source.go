package keyvault

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceKeyVaultCertificateContact() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeyVaultCertificateContactRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.VaultID,
			},

			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"phone": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceKeyVaultCertificateContactRead(d *schema.ResourceData, meta interface{}) error {
	clients := meta.(*clients.Client)
	ctx, cancel := timeouts.ForRead(clients.StopContext, d)
	defer cancel()

	keyVaultId := d.Get("key_vault_id").(string)

	keyVaultUrl, err := getKeyVaultBaseUrl(ctx, d, clients, keyVaultId)
	if err != nil {
		return err
	}
	if keyVaultUrl == "" {
		return nil
	}

	contacts, err := findContactsFromKeyVaultUrl(ctx, d, clients, keyVaultUrl)
	if err != nil {
		return err
	}

	email := d.Get("email").(string)

	if len(*contacts.ContactList) > 0 {
		if index := findContactFromEmail(email, contacts); index != -1 {
			d.Set("email", (*contacts.ContactList)[index].EmailAddress)
			d.Set("name", (*contacts.ContactList)[index].Name)
			d.Set("phone", (*contacts.ContactList)[index].Phone)
		} else {
			d.SetId("")
		}
	} else {
		d.SetId("")
	}

	// This is because azure doesn't have an 'id' for a keyvault certificate contacts
	// In order to compensate for this and allow importing of this resource we are artificially
	// creating an identity for a key vault certificate contacts object
	d.SetId(fmt.Sprintf("%s/contacts/%s", keyVaultId, email))

	return nil
}
