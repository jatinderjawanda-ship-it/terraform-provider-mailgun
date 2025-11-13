# Complete Mailgun Setup Example

This example demonstrates how to provision a complete Mailgun configuration from scratch using Terraform.

## What This Creates

- **Domain**: A new Mailgun domain
- **Webhooks**: Event notifications for delivered and bounced emails
- **SMTP Credentials**: Authentication for sending emails via SMTP
- **Route**: Inbound email routing rule for support emails
- **API Key**: Programmatic access key for CI/CD pipelines

## Prerequisites

- Mailgun account with API access
- Valid API key
- A domain you want to configure (can be a subdomain for testing)

## Usage

1. **Install the provider locally** (from the repository root):
   ```bash
   make install
   ```

2. **Copy the example variables file**:
   ```bash
   cd examples/complete-setup
   cp terraform.tfvars.example terraform.tfvars
   ```

3. **Edit `terraform.tfvars`** with your values:
   ```hcl
   mailgun_api_key = "key-your-actual-api-key"
   domain_name     = "test.yourdomain.com"
   smtp_password   = "your-secure-password"
   ```

4. **Initialize Terraform**:
   ```bash
   terraform init
   ```

5. **Review the plan**:
   ```bash
   terraform plan
   ```

6. **Apply the configuration**:
   ```bash
   terraform apply
   ```

7. **View the outputs**:
   ```bash
   terraform output

   # View sensitive outputs
   terraform output smtp_password
   terraform output api_key_secret
   ```

## Testing

After applying, you can verify the resources were created:

- **Check Mailgun dashboard**: https://app.mailgun.com/
- **Test SMTP credentials**: Use the output credentials to send a test email
- **Test webhooks**: Send an email and verify your webhook endpoint receives events
- **Test routes**: Send an email to `support@your-domain` and verify it's forwarded

## Cleanup

To destroy all resources:

```bash
terraform destroy
```

## Import Existing Resources

If you have existing Mailgun resources, you can import them:

```bash
# Import a webhook
terraform import mailgun_webhook.delivered example.com/delivered

# Import SMTP credentials
terraform import mailgun_smtp_credential.app example.com/app@example.com

# Import a route
terraform import mailgun_route.support 4f3bad2335335426750048c6

# Import an API key
terraform import mailgun_api_key.ci_cd key-abc123
```

## Notes

- **API Key Secret**: Only available on creation. Store it securely!
- **SMTP Password**: Stored in Terraform state (encrypted at rest)
- **Domain Verification**: You'll need to add DNS records to verify your domain
- **Routes**: Use Mailgun's filter expression syntax for matching
