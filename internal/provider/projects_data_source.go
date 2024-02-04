// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	client "github.com/andreaswwilson/todoist-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ProjectsDataSource{}
var _ datasource.DataSourceWithConfigure = &ProjectsDataSource{}

func NewProjectsDataSource() datasource.DataSource {
	return &ProjectsDataSource{}
}

// ProjectsDataSource defines the data source implementation.
type ProjectsDataSource struct {
	client *client.Client
}

// ProjectsDataSourceModel describes the data source data model.
type ProjectsDataSourceModel struct {
	Name           types.String `tfsdk:"name"`
	Id             types.String `tfsdk:"id"`
	CommentCount   types.Int64  `tfsdk:"comment_count"`
	Color          types.String `tfsdk:"color"`
	IsShared       types.Bool   `tfsdk:"is_shared"`
	Order          types.Int64  `tfsdk:"order"`
	IsFavorite     types.Bool   `tfsdk:"is_favorite"`
	IsInboxProject types.Bool   `tfsdk:"is_inbox_project"`
	IsTeamInbox    types.Bool   `tfsdk:"is_team_inbox"`
	ViewStyle      types.String `tfsdk:"view_style"`
	URL            types.String `tfsdk:"url"`
	ParentID       types.String `tfsdk:"parent_id"`
}

func (d *ProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *ProjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Projects data source",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Project name",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Project identifier",
				Required:            true,
			},
			"comment_count": schema.Int64Attribute{
				MarkdownDescription: "Number of comments",
				Computed:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color of the project",
				Computed:            true,
			},
			"is_shared": schema.BoolAttribute{
				MarkdownDescription: "Is the project shared",
				Computed:            true,
			},
			"order": schema.Int64Attribute{
				MarkdownDescription: "Order of the project",
				Computed:            true,
			},
			"is_favorite": schema.BoolAttribute{
				MarkdownDescription: "Is the project a favorite",
				Computed:            true,
			},
			"is_inbox_project": schema.BoolAttribute{
				MarkdownDescription: "Is the project an inbox project",
				Computed:            true,
			},
			"is_team_inbox": schema.BoolAttribute{
				MarkdownDescription: "Is the project a team inbox",
				Computed:            true,
			},
			"view_style": schema.StringAttribute{
				MarkdownDescription: "View style of the project",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL of the project",
				Computed:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "Parent ID of the project",
				Computed:            true,
			}},
	}
}

func (d *ProjectsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.client.GetProject(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project, got error: %s", err))
		return
	}
	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Name = types.StringValue(response.Name)
	data.CommentCount = types.Int64Value(int64(response.CommentCount))
	data.Color = types.StringValue(response.Color)
	data.IsShared = types.BoolValue(response.IsShared)
	data.Order = types.Int64Value(int64(response.Order))
	data.IsFavorite = types.BoolValue(response.IsFavorite)
	data.IsInboxProject = types.BoolValue(response.IsInboxProject)
	data.IsTeamInbox = types.BoolValue(response.IsTeamInbox)
	data.ViewStyle = types.StringValue(response.ViewStyle)
	data.URL = types.StringValue(response.URL)
	data.ParentID = types.StringValue(response.ParentID)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read project data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
