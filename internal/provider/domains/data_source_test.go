// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package domains_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/dimoschi/terraform-provider-mailgun/internal/provider/test_helpers"
)

func TestAccDomainsDataSource(t *testing.T) {
	// Skip if MAILGUN_API_KEY is not set
	if os.Getenv("MAILGUN_API_KEY") == "" {
		t.Skip("MAILGUN_API_KEY environment variable is not set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_helpers.AccPreCheck(t) },
		ProtoV6ProviderFactories: test_helpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccDomainsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mailgun_domains.test", "total_count"),
					resource.TestCheckResourceAttrSet("data.mailgun_domains.test", "items.#"),
				),
			},
		},
	})
}

func TestAccDomainsDataSource_WithLimit(t *testing.T) {
	if os.Getenv("MAILGUN_API_KEY") == "" {
		t.Skip("MAILGUN_API_KEY environment variable is not set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_helpers.AccPreCheck(t) },
		ProtoV6ProviderFactories: test_helpers.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainsDataSourceConfigWithLimit(10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mailgun_domains.test", "limit", "10"),
					resource.TestCheckResourceAttrSet("data.mailgun_domains.test", "total_count"),
				),
			},
		},
	})
}

func testAccDomainsDataSourceConfig() string {
	return `
provider "mailgun" {
  api_key = "` + os.Getenv("MAILGUN_API_KEY") + `"
}

data "mailgun_domains" "test" {}
`
}

func testAccDomainsDataSourceConfigWithLimit(limit int) string {
	apiKey := os.Getenv("MAILGUN_API_KEY")
	return `
provider "mailgun" {
  api_key = "` + apiKey + `"
}

data "mailgun_domains" "test" {
  limit = 10
}
`
}
