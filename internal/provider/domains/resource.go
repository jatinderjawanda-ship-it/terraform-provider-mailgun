// Copyright (c) Hack The Box
// SPDX-License-Identifier: MPL-2.0

package domains

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
	_ resource.Resource                = &DomainResource{}
	_ resource.ResourceWithConfigure   = &DomainResource{}
	_ resource.ResourceWithImportState = &DomainResource{}
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
	resp.Schema = DomainResourceSchema()
}

// Configure adds the provider-configured client to the resource.
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

	// Map response to a plan model
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

	// Map response to a state model
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

	// Map response to a plan model
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

// ImportState imports an existing domain by name
func (r *DomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The ID passed is the domain name
	domainName := req.ID

	// Validate that client is configured
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured Mailgun Client",
			"The Mailgun client has not been properly configured. "+
				"Please report this issue to the provider developers.",
		)
		return
	}

	// Create context with timeout
	importCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Get the domain via Mailgun API
	domainResp, err := r.client.GetDomain(importCtx, domainName, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Domain",
			fmt.Sprintf("Could not import domain %s: %s", domainName, err),
		)
		return
	}

	// Create a new model with imported values
	var state DomainModel
	state.Name = types.StringValue(domainName)
	state = mapDomainResponseToModel(domainResp, state)

	// Save imported state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// mapDomainResponseToModel maps Mailgun API response to Terraform model
func mapDomainResponseToModel(domainResp mtypes.GetDomainResponse, model DomainModel) DomainModel {
	ctx := context.Background()

	// Create disabled as null object (SDK doesn't parse this field)
	disabledAttrTypes := map[string]attr.Type{
		"code":        types.StringType,
		"note":        types.StringType,
		"permanently": types.BoolType,
		"reason":      types.StringType,
		"until":       types.StringType,
	}

	// Map domain fields from SDK response
	domainValue := NewDomainValueMust(model.Domain.AttributeTypes(ctx), map[string]attr.Value{
		"created_at":                    types.StringValue(domainResp.Domain.CreatedAt.String()),
		"disabled":                      types.ObjectNull(disabledAttrTypes),
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

	// Map DNS records from response
	if len(domainResp.ReceivingDNSRecords) > 0 {
		receivingRecords := make([]ReceivingDnsRecordsValue, len(domainResp.ReceivingDNSRecords))
		for i, record := range domainResp.ReceivingDNSRecords {
			cachedList, _ := types.ListValueFrom(ctx, types.StringType, record.Cached)
			receivingRecords[i] = ReceivingDnsRecordsValue{
				Cached:     cachedList,
				IsActive:   types.BoolValue(record.Active),
				Name:       types.StringValue(record.Name),
				Priority:   types.StringValue(record.Priority),
				RecordType: types.StringValue(record.RecordType),
				Valid:      types.StringValue(record.Valid),
				Value:      types.StringValue(record.Value),
				state:      attr.ValueStateKnown,
			}
		}
		receivingList, _ := types.ListValueFrom(ctx, ReceivingDnsRecordsType{
			ObjectType: types.ObjectType{AttrTypes: ReceivingDnsRecordsValue{}.AttributeTypes(ctx)},
		}, receivingRecords)
		model.ReceivingDnsRecords = receivingList
	} else if model.ReceivingDnsRecords.IsUnknown() {
		receivingDnsRecordsType := ReceivingDnsRecordsType{
			ObjectType: types.ObjectType{AttrTypes: ReceivingDnsRecordsValue{}.AttributeTypes(ctx)},
		}
		model.ReceivingDnsRecords = types.ListNull(receivingDnsRecordsType)
	}

	if len(domainResp.SendingDNSRecords) > 0 {
		sendingRecords := make([]SendingDnsRecordsValue, len(domainResp.SendingDNSRecords))
		for i, record := range domainResp.SendingDNSRecords {
			cachedList, _ := types.ListValueFrom(ctx, types.StringType, record.Cached)
			sendingRecords[i] = SendingDnsRecordsValue{
				Cached:     cachedList,
				IsActive:   types.BoolValue(record.Active),
				Name:       types.StringValue(record.Name),
				Priority:   types.StringValue(record.Priority),
				RecordType: types.StringValue(record.RecordType),
				Valid:      types.StringValue(record.Valid),
				Value:      types.StringValue(record.Value),
				state:      attr.ValueStateKnown,
			}
		}
		sendingList, _ := types.ListValueFrom(ctx, SendingDnsRecordsType{
			ObjectType: types.ObjectType{AttrTypes: SendingDnsRecordsValue{}.AttributeTypes(ctx)},
		}, sendingRecords)
		model.SendingDnsRecords = sendingList
	} else if model.SendingDnsRecords.IsUnknown() {
		sendingDnsRecordsType := SendingDnsRecordsType{
			ObjectType: types.ObjectType{AttrTypes: SendingDnsRecordsValue{}.AttributeTypes(ctx)},
		}
		model.SendingDnsRecords = types.ListNull(sendingDnsRecordsType)
	}

	// Set all other fields from plan or to null (SDK doesn't provide them in response)
	// These are request-only fields that don't come back in the response
	if model.DkimHostName.IsUnknown() {
		model.DkimHostName = types.StringNull()
	}
	if model.DkimKeySize.IsUnknown() {
		model.DkimKeySize = types.StringNull()
	}
	if model.DkimSelector.IsUnknown() {
		model.DkimSelector = types.StringNull()
	}
	if model.EncryptIncomingMessage.IsUnknown() {
		model.EncryptIncomingMessage = types.BoolNull()
	}
	if model.ForceDkimAuthority.IsUnknown() {
		model.ForceDkimAuthority = types.BoolNull()
	}
	if model.ForceRootDkimHost.IsUnknown() {
		model.ForceRootDkimHost = types.BoolNull()
	}
	if model.Hextended.IsUnknown() {
		model.Hextended = types.BoolNull()
	}
	if model.HwithDns.IsUnknown() {
		model.HwithDns = types.BoolNull()
	}
	if model.Ips.IsUnknown() {
		model.Ips = types.StringNull()
	}
	if model.Message.IsUnknown() {
		model.Message = types.StringNull()
	}
	if model.PoolId.IsUnknown() {
		model.PoolId = types.StringNull()
	}
	// DNS records already mapped above
	if model.SmtpPassword.IsUnknown() {
		model.SmtpPassword = types.StringNull()
	}
	if model.WebPrefix.IsUnknown() {
		model.WebPrefix = types.StringNull()
	}
	if model.WebScheme.IsUnknown() {
		model.WebScheme = types.StringNull()
	}

	return model
}
