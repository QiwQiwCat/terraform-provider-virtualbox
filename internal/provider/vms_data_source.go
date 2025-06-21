package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
)

var (
	_ datasource.DataSource = &vmsDataSource{}
)

type vmsDataSource struct{}

type vmsDataSourceModel struct {
	Vms []vmsModel `tfsdk:"vms"`
}

type vmsModel struct {
	UUID types.String `tfsdk:"uuid"`
	Name types.String `tfsdk:"name"`
}

func (d *vmsDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vms"
}

func (d *vmsDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Fetches the list of Virtual Machines registered with your Oracle VM VirtualBox installation.",
		Attributes: map[string]schema.Attribute{
			"vms": schema.ListNestedAttribute{
				Description: "List of Virtual Machines.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "The UUID of the Virtual Machine.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the Virtual Machine.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *vmsDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var state vmsDataSourceModel

	output, err := RunVBoxManageWithOutput("list", "vms")
	if err != nil {
		response.Diagnostics.AddError(
			"Error Reading Virtual Machines",
			"An error occurred while trying to read the list of Virtual Machines: "+err.Error(),
		)
		return
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var vmList []vmsModel

	re := regexp.MustCompile(`^"(.*?)"\s+\{([a-fA-F0-9\-]+)\}$`)

	for _, line := range lines {
		match := re.FindStringSubmatch(line)
		if len(match) == 3 {
			name := match[1]
			uuid := match[2]
			vmList = append(vmList, vmsModel{
				Name: types.StringValue(name),
				UUID: types.StringValue(uuid),
			})
		}
	}

	// Map response body to model
	state.Vms = vmList

	// Set state
	diags := response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func NewVmsDataSource() datasource.DataSource {
	return &vmsDataSource{}
}
