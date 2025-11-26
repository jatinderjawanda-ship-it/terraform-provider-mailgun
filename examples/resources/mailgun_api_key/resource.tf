# Create a domain-specific sending key
resource "mailgun_api_key" "sending" {
  role        = "sending"
  kind        = "domain"
  domain_name = "mail.example.com"
  description = "Application sending key"
}

# Create an admin key (use with caution)
resource "mailgun_api_key" "admin" {
  role        = "admin"
  description = "Admin API key for management"
}

# The secret is only available after creation
output "api_key_secret" {
  value     = mailgun_api_key.sending.secret
  sensitive = true
}
