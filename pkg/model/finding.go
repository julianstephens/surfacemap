package model

type Finding struct {
	ID             string
	Type           FindingType
	Severity       Severity
	Title          string
	Description    string
	Recommendation string

	ResourceID   string
	ResourceType ResourceType

	Location FileLocation
	Evidence map[string]any

	Category string
	Tags     []string
}

type FindingType string

const (
	FindingPublicNoAuth      FindingType = "public_no_auth"
	FindingOverPermissioned  FindingType = "over_permissioned"
	FindingUnderPermissioned FindingType = "under_permissioned"
	FindingWildcardIAM       FindingType = "wildcard_iam"
	FindingHardcodedSecret   FindingType = "hardcoded_secret"
	FindingSQLInjection      FindingType = "sql_injection"
	FindingDangerousFunction FindingType = "dangerous_function"
)

type FindingRegistry struct {
	findings   []*Finding
	byType     map[FindingType][]*Finding
	bySeverity map[Severity][]*Finding
	byResource map[string][]*Finding
}

func (r *FindingRegistry) Add(finding *Finding)                       {}
func (r *FindingRegistry) GetAll() []*Finding                         { return nil }
func (r *FindingRegistry) GetByType(t FindingType) []*Finding         { return nil }
func (r *FindingRegistry) GetBySeverity(s Severity) []*Finding        { return nil }
func (r *FindingRegistry) GetByResource(resourceID string) []*Finding { return nil }
func (r *FindingRegistry) Count() int                                 { return 0 }
func (r *FindingRegistry) CountBySeverity(s Severity) int             { return 0 }
