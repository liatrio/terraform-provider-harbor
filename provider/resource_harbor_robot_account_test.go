package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func TestAccHarborRobotAccount_Basic(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)
	resourceName := "harbor_robot_account.robot"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_robot_account"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborRobotAccount_Basic(projectName, robotName, "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckResourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "token"),
				),
			},
		},
	})
}

func TestAccHarborRobotAccount_Full(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)
	resourceName := "harbor_robot_account.robot"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_robot_account"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborRobotAccount_Full(
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

func TestAccHarborRobotAccount_Update(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)
	resourceName := "harbor_robot_account.robot"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_robot_account"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborRobotAccount_Basic(projectName, robotName, "false"),
				Check:  testCheckResourceExists(resourceName),
			},
			{
				Config: testCreateHarborRobotAccount_Basic(projectName, robotName, "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
				),
			},
		},
	})
}

func TestAccHarborRobotAccount_CreateAfterManualDestroy(t *testing.T) {
	var robotId string

	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)
	resourceName := "harbor_robot_account.robot"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_robot_account"),
		Steps: []resource.TestStep{
			{
				Config: testCreateHarborRobotAccount_Basic(projectName, robotName, "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckResourceExists(resourceName),
					testCheckGetResourceId(resourceName, &robotId),
				),
			},
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*harbor.Client)

					err := client.DeleteRobotAccount(robotId)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testCreateHarborRobotAccount_Basic(projectName, robotName, "true"),
				Check:  testCheckResourceExists(resourceName),
			},
		},
	})
}

func testCreateHarborRobotAccount_Basic(projectName string, robotName string, disabled string) string {
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

func testCreateHarborRobotAccount_Full(projectName string, robotName string, description string, disabled string) string {
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
