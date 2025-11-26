# Basic domain
resource "mailgun_domain" "basic" {
  name        = "mail.example.com"
  spam_action = "disabled"
  wildcard    = false
}

# Domain with advanced configuration
resource "mailgun_domain" "advanced" {
  name                          = "notifications.example.com"
  spam_action                   = "tag"
  wildcard                      = true
  use_automatic_sender_security = true
  force_dkim_authority          = true
  dkim_key_size                 = "2048"
  web_scheme                    = "https"
}
