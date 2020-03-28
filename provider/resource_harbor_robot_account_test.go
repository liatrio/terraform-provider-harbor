package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func TestAccHarborRobotAccountBasic(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckHarborRobotAccountDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHarborRobotAccountBasic(projectName, robotName, "false"),
				Check:  testAccCheckHarborRobotAccountExists("harbor_robot_account.robot"),
			},
			//		{
			//			ResourceName:        "keycloak_group.group",
			//			ImportState:         true,
			//			ImportStateVerify:   true,
			//			ImportStateIdPrefix: realmName + "/",
			//		},
		},
	})
}

func TestAccHarborRobotAccountUpdate(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)
	robotName := "robot$terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckHarborRobotAccountDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHarborRobotAccountBasic(projectName, robotName, "false"),
				Check:  testAccCheckHarborRobotAccountExists("harbor_robot_account.robot"),
			},
			{
				Config: testHarborRobotAccountBasic(projectName, robotName, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHarborProjectExists("harbor_robot_account.robot"),
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

func testAccCheckHarborRobotAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getRobotAccountFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func getRobotAccountFromState(s *terraform.State, resourceName string) (*harbor.RobotAccount, error) {
	client := testAccProvider.Meta().(*harbor.Client)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID

	project, err := client.GetRobotAccount(id)
	if err != nil {
		return nil, fmt.Errorf("error getting group with id %s: %s", id, err)
	}

	return project, nil
}

func testAccCheckHarborRobotAccountDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "harbor_robot_account" {
				continue
			}

			id := rs.Primary.ID

			client := testAccProvider.Meta().(*harbor.Client)

			group, _ := client.GetRobotAccount(id)
			if group != nil {
				return fmt.Errorf("robot account with id %s still exists", id)
			}
		}

		return nil
	}
}
