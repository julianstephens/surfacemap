package terraform

import (
	"context"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/julianstephens/go-utils/generic"
	"github.com/julianstephens/go-utils/logger"

	pkgerrors "github.com/julianstephens/surfacemap/pkg/errors"
	"github.com/julianstephens/surfacemap/pkg/model"
)

// GraphBuilder implements ModuleGraphBuilder
type GraphBuilder struct {
	parser   *HCLParser
	logger   *logger.Logger
	resolver ModuleResolver
	visited  map[string]bool
	config   *AnalysisConfig
}

// NewGraphBuilder creates a new graph builder
func NewGraphBuilder(parser *HCLParser, logger *logger.Logger, config *AnalysisConfig) *GraphBuilder {
	return &GraphBuilder{
		parser:   parser,
		logger:   logger.WithField("component", "graph_builder"),
		resolver: NewPathResolver(logger),
		visited:  make(map[string]bool),
		config:   config,
	}
}

// Build constructs a module graph starting from entry point
// Parameters:
//
//	ctx: Context for cancellation
//	entryPath: Starting directory
//	config: Analysis configuration
//
// Returns:
//
//	*ModuleGraph: Complete graph
//	error: Any errors
func (b *GraphBuilder) Build(ctx context.Context, entryPath string, config *AnalysisConfig) (*ModuleGraph, error) {
	b.logger.Infof("building module graph from entry point: %s", entryPath)
	moduleGraph := ModuleGraph{
		Modules: make(map[string]*ModuleNode),
	}

	resolvedEntryPath, err := b.resolver.Resolve(entryPath, entryPath)
	if err != nil {
		return nil, pkgerrors.NewModuleError(
			pkgerrors.ErrModuleResolve,
			"resolve entry path",
			entryPath,
			entryPath,
			err,
		)
	}
	b.logger.Debugf("resolved entry path: %s", resolvedEntryPath)

	rootNode, err := b.buildNode(ctx, resolvedEntryPath, "root", nil, 0)
	if err != nil {
		return nil, pkgerrors.WrapWithContext(err, "failed to build root node from %s", resolvedEntryPath)
	}
	moduleGraph.Root = rootNode
	b.logger.Debugf("built root node: %s with %d resources and %d module calls",
		rootNode.Path, len(rootNode.Resources), len(rootNode.Calls))

	for _, modCall := range rootNode.Calls {
		modNode, err := b.buildNode(ctx, modCall.ResolvedPath, modCall.Name, modCall.Module.VariableValues, 1)
		if err != nil {
			continue
		}
		moduleGraph.Modules[modNode.Path] = modNode
		b.logger.Debugf("built child node: %s with %d resources and %d module calls",
			modNode.Path, len(modNode.Resources), len(modNode.Calls))
	}
	b.logger.Infof("completed building graph: %d modules, %d resources total",
		len(moduleGraph.Modules)+1, rootNode.ResourceCount)

	moduleGraph.Callers = b.buildReverseIndex(&moduleGraph)
	return &moduleGraph, nil
}

// buildNode recursively constructs a module node and its children
// Parameters:
//
//	ctx: Context
//	modulePath: Absolute path to module
//	name: Module name (from call or "root")
//	variables: Variable values from parent
//	depth: Current recursion depth
//
// Returns:
//
//	*ModuleNode: Constructed node
//	error: Any errors
func (b *GraphBuilder) buildNode(ctx context.Context, modulePath, name string, variables map[string]any, depth int) (*ModuleNode, error) {
	if generic.HasKey(b.visited, modulePath) {
		b.logger.Warnf("cycle detected at %s, skipping", modulePath)
		return nil, nil
	}

	b.logger.Debugf("parsing module at %s (depth %d)", modulePath, depth)
	b.visited[modulePath] = true

	if depth > b.config.CallerDiscovery.MaxDepth {
		b.logger.Warnf("max depth exceeded at %s, skipping", modulePath)
		return nil, nil
	}

	tfConfig, err := b.parser.Parse(ctx, modulePath)
	if err != nil {
		return nil, pkgerrors.WrapWithContext(err, "failed to parse module at %s", modulePath)
	}
	b.logger.Debugf("parsed module at %s: %d resources, %d data sources, %d module calls",
		modulePath, len(tfConfig.Resources), len(tfConfig.DataSources), len(tfConfig.Modules))

	node := &ModuleNode{
		Name:   name,
		Path:   modulePath,
		Source: modulePath,
		Type:   ModuleTypeUnknown,

		Calls:    make([]*ModuleCall, 0),
		Children: make([]*ModuleNode, 0),

		Resources:      tfConfig.Resources,
		DataSources:    tfConfig.DataSources,
		Locals:         tfConfig.Locals,
		VariableValues: generic.MergeMap(tfConfig.Variables, variables),
		ResourceCount:  tfConfig.ResourceCount,
		FileCount:      tfConfig.FileCount,
	}

	for _, module := range tfConfig.Modules {
		resolvedSource, err := b.resolver.Resolve(modulePath, module.Source)
		if err != nil {
			b.logger.Warnf("failed to resolve module source %s in %s: %v", module.Source, modulePath, err)
			continue
		}

		if !b.resolver.IsLocal(resolvedSource) {
			b.logger.Warnf("non-local module source %s in %s, skipping", module.Source, modulePath)
			continue
		}

		child, err := b.buildNode(ctx, resolvedSource, module.Name, module.Variables, depth+1)
		if err != nil {
			continue
		}

		inputs, err := b.parser.ExtractModuleInputs(modulePath, module.Name)
		if err != nil {
			b.logger.Warnf("failed to extract module inputs for %s in %s: %v", module.Name, modulePath, err)
			continue
		}

		childCall := &ModuleCall{
			Name:         module.Name,
			Source:       module.Source,
			Inputs:       inputs,
			ResolvedPath: resolvedSource,
			Module:       child,
			Location: model.FileLocation{
				File: module.Source,
				Line: 0, // TODO: Get actual line number from HCL
			},
		}

		if err := b.LinkModules(node, childCall, child); err != nil {
			continue
		}
	}

	return node, nil
}

// BuildFromCaller builds a graph with caller as root, linking to target
//
// Algorithm:
// 1. Parse caller directory
// 2. Find module call to target
// 3. Extract variables from module call
// 4. Parse target with those variables
// 5. Build graph: caller (root) → target (child)
//
// Parameters:
//
//	ctx: Context
//	callerPath: Path to caller
//	targetPath: Path to target module
//
// Returns:
//
//	*ModuleGraph: Two-node graph (caller → target)
//	error: Any errors
func (b *GraphBuilder) BuildFromCaller(ctx context.Context, callerPath, targetPath string) (*ModuleGraph, error) {
	// TODO: Implement BuildFromCaller
	return nil, pkgerrors.NewModuleError(
		pkgerrors.ErrNotImplemented,
		"build from caller",
		callerPath,
		targetPath,
		nil,
	)
}

// LinkModules establishes parent-child relationship
// Parameters:
//
//	parent: Parent module node
//	call: Module call information
//	child: Child module node
//
// Returns:
//
//	error: Any errors
func (b *GraphBuilder) LinkModules(parent *ModuleNode, call *ModuleCall, child *ModuleNode) error {
	child.Parent = parent
	parent.Children = append(parent.Children, child)
	call.Module = child
	call.ResolvedPath = child.Path
	return nil
}

// extractModuleCalls gets module call information from a directory
//
// Algorithm:
//  1. Use tfconfig.LoadModule to get metadata
//  2. For each module call:
//     a. Extract source
//     b. Extract input variables (names and expressions)
//     c. Store location information
//  3. Return list of ModuleCall structs
//
// Parameters:
//
//	modulePath: Path to module directory
//
// Returns:
//
//	[]*ModuleCall: List of module calls
//	error: Any errors
func (b *GraphBuilder) extractModuleCalls(modulePath string) ([]*ModuleCall, error) {
	module, diags := tfconfig.LoadModule(modulePath)
	if diags.HasErrors() {
		return nil, diags.Err()
	}
	calls := make([]*ModuleCall, 0)

	for _, call := range module.ModuleCalls {
		calls = append(calls, &ModuleCall{
			Name:   call.Name,
			Source: call.Source,
			Location: model.FileLocation{
				File: call.Pos.Filename,
				Line: call.Pos.Line,
			},
		})
	}
	return calls, nil
}

// buildReverseIndex creates callers map from graph
//
// Algorithm:
//  1. Traverse entire graph
//  2. For each node with children:
//     a. For each child:
//     - Add parent to child's callers list
//     - Store module call information
//  3. Return callers map
//
// Parameters:
//
//	graph: Module graph
//
// Returns:
//
//	map[string][]*CallerInfo: Callers by module path
func (b *GraphBuilder) buildReverseIndex(graph *ModuleGraph) map[string][]*CallerInfo {
	// TODO: Implement buildReverseIndex
	return make(map[string][]*CallerInfo)
}
