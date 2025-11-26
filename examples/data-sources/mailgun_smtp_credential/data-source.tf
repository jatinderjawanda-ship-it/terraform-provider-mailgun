# Lookup a specific SMTP credential
data "mailgun_smtp_credential" "app" {
  domain = "mail.example.com"
  login  = "app-mailer"
}

output "smtp_full_login" {
  value = data.mailgun_smtp_credential.app.full_login
}
