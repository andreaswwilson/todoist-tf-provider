package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccProjectsResourceConfig("myprojectname"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("todoist_projects.test", "name", "myprojectname"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "todoist_projects.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProjectsResourceConfig("myprojectnametwo"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("todoist_projects.test", "name", "myprojectnametwo"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProjectsResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "todoist_projects" "test" {
  name = %[1]q
}
`, configurableAttribute)
}
