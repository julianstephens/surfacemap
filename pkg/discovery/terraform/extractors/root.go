package extractors

import (
	"github.com/julianstephens/surfacemap/pkg/model"
)

type Extractor[T any] interface {
	Extract(resource model.Resource) (T, error)
}

// ExtractResource is a helper function that takes a generic model.Resource and determines which specific extractor to use based on the resource type. It returns the extracted resource as an interface{} and any error encountered during extraction. This function serves as a central point for dispatching to the appropriate extractor for different Terraform resource types.
// For example, if the resource type is "aws_lambda_function", it will call ExtractLambdaFunction to extract the Lambda function details. If the resource type is not recognized, it returns nil without an error, allowing the caller to handle unsupported resource types gracefully.
func ExtractResource(resource model.Resource) (any, error) {
	switch resource.Type {
	case "aws_lambda_function":
		return ExtractLambdaFunction(resource)
	case "aws_dynamodb_table":
		return ExtractDynamoDBTable(resource)
	default:
		return nil, nil
	}
}
