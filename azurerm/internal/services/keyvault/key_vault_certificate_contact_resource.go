package keyvault

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceKeyVaultCertificateContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyVaultCertificateContactCreate,
		Read:   resourceKeyVaultCertificateContactRead,
		Update: resourceKeyVaultCertificateContactUpdate,
		Delete: resourceKeyVaultCertificateContactDelete,

		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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

func resourceKeyVaultCertificateContactCreate(d *schema.ResourceData, meta interface{}) error {
	clients := meta.(*clients.Client)
	ctx, cancel := timeouts.ForRead(clients.StopContext, d)
	defer cancel()

	// Locking to prevent parallel changes causing issues
	keyVaultId := d.Get("key_vault_id").(string)
	locks.ByName(keyVaultId, keyVaultResourceName)
	defer locks.UnlockByName(keyVaultId, keyVaultResourceName)

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
	name := d.Get("name").(string)
	phone := d.Get("phone").(string)

	contactList := make([]keyvault.Contact, 0)
	if contacts.ContactList != nil {
		contactList = *contacts.ContactList

		if findContactFromEmail(email, contacts) != -1 {
			return tf.ImportAsExistsError("azurerm_key_vault_certificate_contact", email)
		}
	}

	contactList = append(contactList, keyvault.Contact{
		Name:         &name,
		EmailAddress: &email,
		Phone:        &phone,
	})

	contacts = keyvault.Contacts{
		ContactList: &contactList,
	}

	if _, err := clients.KeyVault.ManagementClient.SetCertificateContacts(ctx, keyVaultUrl, contacts); err != nil {
		return err
	}

	// This is because azure doesn't have an 'id' for a keyvault certificate contacts
	// In order to compensate for this and allow importing of this resource we are artificially
	// creating an identity for a key vault certificate contacts object
	d.SetId(fmt.Sprintf("%s/contacts/%s", keyVaultId, email))

	return resourceKeyVaultCertificateContactRead(d, meta)
}

func resourceKeyVaultCertificateContactUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("email") || d.HasChange("name") || d.HasChange("phone") {

		clients := meta.(*clients.Client)
		ctx, cancel := timeouts.ForRead(clients.StopContext, d)
		defer cancel()

		// Locking to prevent parallel changes causing issues
		keyVaultId := d.Get("key_vault_id").(string)
		locks.ByName(keyVaultId, keyVaultResourceName)
		defer locks.UnlockByName(keyVaultId, keyVaultResourceName)

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
		name := d.Get("name").(string)
		phone := d.Get("phone").(string)

		if index := findContactFromEmail(email, contacts); index != -1 {
			(*contacts.ContactList)[index] = keyvault.Contact{
				EmailAddress: &email,
				Name:         &name,
				Phone:        &phone,
			}

			if _, err := clients.KeyVault.ManagementClient.SetCertificateContacts(ctx, keyVaultUrl, contacts); err != nil {
				return err
			}
		}
	}

	return resourceKeyVaultCertificateContactRead(d, meta)
}

func resourceKeyVaultCertificateContactRead(d *schema.ResourceData, meta interface{}) error {
	clients := meta.(*clients.Client)
	ctx, cancel := timeouts.ForRead(clients.StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	email := id.Path["contacts"]
	keyVaultId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s", id.SubscriptionID, id.ResourceGroup, id.Path["vaults"])

	keyVaultUrl, err := getKeyVaultBaseUrl(ctx, d, clients, keyVaultId)
	if err != nil {
		return err
	}
	if keyVaultUrl == "" {
		return nil
	}

	d.Set("key_vault_id", keyVaultId)

	contacts, err := findContactsFromKeyVaultUrl(ctx, d, clients, keyVaultUrl)
	if err != nil {
		return err
	}

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

	return nil
}

func resourceKeyVaultCertificateContactDelete(d *schema.ResourceData, meta interface{}) error {
	clients := meta.(*clients.Client)
	ctx, cancel := timeouts.ForRead(clients.StopContext, d)
	defer cancel()

	// Locking to prevent parallel changes causing issues
	keyVaultId := d.Get("key_vault_id").(string)
	locks.ByName(keyVaultId, keyVaultResourceName)
	defer locks.UnlockByName(keyVaultId, keyVaultResourceName)

	keyVaultUrl, err := getKeyVaultBaseUrl(ctx, d, clients, keyVaultId)
	if err != nil {
		return err
	}

	contacts, err := findContactsFromKeyVaultUrl(ctx, d, clients, keyVaultUrl)
	if err != nil {
		return err
	}

	email := d.Get("email").(string)

	if len(*contacts.ContactList) > 1 {
		if index := findContactFromEmail(email, contacts); index != -1 {
			*contacts.ContactList = append((*contacts.ContactList)[:index], (*contacts.ContactList)[index+1:]...)
			if _, err := clients.KeyVault.ManagementClient.SetCertificateContacts(ctx, keyVaultUrl, contacts); err != nil {
				return err
			}
		}
	} else {
		clients.KeyVault.ManagementClient.DeleteCertificateContacts(ctx, keyVaultUrl)
	}

	return nil
}

func getKeyVaultBaseUrl(ctx context.Context, d *schema.ResourceData, clients *clients.Client, inputVaultId string) (keyVaultUrl string, err error) {
	keyVaultId, err := parse.VaultID(inputVaultId)
	if err != nil {
		return "", err
	}

	ok, err := clients.KeyVault.Exists(ctx, *keyVaultId)
	if err != nil {
		return "", fmt.Errorf("Error checking if key vault %q exists: %v", *keyVaultId, err)
	}
	if !ok {
		log.Printf("[DEBUG] Key Vault %q was not found - removing from state", *keyVaultId)
		d.SetId("")
		return "", nil
	}

	keyVaultBaseUrl, err := clients.KeyVault.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return "", fmt.Errorf("Error retieving base url of key vault %q: %v", *keyVaultId, err)
	}

	keyVaultUrl = *keyVaultBaseUrl

	return
}

func findContactsFromKeyVaultUrl(ctx context.Context, d *schema.ResourceData, clients *clients.Client, keyVaultBaseUrl string) (contacts keyvault.Contacts, err error) {
	contacts, err = clients.KeyVault.ManagementClient.GetCertificateContacts(ctx, keyVaultBaseUrl)
	if err != nil {
		if utils.ResponseWasNotFound(contacts.Response) {
			d.SetId("")
			return keyvault.Contacts{}, nil
		}
		return keyvault.Contacts{}, fmt.Errorf("Error making Read request on Azure KeyVault Certificate Contacts: %+v", err)
	}

	return
}

func findContactFromEmail(email string, contacts keyvault.Contacts) (index int) {
	index = -1
	for i, c := range *contacts.ContactList {
		if email == *c.EmailAddress {
			index = i
			break
		}
	}
	return
}
