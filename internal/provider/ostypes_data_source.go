package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ datasource.DataSource = &ostypesDataSource{}
)

type ostypesDataSource struct{}

type ostypesDataSourceModel struct {
	OsTypes []ostypesModel `tfsdk:"ostypes"`
}

type ostypesModel struct {
	ID           types.String `tfsdk:"id"`
	Description  types.String `tfsdk:"description"`
	FamilyId     types.String `tfsdk:"family_id"`
	FamilyDesc   types.String `tfsdk:"family_desc"`
	Architecture types.String `tfsdk:"architecture"`
	Is64Bit      types.Bool   `tfsdk:"is_64_bit"`
}

func (d *ostypesDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ostypes"
}

func (d *ostypesDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Fetches the list of available OS types supported by your Oracle VM VirtualBox installation.",
		Attributes: map[string]schema.Attribute{
			"ostypes": schema.ListNestedAttribute{
				Description: "List of available OS types.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the OS type.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "A human-readable description of the OS type.",
							Computed:    true,
						},
						"family_id": schema.StringAttribute{
							Description: "The ID of the OS family to which this OS type belongs.",
							Computed:    true,
						},
						"family_desc": schema.StringAttribute{
							Description: "A human-readable description of the OS family.",
							Computed:    true,
						},
						"architecture": schema.StringAttribute{
							Description: "The architecture of the OS type (e.g., x86, x86_64).",
							Computed:    true,
						},
						"is_64_bit": schema.BoolAttribute{
							Description: "Indicates whether the OS type is 64-bit.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ostypesDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var state ostypesDataSourceModel

	output, err := RunVBoxManageWithOutput("list", "ostypes", "-l", "-s")
	if err != nil {
		response.Diagnostics.AddError(
			"Failed to fetch OS types",
			"An error occurred while trying to list available OS types: "+err.Error(),
		)
		return
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var models []ostypesModel

	var current ostypesModel
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			// Push current if filled
			if !current.ID.IsNull() {
				models = append(models, current)
				current = ostypesModel{}
			}
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "ID":
			current.ID = types.StringValue(value)
		case "Description":
			current.Description = types.StringValue(value)
		case "Family ID":
			current.FamilyId = types.StringValue(value)
		case "Family Desc":
			current.FamilyDesc = types.StringValue(value)
		case "Architecture":
			current.Architecture = types.StringValue(value)
		case "64 bit":
			current.Is64Bit = types.BoolValue(value == "true")
		}
	}

	// Add the last block if necessary
	if !current.ID.IsNull() {
		models = append(models, current)
	}

	// Map response body to model
	state.OsTypes = models

	// Set state
	diags := response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func NewOsTypesDataSource() datasource.DataSource {
	return &ostypesDataSource{}
}
