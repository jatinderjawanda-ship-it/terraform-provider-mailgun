# List all domains
data "mailgun_domains" "all" {}

output "all_domain_names" {
  value = [for d in data.mailgun_domains.all.items : d.name]
}

# List domains with filtering
data "mailgun_domains" "filtered" {
  limit  = 10
  state  = "active"
  search = "example"
}
