package servicepreparer

type ServicePreparer interface {
	Prepare(composeServiceName string) error
}
