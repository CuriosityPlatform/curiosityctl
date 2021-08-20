package servicepreparer

import "context"

type ServicePreparer interface {
	Prepare(ctx context.Context, composeServiceName string) error
}
