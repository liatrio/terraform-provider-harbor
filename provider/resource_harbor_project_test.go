package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func TestAccHarborProjectBasic(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckHarborProjectDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHarborProjectBasic(projectName, "true"),
				Check:  testAccCheckHarborProjectExists("harbor_project.project"),
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

func TestAccHarborProjectUpdate(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckHarborProjectDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHarborProjectBasic(projectName, "true"),
				Check:  testAccCheckHarborProjectExists("harbor_project.project"),
			},
			{
				Config: testHarborProjectBasic(projectName, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHarborProjectExists("harbor_project.project"),
					resource.TestCheckResourceAttr("harbor_project.project", "public", "false"),
				),
			},
		},
	})
}

func testHarborProjectBasic(projectName string, public string) string {
	return fmt.Sprintf(`
resource "harbor_project" "project" {
	name     = "%s"
	public   = "%s"
}
	`, projectName, public)
}

func testAccCheckHarborProjectExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getProjectFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func getProjectFromState(s *terraform.State, resourceName string) (*harbor.Project, error) {
	client := testAccProvider.Meta().(*harbor.Client)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID

	project, err := client.GetProject(id)
	if err != nil {
		return nil, fmt.Errorf("error getting group with id %s: %s", id, err)
	}

	return project, nil
}

func testAccCheckHarborProjectDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "harbor_project" {
				continue
			}

			id := rs.Primary.ID

			client := testAccProvider.Meta().(*harbor.Client)

			group, _ := client.GetProject(id)
			if group != nil {
				return fmt.Errorf("group with id %s still exists", id)
			}
		}

		return nil
	}
}
