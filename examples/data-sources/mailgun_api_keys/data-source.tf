# List all API keys
data "mailgun_api_keys" "all" {}

output "api_key_descriptions" {
  value = [for k in data.mailgun_api_keys.all.items : k.description]
}

# List API keys with filtering
data "mailgun_api_keys" "sending_keys" {
  role = "sending"
}
