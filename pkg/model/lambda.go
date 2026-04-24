package model

type LambdaFunction struct {
	Resource

	FunctionName  string    `json:"function_name"` // The name of the Lambda function
	RoleArn       string    `json:"role_arn"`      // The ARN of the IAM role that the Lambda function assumes when it executes
	Architectures *[]string `json:"architectures"` // The instruction set architecture that the function supports, e.g. "x86_64" or "arm64"
	Environment   *struct {
		Variables map[string]string `json:"variables"` // Environment variables for the Lambda function
	} `json:"environment,omitempty"`
	Handler        *string    `json:"handler"`              // The function within your code that Lambda calls to begin execution
	Runtime        *string    `json:"runtime"`              // The runtime environment for the Lambda function, e.g. "nodejs14.x", "python3.8", etc.
	SourceCodePath *string    `json:"source_code_path"`     // The path to the source code for the Lambda function, if available
	VPCConfig      *VPCConfig `json:"vpc_config,omitempty"` // VPC configuration for the Lambda function, if it is associated with a VPC

	ExposureLevel ExposureLevel `json:"exposure_level"` // The exposure level of the Lambda function, e.g. "public", "internal", or "unknown"
}

type VPCConfig struct {
	SubnetIds        *[]string `json:"subnet_ids"`         // The IDs of the subnets associated with the Lambda function
	SecurityGroupIds *[]string `json:"security_group_ids"` // The IDs of the security groups associated with the Lambda function
}

type LambdaTrigger struct {
	Type     TriggerType
	SourceID string
	Public   bool
	AuthType AuthType
}

type TriggerType string

const (
	TriggerTypeAPIGateway      TriggerType = "api_gateway"
	TriggerTypeALB             TriggerType = "alb"
	TriggerTypeS3Event         TriggerType = "s3"
	TriggerTypeCloudWatchEvent TriggerType = "cloudwatch"
	TriggerTypeDynamoDBStream  TriggerType = "dynamodb_stream"
	TriggerTypeSQS             TriggerType = "sqs"
	TriggerTypeSNS             TriggerType = "sns"
	TriggerTypeEventBridge     TriggerType = "eventbridge"
)
