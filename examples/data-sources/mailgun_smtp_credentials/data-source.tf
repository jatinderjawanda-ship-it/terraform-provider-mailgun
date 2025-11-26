# List all SMTP credentials for a domain
data "mailgun_smtp_credentials" "all" {
  domain = "mail.example.com"
}

output "credential_logins" {
  value = [for c in data.mailgun_smtp_credentials.all.items : c.login]
}
