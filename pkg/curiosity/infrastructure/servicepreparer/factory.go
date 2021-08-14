package servicepreparer

import (
	"fmt"

	"curiosity/pkg/common/app/dockerclient"
	"curiosity/pkg/curiosity/app/servicepreparer"
)

const (
	mysqlService = "services-db"
)

func NewFactory(client dockerclient.Client) servicepreparer.Factory {
	return &factory{client: client}
}

type factory struct {
	client dockerclient.Client
}

//nolint:gocritic
func (f *factory) Preparer(service string) (servicepreparer.ServicePreparer, error) {
	switch service {
	case mysqlService:
		return &mysqlPreparer{
			client:     f.client,
			dbUser:     "root",
			dbPassword: "1234",
		}, nil
	}

	return nil, fmt.Errorf("unknown service %s for preparation", service)
}
