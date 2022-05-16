package healthcheck

// IService interface is a collection of function signatures that represents the healthcheck's service contract.
type IService interface {
	GetStatus() error
}
