package terraform

import (
	"github.com/julianstephens/surfacemap/pkg/model"
)

// ModuleGraph represents the complete module dependency graph
type ModuleGraph struct {
	// Root is the entry point module
	Root *ModuleNode

	// Modules maps absolute paths to module nodes
	Modules map[string]*ModuleNode

	// Callers maps module paths to their callers (reverse lookup)
	Callers map[string][]*CallerInfo
}

// ModuleNode represents a single module in the dependency graph
type ModuleNode struct {
	// Path is the absolute filesystem path to this module
	Path string

	// Name is the module identifier (from module block or "root")
	Name string

	// Source is the original source string (e.g., "../.." or "registry/...")
	Source string

	// Type indicates what kind of module this is
	Type ModuleType

	// Parsed content from .tf files
	Resources   []model.Resource
	DataSources []model.Resource
	Locals      map[string]any
	Variables   map[string]*VariableDefinition

	// Variable values (if provided by caller or flags)
	VariableValues map[string]any

	// Module relationships
	Parent   *ModuleNode
	Children []*ModuleNode
	Calls    []*ModuleCall

	// Metadata
	FileCount     int
	ResourceCount int
}

// ModuleType classifies the module's purpose
type ModuleType int

const (
	ModuleTypeUnknown ModuleType = iota
	// ModuleTypeCaller has module calls but few/no resources
	ModuleTypeCaller
	// ModuleTypeReusable has resources and is called by others
	ModuleTypeReusable
	// ModuleTypeStandalone has resources but no callers or calls
	ModuleTypeStandalone
	// ModuleTypeHybrid has both resources and module calls
	ModuleTypeHybrid
)

func (m ModuleType) String() string {
	return ""
}

// ModuleCall represents a module instantiation
type ModuleCall struct {
	// Name is the module identifier in the module block
	Name string

	// Source is the module source path
	Source string

	// Inputs are the variable values passed to the module
	Inputs map[string]*InputValue

	// ResolvedPath is the absolute path after resolution
	ResolvedPath string

	// Module is the resolved module node (if local)
	Module *ModuleNode

	// Location in source file
	Location model.FileLocation
}

// InputValue represents a value passed to a module
type InputValue struct {
	// Name is the variable name
	Name string

	// Value is the resolved value (may be nil if not resolvable)
	Value any

	// Expression is the raw HCL expression
	Expression string

	// IsReference indicates if this references another resource
	IsReference bool
}

// VariableDefinition represents a variable declaration
type VariableDefinition struct {
	// Name is the variable identifier
	Name string

	// Type is the variable type constraint (string, bool, object, etc.)
	Type string

	// Default is the default value (may be nil)
	Default any

	// Required indicates if the variable must be provided
	Required bool

	// Description from variable block
	Description string

	// Location in source file
	Location model.FileLocation
}

// CallerInfo describes a module that calls another module
type CallerInfo struct {
	// CallerPath is the absolute path to the calling module
	CallerPath string

	// CallName is the module block name
	CallName string

	// Variables are the inputs provided by this caller
	Variables map[string]any

	// ModuleCall is the full module call information
	ModuleCall *ModuleCall
}
