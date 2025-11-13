variable "mailgun_api_key" {
  description = "Mailgun API key"
  type        = string
  sensitive   = true
}

variable "domain_name" {
  description = "Domain name to create in Mailgun"
  type        = string
}

variable "smtp_password" {
  description = "SMTP password for the credential"
  type        = string
  sensitive   = true
}
