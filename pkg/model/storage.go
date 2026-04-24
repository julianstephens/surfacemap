package model

type S3Bucket struct {
	Resource

	BucketName string
	Region     string

	Exposure ExposureLevel
}

type DynamoDBTable struct {
	Resource

	TableName string
	Region    string

	Exposure ExposureLevel
}

type RDSInstance struct {
	Resource

	DBInstanceIdentifier string
	DBName               string
	Engine               string
	Region               string
	AvailabilityZone     string
	DBSubnetGroup        string
	PubliclyAccessible   bool
	VPCSecurityGroups    []string

	Exposure ExposureLevel
}
