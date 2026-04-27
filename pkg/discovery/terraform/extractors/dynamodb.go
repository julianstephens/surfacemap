package extractors

import (
	"errors"

	"github.com/julianstephens/surfacemap/pkg/model"
)

var (
	ErrMissingTableName = errors.New("missing name attribute in DynamoDB resource")
	ErrMissingRegion    = errors.New("missing region attribute in DynamoDB resource")
)

type DynamoDBExtractor struct {
}

func (de *DynamoDBExtractor) Extract(resource model.Resource) (*model.DynamoDBTable, error) {
	if resource.Type != "aws_dynamodb_table" {
		return nil, ErrUnsupportedResourceType(resource.Type)
	}

	tableName, ok := resource.Attributes["name"].(string)
	if !ok {
		return nil, ErrMissingTableName
	}

	region, ok := resource.Attributes["region"].(string)
	if !ok {
		return nil, ErrMissingRegion
	}

	return &model.DynamoDBTable{
		Resource:  resource,
		TableName: tableName,
		Region:    region,
	}, nil
}

func ExtractDynamoDBTable(resource model.Resource) (*model.DynamoDBTable, error) {
	extractor := &DynamoDBExtractor{}
	return extractor.Extract(resource)
}
