variable "mailgun_api_key" {
  description = "Mailgun API key"
  type        = string
  sensitive   = true
}

variable "mailgun_region" {
  description = "Mailgun region (US or EU)"
  type        = string
  default     = "US"
}

variable "domain_name" {
  description = "Primary domain name for the basic example"
  type        = string
}

variable "advanced_domain_name" {
  description = "Domain name for the advanced configuration example"
  type        = string
}
