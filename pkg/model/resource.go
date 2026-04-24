package model

type Resource struct {
	ID         string         `json:"id"`         // Terraform resource address, e.g. "aws_lambda_function.my_function"
	Type       string         `json:"type"`       // Resource type, e.g. "aws_lambda_function"
	Attributes map[string]any `json:"attributes"` // Raw Terraform attributes for the resource
	Location   FileLocation   `json:"location"`   // File and line number where the resource is defined in the Terraform code
}

type FileLocation struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

type ResourceType string

const (
	ResourceTypeLambdaFunction ResourceType = "aws_lambda_function"
	ResourceTypeS3Bucket       ResourceType = "aws_s3_bucket"
	ResourceTypeIAMRole        ResourceType = "aws_iam_role"
	ResourceTypeIAMPolicy      ResourceType = "aws_iam_policy"
	ResourceTypeEC2Instance    ResourceType = "aws_ec2_instance"
	ResourceTypeDynamoDBTable  ResourceType = "aws_dynamodb_table"
	ResourceTypeRDSInstance    ResourceType = "aws_rds_instance"
	ResourceTypeSQSQueue       ResourceType = "aws_sqs_queue"
	ResourceTypeSNSResource    ResourceType = "aws_sns_resource"
	ResourceTypeAPIGateway     ResourceType = "aws_api_gateway"
)
