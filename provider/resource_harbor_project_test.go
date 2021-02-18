package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func TestAccHarborProjectBasic(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_project"),
		Steps: []resource.TestStep{
			{
				Config: testHarborProjectBasic(projectName, false, false),
				Check:  testCheckResourceExists("harbor_project.project"),
			},
		},
	})
}

func TestAccHarborProjectUpdate(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_project"),
		Steps: []resource.TestStep{
			{
				Config: testHarborProjectBasic(projectName, false, false),
				Check:  testCheckResourceExists("harbor_project.project"),
			},
			{
				Config: testHarborProjectBasic(projectName, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceExists("harbor_project.project"),
					resource.TestCheckResourceAttr("harbor_project.project", "public", "true"),
				),
			},
		},
	})
}

func TestAccHarborProjectCreateAfterManualDestroy(t *testing.T) {
	var projectID string

	projectName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_project"),
		Steps: []resource.TestStep{
			{
				Config: testHarborProjectBasic(projectName, false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceExists("harbor_project.project"),
					testCheckGetResourceID("harbor_project.project", &projectID),
				),
			},
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*harbor.Client)

					err := client.DeleteProject(projectID)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testHarborProjectBasic(projectName, true, true),
				Check:  testCheckResourceExists("harbor_project.project"),
			},
		},
	})
}

func testHarborProjectBasic(projectName string, public bool, autoScan bool) string {
	return fmt.Sprintf(`
resource "harbor_project" "project" {
	name      = "%s"
	public    = %t
	auto_scan = %t
}
	`, projectName, public, autoScan)
}
