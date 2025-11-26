output "basic_domain_id" {
  description = "The ID of the basic domain"
  value       = mailgun_domain.example.domain.id
}

output "basic_domain_name" {
  description = "The name of the basic domain"
  value       = mailgun_domain.example.name
}

output "basic_domain_state" {
  description = "The state of the basic domain"
  value       = mailgun_domain.example.domain.state
}

output "basic_domain_smtp_login" {
  description = "SMTP login for the basic domain"
  value       = mailgun_domain.example.domain.smtp_login
}

output "basic_domain_receiving_dns_records" {
  description = "Receiving DNS records for the basic domain"
  value       = mailgun_domain.example.receiving_dns_records
}

output "basic_domain_sending_dns_records" {
  description = "Sending DNS records for the basic domain"
  value       = mailgun_domain.example.sending_dns_records
}

output "advanced_domain_id" {
  description = "The ID of the advanced domain"
  value       = mailgun_domain.advanced.domain.id
}

output "advanced_domain_name" {
  description = "The name of the advanced domain"
  value       = mailgun_domain.advanced.name
}

# Outputs from single domain lookup data source
output "looked_up_domain_state" {
  description = "State of the looked up domain"
  value       = data.mailgun_domain.lookup.state
}

output "looked_up_domain_smtp_login" {
  description = "SMTP login for the looked up domain"
  value       = data.mailgun_domain.lookup.smtp_login
}

output "looked_up_domain_created_at" {
  description = "Creation timestamp of the looked up domain"
  value       = data.mailgun_domain.lookup.created_at
}

output "all_domains_count" {
  description = "Total count of all domains in the account"
  value       = data.mailgun_domains.all.total_count
}

output "all_domains_list" {
  description = "List of all domain names"
  value       = [for domain in data.mailgun_domains.all.items : domain.name]
}

output "filtered_domains_count" {
  description = "Count of filtered domains"
  value       = data.mailgun_domains.filtered.total_count
}

output "filtered_domains_list" {
  description = "List of filtered domain names"
  value       = [for domain in data.mailgun_domains.filtered.items : domain.name]
}
