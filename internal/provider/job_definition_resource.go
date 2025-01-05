package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/eagletmt/terraform-provider-kuroko2/internal/kuroko2"
)

var (
	_ resource.Resource                = (*jobDefinitionResource)(nil)
	_ resource.ResourceWithImportState = (*jobDefinitionResource)(nil)
)

type jobDefinitionResource struct {
	client kuroko2.Client
}

func NewJobDefinitionResource() resource.Resource {
	return &jobDefinitionResource{}
}

func (r *jobDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job_definition"
}

func (r *jobDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
			"script": schema.StringAttribute{
				Required: true,
			},
			"admins": schema.ListAttribute{
				ElementType: types.Int64Type,
				Required:    true,
			},
			"cron": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"notify_cancellation": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"suspended": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"prevent_multi": schema.StringAttribute{
				Optional:   true,
				Computed:   true,
				Default:    stringdefault.StaticString("WORKING_OR_ERROR"),
				Validators: []validator.String{stringvalidator.OneOf("NONE", "WORKING_OR_ERROR", "WORKING", "ERROR")},
			},
			"slack_channel": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
		},
	}
}

func (r *jobDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*kuroko2ProviderData).client
}

type jobDefinitionResourceModel struct {
	Id                 types.Int64  `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Script             types.String `tfsdk:"script"`
	Admins             types.List   `tfsdk:"admins"`
	Cron               types.List   `tfsdk:"cron"`
	Tags               types.List   `tfsdk:"tags"`
	NotifyCancellation types.Bool   `tfsdk:"notify_cancellation"`
	Suspended          types.Bool   `tfsdk:"suspended"`
	PreventMulti       types.String `tfsdk:"prevent_multi"`
	SlackChannel       types.String `tfsdk:"slack_channel"`
}

func (r *jobDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data jobDefinitionResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, diags := decodeJobDefinition(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	definition, err := r.client.CreateJobDefinition(ctx, model)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to create a job definition: %s", err.Error()))
		return
	}

	tfmodel, diags := encodeJobDefinition(ctx, definition)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tfmodel)...)
}

func (r *jobDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data jobDefinitionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	definition, err := r.client.GetJobDefinition(ctx, data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to get a job definition: %s", err.Error()))
		return
	}

	tfmodel, diags := encodeJobDefinition(ctx, definition)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, tfmodel)...)
}

func (r *jobDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data jobDefinitionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, diags := decodeJobDefinition(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateJobDefinition(ctx, data.Id.ValueInt64(), model)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to update a job definition: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *jobDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data jobDefinitionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteJobDefinition(ctx, data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to delete a job definition: %s", err.Error()))
		return
	}
}

func (r *jobDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Failed to import a job definition: %s", err.Error()))
		return
	}
	resp.State.SetAttribute(ctx, path.Root("id"), id)
}

func encodeJobDefinition(ctx context.Context, definition kuroko2.JobDefinition) (*jobDefinitionResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	admins, ds := types.ListValueFrom(ctx, types.Int64Type, definition.Admins)
	diags.Append(ds...)
	if diags.HasError() {
		return nil, diags
	}

	cron := types.ListNull(types.StringType)
	if len(definition.Cron) > 0 {
		cron, ds = types.ListValueFrom(ctx, types.StringType, definition.Cron)
		diags.Append(ds...)
		if diags.HasError() {
			return nil, diags
		}
	}

	tags := types.ListNull(types.StringType)
	if len(definition.Tags) > 0 {
		tags, ds = types.ListValueFrom(ctx, types.StringType, definition.Tags)
		diags.Append(ds...)
		if diags.HasError() {
			return nil, diags
		}
	}

	var preventMulti string
	switch definition.PreventMulti {
	case 0:
		preventMulti = "NONE"
	case 1:
		preventMulti = "WORKING_OR_ERROR"
	case 2:
		preventMulti = "WORKING"
	case 3:
		preventMulti = "ERROR"
	default:
		panic(fmt.Sprintf("Unknown prevent_multi value: %d", definition.PreventMulti))
	}

	return &jobDefinitionResourceModel{
		Id:                 types.Int64Value(definition.Id),
		Name:               types.StringValue(definition.Name),
		Description:        types.StringValue(definition.Description),
		Script:             types.StringValue(definition.Script),
		Admins:             admins,
		Cron:               cron,
		Tags:               tags,
		NotifyCancellation: types.BoolValue(definition.NotifyCancellation),
		Suspended:          types.BoolValue(definition.Suspended),
		PreventMulti:       types.StringValue(preventMulti),
		SlackChannel:       types.StringValue(definition.SlackChannel),
	}, diags
}

func decodeJobDefinition(ctx context.Context, data jobDefinitionResourceModel) (kuroko2.JobDefinitionModel, diag.Diagnostics) {
	model := kuroko2.JobDefinitionModel{
		Name:               data.Name.ValueString(),
		Description:        data.Description.ValueString(),
		Script:             data.Script.ValueString(),
		NotifyCancellation: data.NotifyCancellation.ValueBool(),
		Suspended:          data.Suspended.ValueBool(),
		SlackChannel:       data.SlackChannel.ValueString(),
	}
	var diags diag.Diagnostics

	diags.Append(data.Admins.ElementsAs(ctx, &model.Admins, false)...)
	if diags.HasError() {
		return model, diags
	}

	diags.Append(data.Cron.ElementsAs(ctx, &model.Cron, false)...)
	if diags.HasError() {
		return model, diags
	}

	diags.Append(data.Tags.ElementsAs(ctx, &model.Tags, false)...)
	if diags.HasError() {
		return model, diags
	}

	switch data.PreventMulti.ValueString() {
	case "NONE":
		model.PreventMulti = 0
	case "WORKING_OR_ERROR":
		model.PreventMulti = 1
	case "WORKING":
		model.PreventMulti = 2
	case "ERROR":
		model.PreventMulti = 3
	default:
		panic(fmt.Sprintf("Unknown prevent_multi value: %s", data.PreventMulti.ValueString()))
	}

	return model, diags
}
