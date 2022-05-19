package service

type Deployer interface {
	DeployAll() error
	DeployApp(name string) error
}
