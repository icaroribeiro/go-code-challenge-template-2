package resolver_test

type GetHealthCheckQueryResponse struct {
	GetHealthCheck struct {
		Status string
	}
}

var getHealthCheckQuery = `query {
		getHealthCheck { 
			status
		}
	}`
