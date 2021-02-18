package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func TestAccHarborWebhookBasic(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	webhookName := "terraform-" + acctest.RandString(10)
	resourceName := "harbor_webhook.webhook"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_webhook"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborWebhookBasic(
					projectName,
					webhookName,
					[]string{"DELETE_ARTIFACT"},
					"http",
					"http://domain.example/webhook",
				),
				Check: testCheckResourceExists(resourceName),
			},
		},
	})
}

func TestAccHarborWebhookFull(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	webhookName := "terraform-" + acctest.RandString(10)
	resourceName := "harbor_webhook.webhook"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_webhook"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborWebhookFull(
					projectName,
					webhookName,
					[]string{"DELETE_ARTIFACT"},
					"http",
					"http://domain.example/webhook",
					"A Test Webhook",
					true,
					"Authorization: Basic AAAAAAAAAAA",
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", webhookName),
					resource.TestCheckResourceAttr(resourceName, "description", "A Test Webhook"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAccHarborWebhookUpdate(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	webhookName := "terraform-" + acctest.RandString(10)
	resourceName := "harbor_webhook.webhook"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_webhook"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborWebhookBasic(
					projectName,
					webhookName,
					[]string{"DELETE_ARTIFACT"},
					"http",
					"http://domain.example/webhook",
				),
				Check: testCheckResourceExists(resourceName),
			},
			{
				Config: testCreateHarborWebhookBasic(
					projectName,
					webhookName+"2",
					[]string{"PULL_ARTIFACT"},
					"slack",
					"http://domain.example/webhook/test",
				),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", webhookName+"2"),
				),
			},
		},
	})
}

func TestAccHarborWebhookCreateAfterManualDestroy(t *testing.T) {
	var webhookID string

	projectName := "terraform-" + acctest.RandString(10)
	webhookName := "terraform-" + acctest.RandString(10)
	resourceName := "harbor_webhook.webhook"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_webhook"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborWebhookBasic(
					projectName,
					webhookName,
					[]string{"DELETE_ARTIFACT"},
					"http",
					"http://domain.example/webhook",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckResourceExists(resourceName),
					testCheckGetResourceID(resourceName, &webhookID),
				),
			},
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*harbor.Client)

					err := client.DeleteWebhook(webhookID)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testCreateHarborWebhookBasic(
					projectName,
					webhookName,
					[]string{"DELETE_ARTIFACT"},
					"http",
					"http://domain.example/webhook",
				),
				Check: testCheckResourceExists(resourceName),
			},
		},
	})
}

func testCreateHarborWebhookBasic(projectName string, webhookName string, eventTypes []string, targetType string, targetAddress string) string {
	return fmt.Sprintf(`
resource "harbor_project" "project" {
	name     = "%s"
}

resource "harbor_webhook" "webhook" {
	name = "%s"
	project_id = harbor_project.project.id
	event_types = [%s]
	target {
		type = "%s"
		address = "%s"
	}
}
	`, projectName, webhookName, `"`+strings.Join(eventTypes, `","`)+`"`, targetType, targetAddress)
}

func testCreateHarborWebhookFull(
	projectName string,
	webhookName string,
	eventTypes []string,
	targetType string,
	targetAddress string,
	description string,
	enabled bool,
	authHeader string,
	skipVerify bool,
) string {
	return fmt.Sprintf(`
resource "harbor_project" "project" {
	name     = "%s"
}

resource "harbor_webhook" "webhook" {
	name = "%s"
	project_id = harbor_project.project.id
	event_types = [%s]
	description = "%s"
	enabled = %t
	target {
		type = "%s"
		address = "%s"
		auth_header = "%s"
		skip_cert_verify = %t
	}

}
	`, projectName, webhookName, `"`+strings.Join(eventTypes, `","`)+`"`, description, enabled, targetType, targetAddress, authHeader, skipVerify)
}
