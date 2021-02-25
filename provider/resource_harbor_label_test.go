package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func TestAccHarborLabelBasic(t *testing.T) {
	labelName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_label"),
		Steps: []resource.TestStep{
			{
				Config: testHarborLabelBasic(labelName),
				Check:  testCheckResourceExists("harbor_label.label"),
			},
		},
	})
}

func TestAccHarborLabelFull(t *testing.T) {
	labelName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_label"),
		Steps: []resource.TestStep{
			{
				Config: testHarborLabelFull(labelName, "#111111", "Test Description"),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceExists("harbor_label.label"),
					resource.TestCheckResourceAttr("harbor_label.label", "color", "#111111"),
					resource.TestCheckResourceAttr("harbor_label.label", "description", "Test Description"),
				),
			},
		},
	})
}

func TestAccHarborLabelUpdate(t *testing.T) {
	labelName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_label"),
		Steps: []resource.TestStep{
			{
				Config: testHarborLabelFull(labelName, "#111111", "Test Description"),
				Check:  testCheckResourceExists("harbor_label.label"),
			},
			{
				Config: testHarborLabelFull(labelName, "#111111", "Test Description"),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceExists("harbor_label.label"),
					resource.TestCheckResourceAttr("harbor_label.label", "color", "#111111"),
					resource.TestCheckResourceAttr("harbor_label.label", "description", "Test Description"),
				),
			},
		},
	})
}

func TestAccHarborLabelUpdateProjectID(t *testing.T) {
	labelName := "terraform-" + acctest.RandString(10)
	projectName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_label"),
		Steps: []resource.TestStep{
			{
				Config: testHarborLabelWithMultiProject(projectName+"1", projectName+"2", labelName, "harbor_project.project.id"),
				Check:  testCheckResourceExists("harbor_label.label"),
			},
			{
				Config: testHarborLabelWithMultiProject(projectName+"1", projectName+"2", labelName, "harbor_project.projectTwo.id"),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceExists("harbor_label.label"),
					resource.TestCheckResourceAttrPair("harbor_label.label", "project_id", "harbor_project.projectTwo", "id"),
				),
			},
		},
	})
}

func TestAccHarborLabelUpdateScope(t *testing.T) {
	labelName := "terraform-" + acctest.RandString(10)
	projectName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_label"),
		Steps: []resource.TestStep{
			{
				Config: testHarborLabelWithProject(projectName, labelName),
				Check:  testCheckResourceExists("harbor_label.label"),
			},
			{
				Config: testHarborLabelBasic(labelName),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceExists("harbor_label.label"),
					resource.TestCheckResourceAttr("harbor_label.label", "project_id", ""),
				),
			},
		},
	})
}

func TestAccHarborLabelCreateAfterManualDestroy(t *testing.T) {
	var labelID string

	labelName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testCheckResourceDestroy("harbor_label"),
		Steps: []resource.TestStep{
			{
				Config: testHarborLabelBasic(labelName),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceExists("harbor_label.label"),
					testCheckGetResourceID("harbor_label.label", &labelID),
				),
			},
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*harbor.Client)

					err := client.DeleteLabel(labelID)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testHarborLabelBasic(labelName),
				Check:  testCheckResourceExists("harbor_label.label"),
			},
		},
	})
}

func testHarborLabelBasic(name string) string {
	return fmt.Sprintf(`
resource "harbor_label" "label" {
	name      = "%s"
}
	`, name)
}

func testHarborLabelFull(name string, color string, description string) string {
	return fmt.Sprintf(`
resource "harbor_label" "label" {
	name        = "%s"
	color       = "%s"
	description = "%s"
}
	`, name, color, description)
}

func testHarborLabelWithProject(projectName string, name string) string {
	return fmt.Sprintf(`
resource "harbor_project" "project" {
	name = "%s"
}

resource "harbor_label" "label" {
	name        = "%s"
	project_id  = harbor_project.project.id
}
	`, projectName, name)
}

func testHarborLabelWithMultiProject(projectName string, projectNameTwo string, name string, projectReference string) string {
	return fmt.Sprintf(`
resource "harbor_project" "project" {
	name = "%s"
}

resource "harbor_project" "projectTwo" {
	name = "%s"
}

resource "harbor_label" "label" {
	name        = "%s"
	project_id  = %s
}
	`, projectName, projectNameTwo, name, projectReference)
}
