// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_domain

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &DomainResource{}
	_ resource.ResourceWithConfigure = &DomainResource{}
)

// DomainResource is the resource implementation.
type DomainResource struct {
	client *mailgun.Client
}

// Metadata returns the resource type name.
func (r *DomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

// Schema defines the schema for the resource.
func (r *DomainResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DomainResourceSchema(ctx)
}

// Configure adds the provider configured client to the resource.
func (r *DomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*mailgun.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *mailgun.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates a new resource.
func (r *DomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DomainModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that client is configured
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Mailgun Client",
			"The Mailgun client has not been properly configured. "+
				"Please report this issue to the provider developers.",
		)
		return
	}

	// Validate required fields
	if plan.Name.IsNull() || plan.Name.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Domain Name",
			"The domain name is required to create a domain.",
		)
		return
	}

	domainName := plan.Name.ValueString()

	// Build CreateDomainOptions
	opts := &mailgun.CreateDomainOptions{}

	if !plan.SmtpPassword.IsNull() {
		opts.Password = plan.SmtpPassword.ValueString()
	}

	if !plan.SpamAction.IsNull() {
		opts.SpamAction = mtypes.SpamAction(plan.SpamAction.ValueString())
	}

	if !plan.ForceDkimAuthority.IsNull() {
		opts.ForceDKIMAuthority = plan.ForceDkimAuthority.ValueBool()
	}

	if !plan.DkimKeySize.IsNull() {
		// DkimKeySize is stored as string in the model, convert to int
		var keySize int
		switch plan.DkimKeySize.ValueString() {
		case "1024":
			keySize = 1024
		case "2048":
			keySize = 2048
		}
		if keySize != 0 {
			opts.DKIMKeySize = keySize
		}
	}

	if !plan.UseAutomaticSenderSecurity.IsNull() {
		opts.UseAutomaticSenderSecurity = plan.UseAutomaticSenderSecurity.ValueBool()
	}

	// Create context with timeout
	createCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Create the domain via Mailgun API
	domainResp, err := r.client.CreateDomain(createCtx, domainName, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Domain",
			fmt.Sprintf("Could not create domain %s: %s", domainName, err),
		)
		return
	}

	// Map response to plan model
	plan = mapDomainResponseToModel(domainResp, plan)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *DomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DomainModel

	// Read current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that client is configured
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Mailgun Client",
			"The Mailgun client has not been properly configured. "+
				"Please report this issue to the provider developers.",
		)
		return
	}

	// Get domain name from state
	domainName := state.Name.ValueString()
	if domainName == "" {
		resp.Diagnostics.AddError(
			"Missing Domain Name",
			"The domain name is required to read a domain.",
		)
		return
	}

	// Create context with timeout
	readCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Get the domain via Mailgun API
	domainResp, err := r.client.GetDomain(readCtx, domainName, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Domain",
			fmt.Sprintf("Could not read domain %s: %s", domainName, err),
		)
		return
	}

	// Map response to state model
	state = mapDomainResponseToModel(domainResp, state)

	// Save updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates an existing resource.
func (r *DomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DomainModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Note: Mailgun domains are largely immutable after creation.
	// Most changes require domain recreation. For now, we'll just
	// refresh the state by reading the domain.

	// Validate that client is configured
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Mailgun Client",
			"The Mailgun client has not been properly configured. "+
				"Please report this issue to the provider developers.",
		)
		return
	}

	domainName := plan.Name.ValueString()

	// Create context with timeout
	updateCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Get the latest domain state
	domainResp, err := r.client.GetDomain(updateCtx, domainName, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Domain",
			fmt.Sprintf("Could not read domain %s: %s", domainName, err),
		)
		return
	}

	// Map response to plan model
	plan = mapDomainResponseToModel(domainResp, plan)

	// Save updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes an existing resource.
func (r *DomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DomainModel

	// Read current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that client is configured
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Mailgun Client",
			"The Mailgun client has not been properly configured. "+
				"Please report this issue to the provider developers.",
		)
		return
	}

	domainName := state.Name.ValueString()
	if domainName == "" {
		resp.Diagnostics.AddError(
			"Missing Domain Name",
			"The domain name is required to delete a domain.",
		)
		return
	}

	// Create context with timeout
	deleteCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Delete the domain via Mailgun API
	err := r.client.DeleteDomain(deleteCtx, domainName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Domain",
			fmt.Sprintf("Could not delete domain %s: %s", domainName, err),
		)
		return
	}

	// State is automatically removed by Terraform after successful deletion
}

// mapDomainResponseToModel maps Mailgun API response to Terraform model
func mapDomainResponseToModel(domainResp mtypes.GetDomainResponse, model DomainModel) DomainModel {
	// Map domain fields
	domainValue := NewDomainValueMust(model.Domain.AttributeTypes(context.Background()), map[string]attr.Value{
		"created_at":                    types.StringValue(domainResp.Domain.CreatedAt.String()),
		"disabled":                      NewDisabledValueNull(),
		"id":                            types.StringValue(domainResp.Domain.ID),
		"is_disabled":                   types.BoolValue(domainResp.Domain.IsDisabled),
		"name":                          types.StringValue(domainResp.Domain.Name),
		"require_tls":                   types.BoolValue(domainResp.Domain.RequireTLS),
		"skip_verification":             types.BoolValue(domainResp.Domain.SkipVerification),
		"smtp_login":                    types.StringValue(domainResp.Domain.SMTPLogin),
		"smtp_password":                 types.StringValue(domainResp.Domain.SMTPPassword),
		"spam_action":                   types.StringValue(string(domainResp.Domain.SpamAction)),
		"state":                         types.StringValue(domainResp.Domain.State),
		"tracking_host":                 types.StringValue(domainResp.Domain.TrackingHost),
		"type":                          types.StringValue(domainResp.Domain.Type),
		"use_automatic_sender_security": types.BoolValue(domainResp.Domain.UseAutomaticSenderSecurity),
		"web_prefix":                    types.StringValue(domainResp.Domain.WebPrefix),
		"web_scheme":                    types.StringValue(domainResp.Domain.WebScheme),
		"wildcard":                      types.BoolValue(domainResp.Domain.Wildcard),
	})

	model.Domain = domainValue
	model.Name = types.StringValue(domainResp.Domain.Name)
	model.UseAutomaticSenderSecurity = types.BoolValue(domainResp.UseAutomaticSenderSecurity)

	// TODO: Map DNS records when needed
	// model.ReceivingDnsRecords = ...
	// model.SendingDnsRecords = ...

	return model
}
