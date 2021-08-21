package compose

import "github.com/compose-spec/compose-go/types"

type Labels types.Labels

func NewProject(composeProject *types.Project) Project {
	return &project{composeProject: composeProject}
}

type Project interface {
	AwaitableServices() ([]Service, error)
	BootableServices() ([]Service, error)
	WithAwaitableServices(func(service Service) error) error
	WithBootableServices(func(service Service) error) error
}

type project struct {
	composeProject *types.Project
}

func (project *project) AwaitableServices() ([]Service, error) {
	var services []Service
	for _, config := range project.composeProject.AllServices() {
		_, exists := config.Labels[waitLabel]
		if exists {
			services = append(services, &service{composeService: config})
		}
	}
	return services, nil
}

func (project *project) BootableServices() ([]Service, error) {
	var services []Service
	for _, config := range project.composeProject.AllServices() {
		_, exists := config.Labels[waitLabel]
		if exists {
			services = append(services, &service{composeService: config})
		}
	}
	return services, nil
}

func (project *project) WithAwaitableServices(f func(service Service) error) error {
	for _, config := range project.composeProject.AllServices() {
		_, exists := config.Labels[bootLabel]
		if exists {
			err := f(&service{composeService: config})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (project *project) WithBootableServices(f func(service Service) error) error {
	for _, config := range project.composeProject.AllServices() {
		_, exists := config.Labels[bootLabel]
		if exists {
			err := f(&service{composeService: config})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type Service interface {
	Name() string
	ContainerName() string
	Labels() Labels
}

type service struct {
	composeService types.ServiceConfig
}

func (service *service) Name() string {
	return service.composeService.Name
}

func (service *service) ContainerName() string {
	return service.composeService.ContainerName
}

func (service *service) Labels() Labels {
	return Labels(service.composeService.Labels)
}
