# Lookup a specific domain by name
data "mailgun_domain" "example" {
  name = "mail.example.com"
}

output "domain_state" {
  value = data.mailgun_domain.example.state
}

output "domain_smtp_login" {
  value = data.mailgun_domain.example.smtp_login
}
