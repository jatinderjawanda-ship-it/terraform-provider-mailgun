// Copyright (c) Hack The Box
// SPDX-License-Identifier: MPL-2.0

package domains

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// DomainModel represents the Terraform state for a domain resource
type DomainModel struct {
	Name                       types.String `tfsdk:"name"`
	SpamAction                 types.String `tfsdk:"spam_action"`
	Wildcard                   types.Bool   `tfsdk:"wildcard"`
	ForceDkimAuthority         types.Bool   `tfsdk:"force_dkim_authority"`
	DkimKeySize                types.String `tfsdk:"dkim_key_size"`
	Ips                        types.String `tfsdk:"ips"`
	PoolId                     types.String `tfsdk:"pool_id"`
	WebScheme                  types.String `tfsdk:"web_scheme"`
	WebPrefix                  types.String `tfsdk:"web_prefix"`
	UseAutomaticSenderSecurity types.Bool   `tfsdk:"use_automatic_sender_security"`
	SmtpPassword               types.String `tfsdk:"smtp_password"`
	DkimSelector               types.String `tfsdk:"dkim_selector"`
	DkimHostName               types.String `tfsdk:"dkim_host_name"`
	ForceRootDkimHost          types.Bool   `tfsdk:"force_root_dkim_host"`
	EncryptIncomingMessage     types.Bool   `tfsdk:"encrypt_incoming_message"`
	Hextended                  types.Bool   `tfsdk:"hextended"`
	HwithDns                   types.Bool   `tfsdk:"hwith_dns"`
	Message                    types.String `tfsdk:"message"`
	// Computed attributes
	Domain              DomainValue `tfsdk:"domain"`
	ReceivingDnsRecords types.List  `tfsdk:"receiving_dns_records"`
	SendingDnsRecords   types.List  `tfsdk:"sending_dns_records"`
}

// DomainValue represents the nested domain object
type DomainValue struct {
	CreatedAt                  types.String `tfsdk:"created_at"`
	Disabled                   types.Object `tfsdk:"disabled"`
	Id                         types.String `tfsdk:"id"`
	IsDisabled                 types.Bool   `tfsdk:"is_disabled"`
	Name                       types.String `tfsdk:"name"`
	RequireTls                 types.Bool   `tfsdk:"require_tls"`
	SkipVerification           types.Bool   `tfsdk:"skip_verification"`
	SmtpLogin                  types.String `tfsdk:"smtp_login"`
	SmtpPassword               types.String `tfsdk:"smtp_password"`
	SpamAction                 types.String `tfsdk:"spam_action"`
	State                      types.String `tfsdk:"state"`
	TrackingHost               types.String `tfsdk:"tracking_host"`
	DomainType                 types.String `tfsdk:"type"`
	UseAutomaticSenderSecurity types.Bool   `tfsdk:"use_automatic_sender_security"`
	WebPrefix                  types.String `tfsdk:"web_prefix"`
	WebScheme                  types.String `tfsdk:"web_scheme"`
	Wildcard                   types.Bool   `tfsdk:"wildcard"`
	state                      attr.ValueState
}

// AttributeTypes returns the attribute types for DomainValue
func (v DomainValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"created_at": types.StringType,
		"disabled": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"code":        types.StringType,
				"note":        types.StringType,
				"permanently": types.BoolType,
				"reason":      types.StringType,
				"until":       types.StringType,
			},
		},
		"id":                            types.StringType,
		"is_disabled":                   types.BoolType,
		"name":                          types.StringType,
		"require_tls":                   types.BoolType,
		"skip_verification":             types.BoolType,
		"smtp_login":                    types.StringType,
		"smtp_password":                 types.StringType,
		"spam_action":                   types.StringType,
		"state":                         types.StringType,
		"tracking_host":                 types.StringType,
		"type":                          types.StringType,
		"use_automatic_sender_security": types.BoolType,
		"web_prefix":                    types.StringType,
		"web_scheme":                    types.StringType,
		"wildcard":                      types.BoolType,
	}
}

// Type returns the type for DomainValue
func (v DomainValue) Type(ctx context.Context) attr.Type {
	return types.ObjectType{
		AttrTypes: v.AttributeTypes(ctx),
	}
}

// ToObjectValue converts DomainValue to an ObjectValue
func (v DomainValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(v.AttributeTypes(ctx), map[string]attr.Value{
		"created_at":                    v.CreatedAt,
		"disabled":                      v.Disabled,
		"id":                            v.Id,
		"is_disabled":                   v.IsDisabled,
		"name":                          v.Name,
		"require_tls":                   v.RequireTls,
		"skip_verification":             v.SkipVerification,
		"smtp_login":                    v.SmtpLogin,
		"smtp_password":                 v.SmtpPassword,
		"spam_action":                   v.SpamAction,
		"state":                         v.State,
		"tracking_host":                 v.TrackingHost,
		"type":                          v.DomainType,
		"use_automatic_sender_security": v.UseAutomaticSenderSecurity,
		"web_prefix":                    v.WebPrefix,
		"web_scheme":                    v.WebScheme,
		"wildcard":                      v.Wildcard,
	})
}

// NewDomainValueMust creates a DomainValue from a map of attribute values
func NewDomainValueMust(attrTypes map[string]attr.Type, attributes map[string]attr.Value) DomainValue {
	return DomainValue{
		CreatedAt:                  attributes["created_at"].(types.String),
		Disabled:                   attributes["disabled"].(types.Object),
		Id:                         attributes["id"].(types.String),
		IsDisabled:                 attributes["is_disabled"].(types.Bool),
		Name:                       attributes["name"].(types.String),
		RequireTls:                 attributes["require_tls"].(types.Bool),
		SkipVerification:           attributes["skip_verification"].(types.Bool),
		SmtpLogin:                  attributes["smtp_login"].(types.String),
		SmtpPassword:               attributes["smtp_password"].(types.String),
		SpamAction:                 attributes["spam_action"].(types.String),
		State:                      attributes["state"].(types.String),
		TrackingHost:               attributes["tracking_host"].(types.String),
		DomainType:                 attributes["type"].(types.String),
		UseAutomaticSenderSecurity: attributes["use_automatic_sender_security"].(types.Bool),
		WebPrefix:                  attributes["web_prefix"].(types.String),
		WebScheme:                  attributes["web_scheme"].(types.String),
		Wildcard:                   attributes["wildcard"].(types.Bool),
		state:                      attr.ValueStateKnown,
	}
}

// DisabledValue represents the nested disabled object
type DisabledValue struct {
	Code        types.String `tfsdk:"code"`
	Note        types.String `tfsdk:"note"`
	Permanently types.Bool   `tfsdk:"permanently"`
	Reason      types.String `tfsdk:"reason"`
	Until       types.String `tfsdk:"until"`
	state       attr.ValueState
}

// AttributeTypes returns the attribute types for DisabledValue
func (v DisabledValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"code":        types.StringType,
		"note":        types.StringType,
		"permanently": types.BoolType,
		"reason":      types.StringType,
		"until":       types.StringType,
	}
}

// ToObjectValue converts DisabledValue to an ObjectValue
func (v DisabledValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(v.AttributeTypes(ctx), map[string]attr.Value{
		"code":        v.Code,
		"note":        v.Note,
		"permanently": v.Permanently,
		"reason":      v.Reason,
		"until":       v.Until,
	})
}

// NewDisabledValueNull creates a null DisabledValue
func NewDisabledValueNull() DisabledValue {
	return DisabledValue{
		Code:        types.StringNull(),
		Note:        types.StringNull(),
		Permanently: types.BoolNull(),
		Reason:      types.StringNull(),
		Until:       types.StringNull(),
		state:       attr.ValueStateNull,
	}
}

// ReceivingDnsRecordsValue represents a DNS record for receiving
type ReceivingDnsRecordsValue struct {
	Cached     types.List   `tfsdk:"cached"`
	IsActive   types.Bool   `tfsdk:"is_active"`
	Name       types.String `tfsdk:"name"`
	Priority   types.String `tfsdk:"priority"`
	RecordType types.String `tfsdk:"record_type"`
	Valid      types.String `tfsdk:"valid"`
	Value      types.String `tfsdk:"value"`
	state      attr.ValueState
}

// AttributeTypes returns the attribute types for ReceivingDnsRecordsValue
func (v ReceivingDnsRecordsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"cached":      types.ListType{ElemType: types.StringType},
		"is_active":   types.BoolType,
		"name":        types.StringType,
		"priority":    types.StringType,
		"record_type": types.StringType,
		"valid":       types.StringType,
		"value":       types.StringType,
	}
}

// ToObjectValue converts ReceivingDnsRecordsValue to an ObjectValue
func (v ReceivingDnsRecordsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(v.AttributeTypes(ctx), map[string]attr.Value{
		"cached":      v.Cached,
		"is_active":   v.IsActive,
		"name":        v.Name,
		"priority":    v.Priority,
		"record_type": v.RecordType,
		"valid":       v.Valid,
		"value":       v.Value,
	})
}

// ReceivingDnsRecordsType represents the list type for receiving DNS records
type ReceivingDnsRecordsType struct {
	types.ObjectType
}

// SendingDnsRecordsValue represents a DNS record for sending
type SendingDnsRecordsValue struct {
	Cached     types.List   `tfsdk:"cached"`
	IsActive   types.Bool   `tfsdk:"is_active"`
	Name       types.String `tfsdk:"name"`
	Priority   types.String `tfsdk:"priority"`
	RecordType types.String `tfsdk:"record_type"`
	Valid      types.String `tfsdk:"valid"`
	Value      types.String `tfsdk:"value"`
	state      attr.ValueState
}

// AttributeTypes returns the attribute types for SendingDnsRecordsValue
func (v SendingDnsRecordsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"cached":      types.ListType{ElemType: types.StringType},
		"is_active":   types.BoolType,
		"name":        types.StringType,
		"priority":    types.StringType,
		"record_type": types.StringType,
		"valid":       types.StringType,
		"value":       types.StringType,
	}
}

// ToObjectValue converts SendingDnsRecordsValue to an ObjectValue
func (v SendingDnsRecordsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(v.AttributeTypes(ctx), map[string]attr.Value{
		"cached":      v.Cached,
		"is_active":   v.IsActive,
		"name":        v.Name,
		"priority":    v.Priority,
		"record_type": v.RecordType,
		"valid":       v.Valid,
		"value":       v.Value,
	})
}

// SendingDnsRecordsType represents the list type for sending DNS records
type SendingDnsRecordsType struct {
	types.ObjectType
}
