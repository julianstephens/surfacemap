package model

type SQSQueue struct {
	Resource

	QueueName string
	Region    string

	Exposure ExposureLevel
}
