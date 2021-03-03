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
				Config: testHarborProjectBasic(projectName),
				Check:  testCheckResourceExists("harbor_project.project"),
			},
		},
	})
}

func TestAccHarborProjectFull(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_project"),
		Steps: []resource.TestStep{
			{
				Config: testHarborProjectFull(projectName, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceExists("harbor_project.project"),
					resource.TestCheckResourceAttr("harbor_project.project", "public", "true"),
					resource.TestCheckResourceAttr("harbor_project.project", "auto_scan", "true"),
				),
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
				Config: testHarborProjectFull(projectName, false, false),
				Check:  testCheckResourceExists("harbor_project.project"),
			},
			{
				Config: testHarborProjectFull(projectName, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceExists("harbor_project.project"),
					resource.TestCheckResourceAttr("harbor_project.project", "public", "true"),
					resource.TestCheckResourceAttr("harbor_project.project", "auto_scan", "true"),
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
				Config: testHarborProjectFull(projectName, false, false),
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
				Config: testHarborProjectFull(projectName, true, true),
				Check:  testCheckResourceExists("harbor_project.project"),
			},
		},
	})
}

func TestAccHarborProjectImportAfterManualCreate(t *testing.T) {
	projectName := "terraform-" + acctest.RandString(10)

	client := testAccProvider.Meta().(*harbor.Client)
	project := &harbor.ProjectReq{}
	project.ProjectName = projectName

	location, err := client.NewProject(project)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_project"),
		Steps: []resource.TestStep{
			{
				Config:        testHarborProjectBasic(projectName),
				ResourceName:  "harbor_project.project",
				ImportStateId: location,
				ImportState:   true,
			},
		},
	})
	err = client.DeleteProject(location)
	if err != nil {
		t.Fatal(err)
	}
}

func testHarborProjectBasic(projectName string) string {
	return fmt.Sprintf(`
resource "harbor_project" "project" {
	name      = "%s"
}
	`, projectName)
}

func testHarborProjectFull(projectName string, public bool, autoScan bool) string {
	return fmt.Sprintf(`
resource "harbor_project" "project" {
	name      = "%s"
	public    = %t
	auto_scan = %t
}
	`, projectName, public, autoScan)
}
