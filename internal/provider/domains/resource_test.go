// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package domains_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/dimoschi/terraform-provider-mailgun/internal/provider/test_helpers"
)

func TestAccDomainResource(t *testing.T) {
	// Skip if MAILGUN_API_KEY is not set
	if os.Getenv("MAILGUN_API_KEY") == "" {
		t.Skip("MAILGUN_API_KEY environment variable is not set")
	}

	domainName := test_helpers.RandomDomainName()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_helpers.AccPreCheck(t) },
		ProtoV6ProviderFactories: test_helpers.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDomainResourceDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDomainResourceConfig(domainName, "disabled", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mailgun_domain.test", "name", domainName),
					resource.TestCheckResourceAttr("mailgun_domain.test", "spam_action", "disabled"),
					resource.TestCheckResourceAttr("mailgun_domain.test", "wildcard", "false"),
					resource.TestCheckResourceAttrSet("mailgun_domain.test", "domain.id"),
					resource.TestCheckResourceAttrSet("mailgun_domain.test", "domain.name"),
					resource.TestCheckResourceAttrSet("mailgun_domain.test", "domain.state"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "mailgun_domain.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"smtp_password", "dkim_key_size", "force_dkim_authority", "web_scheme", "web_prefix"},
			},
			// Update and Read testing
			{
				Config: testAccDomainResourceConfig(domainName, "tag", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mailgun_domain.test", "name", domainName),
					resource.TestCheckResourceAttr("mailgun_domain.test", "spam_action", "tag"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDomainResource_Wildcard(t *testing.T) {
	if os.Getenv("MAILGUN_API_KEY") == "" {
		t.Skip("MAILGUN_API_KEY environment variable is not set")
	}

	domainName := test_helpers.RandomDomainName()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_helpers.AccPreCheck(t) },
		ProtoV6ProviderFactories: test_helpers.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDomainResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainResourceConfig(domainName, "disabled", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mailgun_domain.test", "name", domainName),
					resource.TestCheckResourceAttr("mailgun_domain.test", "wildcard", "true"),
				),
			},
		},
	})
}

func testAccDomainResourceConfig(domainName, spamAction string, wildcard bool) string {
	return fmt.Sprintf(`
provider "mailgun" {
  api_key = "%s"
}

resource "mailgun_domain" "test" {
  name        = "%s"
  spam_action = "%s"
  wildcard    = %t
}
`, os.Getenv("MAILGUN_API_KEY"), domainName, spamAction, wildcard)
}

func testAccCheckDomainResourceDestroy(s *terraform.State) error {
	// Domain deletion is handled by the provider
	// This is a placeholder for more complex destroy checks if needed
	return nil
}
