package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccHarborRobotAccountBasic(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_robot_account"),
		Steps: []resource.TestStep{
			{
				Config: testHarborRobotAccountBasic(projectName, robotName, "false"),
				Check:  testCheckResourceExists("harbor_robot_account.robot"),
			},
		},
	})
}

func TestAccHarborRobotAccountUpdate(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_robot_account"),
		Steps: []resource.TestStep{
			{
				Config: testHarborRobotAccountBasic(projectName, robotName, "false"),
				Check:  testCheckResourceExists("harbor_robot_account.robot"),
			},
			{
				Config: testHarborRobotAccountBasic(projectName, robotName, "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceExists("harbor_robot_account.robot"),
					resource.TestCheckResourceAttr("harbor_robot_account.robot", "disabled", "true"),
				),
			},
		},
	})
}

func testHarborRobotAccountBasic(projectName string, robotName string, disabled string) string {
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
