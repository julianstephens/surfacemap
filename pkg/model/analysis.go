package model

import "time"

type ExposureLevel string

const (
	ExposureLevelUnknown  ExposureLevel = "unknown"
	ExposureLevelPublic   ExposureLevel = "public"
	ExposureLevelInternal ExposureLevel = "internal"
)

type CodeAnalysis struct {
	Language       string
	SourceFiles    []string
	AWSCalls       []AWSCall
	SecurityIssues []SecurityIssue
	Dependencies   []string
}

type AWSCall struct {
	Service    string
	Operation  string
	Resource   string
	Parameters map[string]any
	Location   FileLocation
}

type SecurityIssue struct {
	Type           IssueType
	Severity       Severity
	Title          string
	Description    string
	Recommendation string

	Location FileLocation

	CodeSnippet string
	LineContent string
}

type IssueType string

const (
	IssueHardcodedCredential IssueType = "hardcoded_credential"
	IssueSQLInjection        IssueType = "sql_injection"
	IssueLogsSensitiveData   IssueType = "logs_sensitive_data"
	IssueDangerousFunction   IssueType = "dangerous_function"
	IssueInsecureDeserialize IssueType = "insecure_deserialize"
)

type Severity int

const (
	SeverityLow Severity = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

func (s Severity) String() string {
	return []string{"Low", "Medium", "High", "Critical"}[s]
}

type AnalysisResult struct {
	RootPath   string
	AnalyzedAt time.Time
	Duration   time.Duration

	Resources *ResourceCollection

	Findings    *FindingRegistry
	IAMAnalysis []*IAMAnalysis
	DataFlows   []*DataFlow

	Stats *AnalysisStats
}

type IAMAnalysis struct {
	Lambda *LambdaFunction
	Role   *IAMRole

	GrantedPermissions []Permission
	UsedPermissions    []Permission
	OverPermissioned   []Permission
	UnderPermissioned  []Permission

	MinimalPolicy *PolicyDocument
}

type DataFlow struct {
	EntryPoint string
	Path       []FlowHop
	DataStores []string

	HasAuthentication bool
	HasEncryption     bool
}

type FlowHop struct {
	ResourceID   string
	ResourceType ResourceType
	Action       string
	Permissions  []Permission
}

type AnalysisStats struct {
	TotalResources  int
	ResourcesByType map[ResourceType]int

	PublicEndpoints  int
	InternalServices int

	FindingsCount      int
	FindingsBySeverity map[Severity]int
	FindingsByType     map[FindingType]int

	IAMRolesAnalyzed  int
	OverPermissioned  int
	UnderPermissioned int
}
