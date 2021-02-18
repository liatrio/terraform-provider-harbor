package provider

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			"project_id": {
				Type:         schema.TypeString,
				Description:  "ID of the project the webhook corresponds to, in the form '/projects/${ID_NUMBER}'",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^/projects/[0-9]+$`), "validation error: project_id should be of the form '/projects/${ID_NUMBER}'"),
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "Display name of the webhook.",
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Contains a description of the webhook.",
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "When true, webhooks are enabled.",
				Optional:    true,
				Default:     true,
			},
			"event_types": {
				Type:        schema.TypeSet,
				Description: "The list of events that will cause this webhook to trigger.",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice(
						[]string{
							"DELETE_ARTIFACT",
							"PULL_ARTIFACT",
							"PUSH_ARTIFACT",
							"DELETE_CHART",
							"DOWNLOAD_CHART",
							"UPLOAD_CHART",
							"QUOTA_EXCEED",
							"QUOTA_WARNING",
							"REPLICATION",
							"SCANNING_FAILED",
							"SCANNING_COMPLETED",
							"TAG_RETENTION",
						},
						false,
					),
				},
			},
			"target": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Description:  "The type of notification to send.",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"http", "slack"}, false),
						},
						"auth_header": {
							Type:        schema.TypeString,
							Description: "The webhook auth header.",
							Optional:    true,
						},
						"skip_cert_verify": {
							Type:        schema.TypeBool,
							Description: "If true, skips tls certificate verification of webhook target.",
							Optional:    true,
							Default:     false,
						},
						"address": {
							Type:         schema.TypeString,
							Description:  "The webhook target address.",
							Required:     true,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},
					},
				},
			},
		},
	}
}

func mapDataToWebhook(d *schema.ResourceData, webhook *harbor.Webhook) error {
	targets := &[]harbor.WebhookTargetObj{}
	mapDataToWebhookTargets(d, targets)

	webhook.Name = d.Get("name").(string)
	webhook.Description = d.Get("description").(string)
	webhook.Enabled = d.Get("enabled").(bool)
	webhook.Targets = *targets

	events := d.Get("event_types").(*schema.Set).List()
	webhook.EventTypes = make([]string, len(events))
	for i, arg := range events {
		webhook.EventTypes[i] = arg.(string)
	}

	return nil
}

func mapDataToWebhookTargets(d *schema.ResourceData, targets *[]harbor.WebhookTargetObj) {
	v, ok := d.GetOk("target")
	if !ok {
		return
	}

	for _, dataTarget := range v.(*schema.Set).List() {
		dataTarget := dataTarget.(map[string]interface{})
		target := harbor.WebhookTargetObj{}

		target.Type = dataTarget["type"].(string)
		target.Address = dataTarget["address"].(string)
		target.AuthHeader = dataTarget["auth_header"].(string)
		target.SkipCertVerify = dataTarget["skip_cert_verify"].(bool)

		*targets = append(*targets, target)
	}
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

	targets := []interface{}{}
	for _, target := range webhook.Targets {
		targetData := map[string]interface{}{
			"type":             target.Type,
			"address":          target.Address,
			"auth_header":      target.AuthHeader,
			"skip_cert_verify": target.SkipCertVerify,
		}
		targets = append(targets, targetData)
	}
	err = d.Set("target", targets)
	if err != nil {
		return err
	}

	err = d.Set("event_types", webhook.EventTypes)
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
	err := mapDataToWebhook(d, webhook)
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
	err := mapDataToWebhook(d, webhook)
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
