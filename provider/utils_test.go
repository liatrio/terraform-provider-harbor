package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func testCheckResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*harbor.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		resourceID := rs.Primary.ID

		_, err := client.GetResource(resourceID)
		if err != nil {
			return fmt.Errorf("error getting project with id %s: %s", resourceID, err)
		}

		return nil
	}
}

func testCheckGetResourceID(resourceName string, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID

		return nil
	}
}

func testCheckResourceDestroy(resourceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			id := rs.Primary.ID

			client := testAccProvider.Meta().(*harbor.Client)

			resource, _ := client.GetResource(id)
			if resource != nil {
				return fmt.Errorf("project with id %s still exists", id)
			}
		}

		return nil
	}
}
