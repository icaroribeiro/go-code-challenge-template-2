package resolver_test

var getHealthCheckQuery = `query {
	getHealthCheck { 
		status
	}
}`

type GetHealthCheckResponse struct {
	GetHealthCheck struct {
		Status string
	}
}
