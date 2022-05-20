package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const (
	AppID = "curiosityctl"
)

func parseConfig() (*config, error) {
	c := &config{
		KubernetesManifestsBasePath: "kubernetes/app",
		DockerRegistry:              "registry.makerov.space:5000",
	}

	if err := envconfig.Process(AppID, c); err != nil {
		return nil, errors.Wrap(err, "failed to parse env")
	}
	return c, nil
}

type config struct {
	PlatformRoot string `envconfig:"platform_root" required:"1"`

	KubernetesManifestsBasePath string
	DockerRegistry              string `envconfig:"docker_registry"`
}
