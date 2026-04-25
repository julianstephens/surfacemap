package terraform

import (
	"context"

	"github.com/julianstephens/surfacemap/pkg/model"
)

type TerraformParser interface {
	// Parse loads and parses all .tf files in a directory
	Parse(ctx context.Context, rootPath string) (*TerraformConfig, error)

	// LoadVariables loads variable values from tfvars files
	LoadVariables(varFiles ...string) error

	// ResolveModule inlines a local module into the main config
	ResolveModule(modulePath string) (*ModuleConfig, error)
}

type TerraformConfig struct {
	Resources []model.Resource
	Modules   []Module
	Variables map[string]any
	Locals    map[string]any
}

type HCLParser struct{}

func NewHCLParser() *HCLParser {
	return &HCLParser{}
}

func (p *HCLParser) Parse(ctx context.Context, rootPath string) (*TerraformConfig, error) {
	return nil, nil
}

func (p *HCLParser) LoadVariables(varFiles ...string) error {
	return nil
}

func (p *HCLParser) ResolveModule(modulePath string) (*ModuleConfig, error) {
	return nil, nil
}
