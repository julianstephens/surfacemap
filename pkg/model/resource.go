package model

type Resource struct {
	ID         string                      `json:"id"`                // Terraform resource address, e.g. "aws_lambda_function.my_function"
	Type       string                      `json:"type" hcl:",label"` // Resource type, e.g. "aws_lambda_function"
	Attributes map[string]any              `json:"attributes"`        // Raw Terraform arguments for the resource
	Blocks     map[string][]map[string]any `json:"blocks"`
	Location   FileLocation                `json:"location"` // File and line number where the resource is defined in the Terraform code
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
