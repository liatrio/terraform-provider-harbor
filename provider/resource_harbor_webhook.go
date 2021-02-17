package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func resourceWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceWebhookCreate,
		Read:   resourceWebhookRead,
		Update: resourceWebhookUpdate,
		Delete: resourceWebhookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Display name of the webhook.",
				Required:    true,
				ForceNew:    false,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Contains a description of the webhook.",
				Optional:    true,
				Default:     false,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "When true, webhooks are enabled.",
				Optional:    true,
				Default:     true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Contains a description of the notify type.",
				Optional:    true,
				Default:     false,
			},
			"auth_header": {
				Type:        schema.TypeString,
				Description: "Contains a description of the authentication header.",
				Optional:    true,
				Default:     false,
			},
			"skip_cert_verify": {
				Type:        schema.TypeBool,
				Description: "Determines whether or not to skip the cert verification.",
				Required:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "Contains a description of the webhook target address.",
				Required:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "Contains the project id associated with the webhook.",
				Required:    true,
			},
		},
	}
}

func mapDataToWebhookReq(d *schema.ResourceData, webhook *harbor.Webhook) error {
	webhook.Name = d.Get("name").(string)

	webhook.Description = d.Get("description").(string)

	webhook.Enabled = d.Get("enabled").(bool)

	webhook.Metadata = harbor.WebhookTargetObj{
		Type:           d.Get("type").(string),
		AuthHeader:     d.Get("auth_header").(string),
		SkipCertVerify: d.Get("skip_cert_verify").(bool),
		Address:        d.Get("address").(string),
	}
	return nil
}

func mapWebhookToData(d *schema.ResourceData, webhook *harbor.Webhook) error {
	err := d.Set("name", webhook.Name)
	if err != nil {
		return err
	}

	err = d.Set("description", webhook.Description)
	if err != nil {
		return err
	}

	err = d.Set("enabled", webhook.Enabled)
	if err != nil {
		return err
	}

	Type := webhook.Metadata.Type
	if err != nil {
		return err
	}
	err = d.Set("type", Type)
	if err != nil {
		return err
	}

	authHeader := webhook.Metadata.AuthHeader
	err = d.Set("auth_header", authHeader)
	if err != nil {
		return err
	}

	skipCertVerify := webhook.Metadata.SkipCertVerify
	err = d.Set("skip_cert_verify", skipCertVerify)
	if err != nil {
		return err
	}

	address := webhook.Metadata.Address
	err = d.Set("address", address)
	if err != nil {
		return err
	}
	return nil
}

func resourceWebhookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)
	webhookID := d.Id()
	webhook, err := client.GetWebhook(webhookID)
	if err != nil {
		return handleNotFoundError(err, d)
	}

	return mapWebhookToData(d, webhook)
}

func resourceWebhookCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	webhook := &harbor.Webhook{}
	err := mapDataToWebhookReq(d, webhook)
	if err != nil {
		return err
	}

	location, err := client.NewWebhook(d.Get("project_id").(string), webhook)
	if err != nil {
		return err
	}

	d.SetId(location)
	return resourceWebhookRead(d, meta)
}

func resourceWebhookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	webhook := &harbor.Webhook{}
	err := mapDataToWebhookReq(d, webhook)
	if err != nil {
		return err
	}

	err = client.UpdateWebhook(d.Id(), webhook)
	if err != nil {
		return err
	}

	return resourceWebhookRead(d, meta)
}

func resourceWebhookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	err := client.DeleteWebhook(d.Id())
	if err != nil {
		return handleNotFoundError(err, d)
	}

	d.SetId("")
	return nil
}
