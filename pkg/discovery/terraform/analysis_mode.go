package terraform

// AnalysisMode defines how to analyze a module
type AnalysisMode int

const (
	// AnalysisModeAuto automatically determines the best mode
	AnalysisModeAuto AnalysisMode = iota

	// AnalysisModeStructure analyzes module structure without specific variables
	AnalysisModeStructure

	// AnalysisModeSingleCaller analyzes with one caller's variables
	AnalysisModeSingleCaller

	// AnalysisModeAllCallers analyzes all instantiations separately
	AnalysisModeAllCallers

	// AnalysisModeHypothetical analyzes with user-provided variables
	AnalysisModeHypothetical

	// AnalysisModeFollowCalls follows module calls from entry (original behavior)
	AnalysisModeFollowCalls
)

func (m AnalysisMode) String() string {
	return ""
}

// AnalysisRequest contains all parameters for an analysis
type AnalysisRequest struct {
	// EntryPath is the directory to analyze
	EntryPath string

	// Mode specifies how to analyze (Auto determines automatically)
	Mode AnalysisMode

	// CallerPath is the specific caller to use (for SingleCaller mode)
	CallerPath string

	// Variables are explicitly provided variable values (for Hypothetical mode)
	Variables map[string]any

	// Config contains analysis configuration
	Config *AnalysisConfig
}

// AnalysisResult contains the results of an analysis
type AnalysisResult struct {
	// Mode is the actual mode used (if auto-detected)
	Mode AnalysisMode

	// Graph is the complete module graph
	Graph *ModuleGraph

	// Config is the flattened Terraform configuration
	Config *TerraformConfig

	// Metadata about the analysis
	Metadata *AnalysisMetadata

	// Suggestions for alternative analyses
	Suggestions []string
}

// AnalysisMetadata provides context about the analysis
type AnalysisMetadata struct {
	// EntryPath is the analyzed directory
	EntryPath string

	// EntryType is the detected module type
	EntryType ModuleType

	// ModuleCount is the number of modules analyzed
	ModuleCount int

	// CallersFound is the number of callers discovered
	CallersFound int

	// VariableSource indicates where variables came from
	VariableSource VariableSource

	// VariableSourcePath is the path to the variable source
	VariableSourcePath string
}

// VariableSource indicates how variables were resolved
type VariableSource int

const (
	VariableSourceNone VariableSource = iota
	VariableSourceModuleCall
	VariableSourceTfVars
	VariableSourceCLI
	VariableSourceDefaults
)

func (v VariableSource) String() string {
	return ""
}
