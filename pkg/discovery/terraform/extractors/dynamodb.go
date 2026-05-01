package extractors

import (
	"github.com/julianstephens/surfacemap/pkg/errors"
	"github.com/julianstephens/surfacemap/pkg/model"
)

type DynamoDBExtractor struct {
}

func (de *DynamoDBExtractor) Extract(resource model.Resource) (*model.DynamoDBTable, error) {
	if resource.Type != "aws_dynamodb_table" {
		return nil, errors.NewExtractionError(
			errors.ErrUnsupportedResourceType,
			resource.ID,
			resource.Type,
			resource.Location,
			"expected aws_dynamodb_table",
			nil,
		)
	}

	tableName, ok := resource.Attributes["name"].(string)
	if !ok {
		return nil, errors.NewAttributeErrorWithType(
			errors.ErrMissingAttribute,
			resource.ID,
			"name",
			"string",
			resource.Attributes["name"],
			resource.Location,
		)
	}

	region, ok := resource.Attributes["region"].(string)
	if !ok {
		return nil, errors.NewAttributeErrorWithType(
			errors.ErrMissingAttribute,
			resource.ID,
			"region",
			"string",
			resource.Attributes["region"],
			resource.Location,
		)
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
