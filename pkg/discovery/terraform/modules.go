package terraform

import "github.com/julianstephens/surfacemap/pkg/model"

type Module struct {
	Name      string
	Source    string
	Variables map[string]any
	Resources []model.Resource
	Modules   []Module
	Locals    map[string]any
}

type ModuleConfig struct {
	Resources []model.Resource
	Modules   []Module
	Variables map[string]any
	Locals    map[string]any
}
