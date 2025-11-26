# Create an SMTP credential for a domain
resource "mailgun_smtp_credential" "app" {
  domain   = "mail.example.com"
  login    = "app-mailer"
  password = var.smtp_password
}

# The full SMTP login will be: app-mailer@mail.example.com
output "smtp_full_login" {
  value = mailgun_smtp_credential.app.full_login
}
