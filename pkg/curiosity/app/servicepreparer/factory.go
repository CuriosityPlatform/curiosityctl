package servicepreparer

type Factory interface {
	Preparer(service string) (ServicePreparer, error)
}
