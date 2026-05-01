package terraform

// CallerSearchResult contains results from searching for callers
type CallerSearchResult struct {
	// Callers found that reference the target module
	Callers []*CallerInfo

	// SearchedPaths are directories that were checked
	SearchedPaths []string

	// Errors encountered during search (non-fatal)
	Errors []error
}

// CallerMatch represents a potential caller match
type CallerMatch struct {
	// Path to the directory with module call
	Path string

	// ModuleCalls in this directory
	ModuleCalls []*ModuleCall

	// Score indicates match confidence (0.0 - 1.0)
	Score float64
}
