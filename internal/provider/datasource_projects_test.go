package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccProjcetsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.todoist_projects.test", "id", "2328048508"),
					resource.TestCheckResourceAttr("data.todoist_projects.test", "name", "hehe"),
				),
			},
		},
	})
}

const testAccProjcetsDataSourceConfig = `
data "todoist_projects" "test" {
  id = "2328048508"
}
`
