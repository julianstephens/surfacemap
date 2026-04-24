package model

type IAMRole struct {
	Resource

	RoleName         string
	AssumeRolePolicy *PolicyDocument
	InlinePolicies   []*IAMPolicy
	AttachedPolicies []*IAMPolicy

	AllPermissions []Permission
}

type IAMPolicy struct {
	PolicyName     string
	PolicyDocument *PolicyDocument
	IsInline       bool
	Location       FileLocation
}

type PolicyDocument struct {
	Version   string
	Statement []PolicyStatement
}

type PolicyStatement struct {
	Sid        string
	Effect     string
	Actions    []string
	Resources  []string
	Conditions map[string]any

	HasWildcardAction   bool
	HasWildcardResource bool
}

type Permission struct {
	Service  string
	Action   string
	Resource string
	Effect   string

	SourcePolicy string
	SourceType   PermissionSource
}

type PermissionSource string

const (
	PermissionSourceIAM  PermissionSource = "iam"
	PermissionSourceCode PermissionSource = "code"
)
