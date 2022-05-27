package resolver_test

var getHealthCheckQuery = `query {
		getHealthCheck { 
			status
		}
	}`

type GetHealthCheckQueryResponse struct {
	GetHealthCheck struct {
		Status string
	}
}
