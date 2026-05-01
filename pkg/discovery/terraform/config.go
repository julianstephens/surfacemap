package terraform

import "time"

// AnalysisConfig controls analysis behavior
type AnalysisConfig struct {
	// CallerDiscovery controls caller search behavior
	CallerDiscovery CallerDiscoveryConfig

	// ModuleResolution controls how modules are resolved
	ModuleResolution ModuleResolutionConfig

	// VariableResolution controls variable resolution
	VariableResolution VariableResolutionConfig

	// Performance tuning
	MaxConcurrentParsers int
	Timeout              time.Duration
}

// CallerDiscoveryConfig controls how callers are discovered
type CallerDiscoveryConfig struct {
	// Enabled determines if caller discovery runs
	Enabled bool

	// SearchScope defines where to look for callers
	// Relative to entry point (e.g., "..", "../..")
	SearchScope []string

	// MaxDepth limits directory traversal depth
	MaxDepth int

	// Patterns are glob patterns for caller directories
	// Examples: "build/*", "environments/*", "*/terraform"
	Patterns []string
}

// ModuleResolutionConfig controls module resolution
type ModuleResolutionConfig struct {
	// AllowRemote enables remote module resolution (registry, git)
	AllowRemote bool

	// CachePath is where remote modules are cached
	CachePath string

	// FollowSymlinks determines if symlinks are followed
	FollowSymlinks bool

	// MaxDepth limits module nesting depth (cycle protection)
	MaxDepth int
}

// VariableResolutionConfig controls variable resolution
type VariableResolutionConfig struct {
	// AllowUnresolved allows analysis with unresolved variables
	AllowUnresolved bool

	// TfVarsFiles are paths to .tfvars files to load
	TfVarsFiles []string

	// EnvironmentVariables enables TF_VAR_* resolution
	EnvironmentVariables bool

	// Priority determines order: CLI > ModuleCall > TfVars > Env > Defaults
	// This is fixed, but can be documented here
}

// DefaultAnalysisConfig returns sensible defaults
func DefaultAnalysisConfig() *AnalysisConfig {
	return &AnalysisConfig{}
}
