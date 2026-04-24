package model

type APIGateway struct {
	Resource

	APIID       string
	Name        string
	Routes      []*APIRoute
	Authorizers []*Authorizer

	Exposure ExposureLevel
}

type APIRoute struct {
	RouteID       string
	RouteKey      string
	Method        string
	Path          string
	IntegrationID string
	Integration   *LambdaFunction
	AuthorizerID  *string
	Authorizer    *Authorizer
	Location      FileLocation
}

type Authorizer struct {
	AuthorizerID string
	Name         string
	Type         AuthType
	Config       map[string]any
}
