// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"os/exec"
	"strings"
)

// Ensure virtualboxProvider satisfies various provider interfaces.
var _ provider.Provider = &virtualboxProvider{}

// virtualboxProvider defines the provider implementation.
type virtualboxProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// virtualboxProviderModel describes the provider data model.
type virtualboxProviderModel struct{}

func (p *virtualboxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "virtualbox"
	resp.Version = p.version
}

func (p *virtualboxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *virtualboxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *virtualboxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVmResource,
	}
}

func (p *virtualboxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVmsDataSource,
		NewOsTypesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &virtualboxProvider{
			version: version,
		}
	}
}

func RunVBoxManage(args ...string) error {
	cmd := exec.Command("VBoxManage", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func RunVBoxManageWithOutput(args ...string) (string, error) {
	cmd := exec.Command("VBoxManage", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func ParseShowVMInfo(output string) map[string]string {
	lines := strings.Split(output, "\n")
	kv := map[string]string{}
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
			kv[key] = val
		}
	}
	return kv
}

func GetVMInfoFromOutput(output map[string]string, key string) string {
	if val, ok := output[key]; ok {
		return val
	}
	return ""
}
