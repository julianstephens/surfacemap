package model

// Resource represents a Terraform resource with its attributes and location in the code
type Resource struct {
	// The ID is the Terraform resource address, which uniquely identifies the resource within the Terraform configuration. It typically follows the format "resource_type.resource_name", e.g. "aws_lambda_function.my_function".
	ID string `json:"id"` // Terraform resource address, e.g. "aws_lambda_function.my_function"
	// The Type is the Terraform resource type, which indicates the kind of resource being defined, e.g. "aws_lambda_function". This is used to determine how to extract and model the resource.
	Type string `json:"type" hcl:",label"` // Resource type, e.g. "aws_lambda_function"
	// Attributes is a map of the raw Terraform arguments for the resource, where the keys are the argument names and the values are their corresponding values. This allows for flexible storage of any attributes defined in the Terraform code.
	Attributes map[string]any `json:"attributes"` // Raw Terraform arguments for the resource
	// Blocks is a map of nested blocks within the resource, where the keys are the block types (e.g. "environment", "vpc_config") and the values are slices of maps representing each block instance. Each block instance is a map of its own attributes. This structure allows for representing complex nested configurations in Terraform.
	Blocks map[string][]map[string]any `json:"blocks"`
	// Location provides the file and line number where the resource is defined in the Terraform code, which is useful for tracing back to the source configuration and for debugging purposes.
	Location FileLocation `json:"location"` // File and line number where the resource is defined in the Terraform code
}

type FileLocation struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

type ResourceType string

const (
	ResourceTypeLambdaFunction  ResourceType = "aws_lambda_function"
	ResourceTypeS3Bucket        ResourceType = "aws_s3_bucket"
	ResourceTypeIAMRole         ResourceType = "aws_iam_role"
	ResourceTypeIAMPolicy       ResourceType = "aws_iam_policy"
	ResourceTypeEC2Instance     ResourceType = "aws_ec2_instance"
	ResourceTypeDynamoDBTable   ResourceType = "aws_dynamodb_table"
	ResourceTypeRDSInstance     ResourceType = "aws_rds_instance"
	ResourceTypeSQSQueue        ResourceType = "aws_sqs_queue"
	ResourceTypeSNSSubscription ResourceType = "aws_sns_subscription"
	ResourceTypeAPIGateway      ResourceType = "aws_api_gateway"
)

type ResourceCollection struct {
	Lambdas          []*LambdaFunction
	S3Buckets        []*S3Bucket
	IAMRoles         []*IAMRole
	IAMPolicies      []*IAMPolicy
	EC2Instances     []*EC2Instance
	DynamoDBTables   []*DynamoDBTable
	RDSInstances     []*RDSInstance
	SQSQueues        []*SQSQueue
	SNSSubscriptions []*SNSSubscription
	APIGateways      []*APIGateway
}

type AuthType string

const (
	AuthTypeNone    AuthType = "none"
	AuthTypeIAM     AuthType = "iam"
	AuthTypeJWT     AuthType = "jwt"
	AuthTypeLambda  AuthType = "lambda"
	AuthTypeCognito AuthType = "cognito"
	AuthTypeAPIKey  AuthType = "api_key"
	AuthTypeCustom  AuthType = "custom"
)
