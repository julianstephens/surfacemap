package model

type EC2Instance struct {
	Resource

	InstanceID   string
	InstanceType string
	State        string
	PublicIP     *string
	PrivateIP    *string
	VPCID        *string
	SubnetID     *string

	Exposure ExposureLevel
}
