# Domains Package

This package implements the Mailgun domain resource and domains data source for the Terraform Mailgun provider.

## Files

- **resource.go** - Full CRUD implementation for the `mailgun_domain` resource
  - Create: Creates a new domain with configurable spam action, wildcard, DKIM settings, etc.
  - Read: Fetches domain state from Mailgun API
  - Update: Updates domain configuration (limited by Mailgun API)
  - Delete: Removes the domain from Mailgun
  - Import: Supports `terraform import` by domain name

- **data_source.go** - Implementation for the `mailgun_domains` data source
  - Lists all domains in the Mailgun account
  - Supports filtering with limit, skip, search, and other parameters

- **schema.go** - Terraform schema definitions
  - DomainResourceSchema: Defines all input and computed attributes for the resource
  - DomainsDataSourceSchema: Defines all attributes for the data source
  - Includes nested objects for domain details and DNS records

- **model.go** - Go structs for Terraform state management
  - DomainModel: State model for domain resource
  - DomainsModel: State model for domains data source
  - Custom value types for complex nested objects (DomainValue, ItemsValue, etc.)

- **resource_test.go** - Acceptance tests for domain resource
  - Tests basic CRUD operations
  - Tests wildcard domains
  - Tests import functionality

- **data_source_test.go** - Acceptance tests for domains data source
  - Tests basic data fetching
  - Tests with limit parameter

- **helpers_test.go** - Test utilities and provider factories

## Running Tests

### Unit Tests
```bash
make test
```

### Acceptance Tests
Acceptance tests require a valid Mailgun API key:

```bash
export MAILGUN_API_KEY="your-api-key-here"
make testacc
```

**Warning**: Acceptance tests create real Mailgun domains and may incur costs.

## Usage Examples

### Resource
```hcl
resource "mailgun_domain" "example" {
  name                         = "example.com"
  spam_action                  = "disabled"
  wildcard                     = true
  force_dkim_authority         = true
  dkim_key_size                = "2048"
  use_automatic_sender_security = true
}
```

### Data Source
```hcl
data "mailgun_domains" "all" {
  limit = 100
}

output "domain_count" {
  value = data.mailgun_domains.all.total_count
}
```

## Implementation Notes

- All code is manually implemented using the `mailgun-go/v5` SDK
- DNS records (receiving and sending) are automatically populated as computed attributes
- The domain resource is mostly immutable after creation; updates may require recreation
- Import is supported using the domain name as the identifier
