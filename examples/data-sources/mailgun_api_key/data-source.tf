# Lookup a specific API key by ID
data "mailgun_api_key" "example" {
  id = "key-abc123def456"
}

output "api_key_role" {
  value = data.mailgun_api_key.example.role
}
