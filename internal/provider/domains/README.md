# Domains Package - Provider Structure Reference

This package implements the Mailgun domains resource and data sources following the **Cloudflare Provider Pattern** - a best practice structure adopted by top-tier Terraform providers.

## 📁 File Structure

```
domains/
├── resource.go                    # Domain resource CRUD implementation
├── resource_model.go              # Resource Terraform models
├── resource_schema.go             # Resource schema definitions
├── resource_test.go               # Resource tests (unit + acceptance)
│
├── data_source.go                 # Single domain lookup (by name)
├── data_source_model.go           # Single domain data source model
├── data_source_schema.go          # Single domain schema
│
├── list_data_source.go            # List all domains
├── list_data_source_model.go      # List data source model
├── list_data_source_schema.go     # List schema with filtering
├── data_source_test.go            # Data source tests (unit + acceptance)
│
└── README.md                      # This file
```

## 🎯 Key Principles

### 1. Separation of Concerns
Each file has a single, clear responsibility:
- **Models**: Terraform state structure
- **Schemas**: Terraform configuration definitions
- **Implementation**: Business logic (CRUD operations, API calls)

### 2. Single vs. List Data Sources
Following Cloudflare's pattern, we separate:
- **Single resource lookup**: `mailgun_domain` - Fetch ONE domain by name
- **List resources**: `mailgun_domains` - Fetch MULTIPLE domains with filtering/pagination

### 3. Consistent Naming
- **Resource**: `DomainResource`, `DomainModel`, `DomainResourceSchema()`
- **Single Data Source**: `DomainDataSource`, `DomainDataSourceModel`, `DomainDataSourceSchema()`
- **List Data Source**: `ListDataSource`, `DomainsModel`, `DomainsListDataSourceSchema()`

## ✅ Checklist for New Resources

When implementing a new resource (e.g., `webhooks`), create:

- [ ] `webhooks/resource.go` - CRUD implementation
- [ ] `webhooks/resource_model.go` - WebhookModel + nested types
- [ ] `webhooks/resource_schema.go` - Schema + NewWebhookResource()
- [ ] `webhooks/resource_test.go` - Unit + acceptance tests
- [ ] `webhooks/data_source.go` - Single webhook lookup (if needed)
- [ ] `webhooks/data_source_model.go` - WebhookDataSourceModel (if needed)
- [ ] `webhooks/data_source_schema.go` - Schema + NewWebhookDataSource() (if needed)
- [ ] `webhooks/list_data_source.go` - List webhooks (if needed)
- [ ] `webhooks/list_data_source_model.go` - WebhooksModel + ItemsValue (if needed)
- [ ] `webhooks/list_data_source_schema.go` - Schema + NewWebhooksListDataSource() (if needed)
- [ ] `webhooks/data_source_test.go` - Data source tests (if applicable)
- [ ] Register in `internal/provider/provider.go`

## 🎓 Why This Pattern?

1. **Maintainability**: Each file is focused and easy to understand
2. **Scalability**: Adding new resources follows a clear template
3. **Testability**: Unit tests can target specific components
4. **Industry Standard**: Matches Cloudflare and other top providers
5. **Clear Boundaries**: Schema, model, and logic are cleanly separated

---

**Use this structure for all future resources!**
