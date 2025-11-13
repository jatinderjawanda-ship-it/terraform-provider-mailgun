// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package domains

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DomainResourceSchema returns the schema for the domain resource
func DomainResourceSchema() rschema.Schema {
	return rschema.Schema{
		Description: "Manages a Mailgun domain. Domains are used to send and receive email.",
		Attributes: map[string]rschema.Attribute{
			"name": rschema.StringAttribute{
				Description: "The domain name to be used for sending and receiving email.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"spam_action": rschema.StringAttribute{
				Description: "Spam filter action for new domain. Options: 'disabled', 'tag', or 'delete'.",
				Optional:    true,
				Computed:    true,
			},
			"wildcard": rschema.BoolAttribute{
				Description: "Determines whether the domain will accept email for sub-domains.",
				Optional:    true,
				Computed:    true,
			},
			"force_dkim_authority": rschema.BoolAttribute{
				Description: "If set to true, the domain will be the DKIM authority for itself even if the root domain is registered on the same mailgun account.",
				Optional:    true,
				Computed:    true,
			},
			"dkim_key_size": rschema.StringAttribute{
				Description: "DKIM key size (1024 or 2048).",
				Optional:    true,
				Computed:    true,
			},
			"ips": rschema.StringAttribute{
				Description: "Comma-separated list of IPs to be assigned to this domain.",
				Optional:    true,
			},
			"pool_id": rschema.StringAttribute{
				Description: "The id of the IP Pool that you wish to assign to the domain.",
				Optional:    true,
			},
			"web_scheme": rschema.StringAttribute{
				Description: "Web scheme for tracking links (http or https).",
				Optional:    true,
				Computed:    true,
			},
			"web_prefix": rschema.StringAttribute{
				Description: "Web prefix for tracking links.",
				Optional:    true,
				Computed:    true,
			},
			"use_automatic_sender_security": rschema.BoolAttribute{
				Description: "Whether to use automatic sender security (SPF/DKIM/DMARC).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"smtp_password": rschema.StringAttribute{
				Description: "Password for SMTP authentication.",
				Optional:    true,
				Sensitive:   true,
			},
			"dkim_selector": rschema.StringAttribute{
				Description: "DKIM selector (custom CNAME for DKIM).",
				Optional:    true,
			},
			"dkim_host_name": rschema.StringAttribute{
				Description: "DKIM host name.",
				Optional:    true,
			},
			"force_root_dkim_host": rschema.BoolAttribute{
				Description: "Force using root domain for DKIM.",
				Optional:    true,
			},
			"encrypt_incoming_message": rschema.BoolAttribute{
				Description: "Whether to encrypt incoming messages.",
				Optional:    true,
			},
			"hextended": rschema.BoolAttribute{
				Description: "Whether to use extended headers.",
				Optional:    true,
			},
			"hwith_dns": rschema.BoolAttribute{
				Description: "Whether to use DNS for headers.",
				Optional:    true,
			},
			"message": rschema.StringAttribute{
				Description: "Custom message for domain creation.",
				Optional:    true,
			},
			// Computed attributes from API response
			"domain": DomainNestedObject(),
			"receiving_dns_records": rschema.ListNestedAttribute{
				Description: "DNS records for receiving email.",
				Computed:    true,
				NestedObject: rschema.NestedAttributeObject{
					Attributes: map[string]rschema.Attribute{
						"cached": rschema.ListAttribute{
							Description: "Cached DNS values.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"is_active": rschema.BoolAttribute{
							Description: "Whether the DNS record is active.",
							Computed:    true,
						},
						"name": rschema.StringAttribute{
							Description: "DNS record name.",
							Computed:    true,
						},
						"priority": rschema.StringAttribute{
							Description: "DNS record priority.",
							Computed:    true,
						},
						"record_type": rschema.StringAttribute{
							Description: "DNS record type.",
							Computed:    true,
						},
						"valid": rschema.StringAttribute{
							Description: "Whether the DNS record is valid.",
							Computed:    true,
						},
						"value": rschema.StringAttribute{
							Description: "DNS record value.",
							Computed:    true,
						},
					},
				},
			},
			"sending_dns_records": rschema.ListNestedAttribute{
				Description: "DNS records for sending email.",
				Computed:    true,
				NestedObject: rschema.NestedAttributeObject{
					Attributes: map[string]rschema.Attribute{
						"cached": rschema.ListAttribute{
							Description: "Cached DNS values.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"is_active": rschema.BoolAttribute{
							Description: "Whether the DNS record is active.",
							Computed:    true,
						},
						"name": rschema.StringAttribute{
							Description: "DNS record name.",
							Computed:    true,
						},
						"priority": rschema.StringAttribute{
							Description: "DNS record priority.",
							Computed:    true,
						},
						"record_type": rschema.StringAttribute{
							Description: "DNS record type.",
							Computed:    true,
						},
						"valid": rschema.StringAttribute{
							Description: "Whether the DNS record is valid.",
							Computed:    true,
						},
						"value": rschema.StringAttribute{
							Description: "DNS record value.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// DomainNestedObject returns the nested object schema for domain details
func DomainNestedObject() rschema.SingleNestedAttribute {
	return rschema.SingleNestedAttribute{
		Description: "Domain details returned from the API.",
		Computed:    true,
		Attributes: map[string]rschema.Attribute{
			"created_at": rschema.StringAttribute{
				Description: "Time the domain was created.",
				Computed:    true,
			},
			"disabled": rschema.SingleNestedAttribute{
				Description: "Information about domain being disabled.",
				Computed:    true,
				Attributes: map[string]rschema.Attribute{
					"code": rschema.StringAttribute{
						Description: "Disable code.",
						Computed:    true,
					},
					"note": rschema.StringAttribute{
						Description: "Disable note.",
						Computed:    true,
					},
					"permanently": rschema.BoolAttribute{
						Description: "Whether disabled permanently.",
						Computed:    true,
					},
					"reason": rschema.StringAttribute{
						Description: "Disable reason.",
						Computed:    true,
					},
					"until": rschema.StringAttribute{
						Description: "Disabled until timestamp.",
						Computed:    true,
					},
				},
			},
			"id": rschema.StringAttribute{
				Description: "Domain ID.",
				Computed:    true,
			},
			"is_disabled": rschema.BoolAttribute{
				Description: "Whether the domain is disabled.",
				Computed:    true,
			},
			"name": rschema.StringAttribute{
				Description: "Domain name.",
				Computed:    true,
			},
			"require_tls": rschema.BoolAttribute{
				Description: "Whether TLS is required.",
				Computed:    true,
			},
			"skip_verification": rschema.BoolAttribute{
				Description: "Whether to skip verification.",
				Computed:    true,
			},
			"smtp_login": rschema.StringAttribute{
				Description: "SMTP login username.",
				Computed:    true,
			},
			"smtp_password": rschema.StringAttribute{
				Description: "SMTP password.",
				Computed:    true,
				Sensitive:   true,
			},
			"spam_action": rschema.StringAttribute{
				Description: "Spam action setting.",
				Computed:    true,
			},
			"state": rschema.StringAttribute{
				Description: "Domain state.",
				Computed:    true,
			},
			"tracking_host": rschema.StringAttribute{
				Description: "Tracking host.",
				Computed:    true,
			},
			"type": rschema.StringAttribute{
				Description: "Domain type.",
				Computed:    true,
			},
			"use_automatic_sender_security": rschema.BoolAttribute{
				Description: "Whether automatic sender security is enabled.",
				Computed:    true,
			},
			"web_prefix": rschema.StringAttribute{
				Description: "Web prefix for tracking.",
				Computed:    true,
			},
			"web_scheme": rschema.StringAttribute{
				Description: "Web scheme (http or https).",
				Computed:    true,
			},
			"wildcard": rschema.BoolAttribute{
				Description: "Whether wildcard is enabled.",
				Computed:    true,
			},
		},
	}
}

// DomainsDataSourceSchema returns the schema for the domains data source
func DomainsDataSourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Retrieves a list of all domains in your Mailgun account.",
		Attributes: map[string]schema.Attribute{
			"authority": schema.StringAttribute{
				Description: "Filter by authority.",
				Optional:    true,
			},
			"include_subaccounts": schema.BoolAttribute{
				Description: "Whether to include subaccounts.",
				Optional:    true,
			},
			"limit": schema.Int64Attribute{
				Description: "Maximum number of domains to return (default: 100).",
				Optional:    true,
			},
			"search": schema.StringAttribute{
				Description: "Search query to filter domains.",
				Optional:    true,
			},
			"skip": schema.Int64Attribute{
				Description: "Number of domains to skip.",
				Optional:    true,
			},
			"sort": schema.StringAttribute{
				Description: "Sort order for results.",
				Optional:    true,
			},
			"state": schema.StringAttribute{
				Description: "Filter by domain state.",
				Optional:    true,
			},
			"total_count": schema.Int64Attribute{
				Description: "Total number of domains returned.",
				Computed:    true,
			},
			"items": schema.ListNestedAttribute{
				Description: "List of domains.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Description: "Time the domain was created.",
							Computed:    true,
						},
						"disabled": schema.SingleNestedAttribute{
							Description: "Information about domain being disabled.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"code": schema.StringAttribute{
									Description: "Disable code.",
									Computed:    true,
								},
								"note": schema.StringAttribute{
									Description: "Disable note.",
									Computed:    true,
								},
								"permanently": schema.BoolAttribute{
									Description: "Whether disabled permanently.",
									Computed:    true,
								},
								"reason": schema.StringAttribute{
									Description: "Disable reason.",
									Computed:    true,
								},
								"until": schema.StringAttribute{
									Description: "Disabled until timestamp.",
									Computed:    true,
								},
							},
						},
						"id": schema.StringAttribute{
							Description: "Domain ID.",
							Computed:    true,
						},
						"is_disabled": schema.BoolAttribute{
							Description: "Whether the domain is disabled.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Domain name.",
							Computed:    true,
						},
						"require_tls": schema.BoolAttribute{
							Description: "Whether TLS is required.",
							Computed:    true,
						},
						"skip_verification": schema.BoolAttribute{
							Description: "Whether to skip verification.",
							Computed:    true,
						},
						"smtp_login": schema.StringAttribute{
							Description: "SMTP login username.",
							Computed:    true,
						},
						"smtp_password": schema.StringAttribute{
							Description: "SMTP password.",
							Computed:    true,
							Sensitive:   true,
						},
						"spam_action": schema.StringAttribute{
							Description: "Spam action setting.",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "Domain state.",
							Computed:    true,
						},
						"tracking_host": schema.StringAttribute{
							Description: "Tracking host.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Domain type.",
							Computed:    true,
						},
						"use_automatic_sender_security": schema.BoolAttribute{
							Description: "Whether automatic sender security is enabled.",
							Computed:    true,
						},
						"web_prefix": schema.StringAttribute{
							Description: "Web prefix for tracking.",
							Computed:    true,
						},
						"web_scheme": schema.StringAttribute{
							Description: "Web scheme (http or https).",
							Computed:    true,
						},
						"wildcard": schema.BoolAttribute{
							Description: "Whether wildcard is enabled.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// NewDomainResource creates a new domain resource
func NewDomainResource() resource.Resource {
	return &DomainResource{}
}

// NewDomainsDataSource creates a new domains data source
func NewDomainsDataSource() datasource.DataSource {
	return &DataSource{}
}
