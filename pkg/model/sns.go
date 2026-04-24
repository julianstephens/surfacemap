package model

type SNSSubscription struct {
	Resource

	SubscriptionARN string
	Protocol        SNSProtocol
	Endpoint        string
	TopicARN        string

	Exposure ExposureLevel
}

type SNSProtocol string

const (
	SNSProtocolApplication SNSProtocol = "application"
	SNSProtocolFirehose    SNSProtocol = "firehose"
	SNSProtocolLambda      SNSProtocol = "lambda"
	SNSProtocolSQS         SNSProtocol = "sqs"
	SNSProtocolSMS         SNSProtocol = "sms"
	SNSProtocolEmail       SNSProtocol = "email"
	SNSProtocolEmailJSON   SNSProtocol = "email_json"
	SNSProtocolHTTP        SNSProtocol = "http"
	SNSProtocolHTTPS       SNSProtocol = "https"
)
