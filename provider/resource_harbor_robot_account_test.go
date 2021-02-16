package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func TestAccHarborRobotAccountBasic(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)
	resourceName := "harbor_robot_account.robot"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_robot_account"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborRobotAccountBasic(projectName, robotName, "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckResourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "token"),
				),
			},
		},
	})
}

func TestAccHarborRobotAccountExpiresAt(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)
	resourceName := "harbor_robot_account.robot"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_robot_account"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborRobotAccountExpiration(projectName, robotName, "2035-01-01T00:00:00Z"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckResourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "token"),
					resource.TestCheckResourceAttr(resourceName, "expires_at", "2035-01-01T00:00:00Z"),
				),
			},
		},
	})
}

func TestAccHarborRobotAccountExpiresAt64BitTimestamp(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)
	resourceName := "harbor_robot_account.robot"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_robot_account"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborRobotAccountExpiration(projectName, robotName, "2040-01-01T00:00:00Z"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckResourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "token"),
					resource.TestCheckResourceAttr(resourceName, "expires_at", "2040-01-01T00:00:00Z"),
				),
			},
		},
	})
}

func TestAccHarborRobotAccountFull(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)
	resourceName := "harbor_robot_account.robot"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_robot_account"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborRobotAccountFull(
					projectName,
					robotName,
					"A Test Robot Account",
					"true",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", robotName),
					resource.TestCheckResourceAttr(resourceName, "description", "A Test Robot Account"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "token"),
				),
			},
		},
	})
}

func TestAccHarborRobotAccountUpdate(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)
	resourceName := "harbor_robot_account.robot"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_robot_account"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborRobotAccountBasic(projectName, robotName, "false"),
				Check:  testCheckResourceExists(resourceName),
			},
			{
				Config: testCreateHarborRobotAccountBasic(projectName, robotName, "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
				),
			},
		},
	})
}

func TestAccHarborRobotAccountCreateAfterManualDestroy(t *testing.T) {
	var robotID string

	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)
	resourceName := "harbor_robot_account.robot"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_robot_account"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborRobotAccountBasic(projectName, robotName, "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckResourceExists(resourceName),
					testCheckGetResourceID(resourceName, &robotID),
				),
			},
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*harbor.Client)

					err := client.DeleteRobotAccount(robotID)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testCreateHarborRobotAccountBasic(projectName, robotName, "true"),
				Check:  testCheckResourceExists(resourceName),
			},
		},
	})
}

func testCreateHarborRobotAccountBasic(projectName string, robotName string, disabled string) string {
	return fmt.Sprintf(`
resource "harbor_project" "project" {
	name     = "%s"
}

resource "harbor_robot_account" "robot" {
	name = "%s"
	project_id = harbor_project.project.id
	disabled = %s
	access {
		resource = "image"
		action = "pull"
	}
}
	`, projectName, robotName, disabled)
}
func testCreateHarborRobotAccountExpiration(projectName string, robotName string, expirationTime string) string {
	return fmt.Sprintf(`
resource "harbor_project" "project" {
	name     = "%s"
}

resource "harbor_robot_account" "robot" {
	name = "%s"
	project_id = harbor_project.project.id
	expires_at = "%s"
	access {
		resource = "image"
		action = "pull"
	}
}
	`, projectName, robotName, expirationTime)
}

func testCreateHarborRobotAccountFull(projectName string, robotName string, description string, disabled string) string {
	return fmt.Sprintf(`
resource "harbor_project" "project" {
	name     = "%s"
}

resource "harbor_robot_account" "robot" {
	project_id = harbor_project.project.id

	name = "%s"
	description = "%s"
	disabled = %s
	access {
		resource = "image"
		action = "pull"
	}
	access {
		resource = "image"
		action = "push"
	}
	access {
		resource = "helm-chart"
		action = "pull"
	}
	access {
		resource = "helm-chart"
		action = "push"
	}
}
	`, projectName, robotName, description, disabled)
}
