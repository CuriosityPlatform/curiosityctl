package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func parseConfig() (*config, error) {
	c := &config{
		KubernetesManifestsBasePath: "kubernetes/app",
	}

	if err := envconfig.Process(appID, c); err != nil {
		return nil, errors.Wrap(err, "failed to parse env")
	}
	return c, nil
}

type config struct {
	PlatformRoot string `envconfig:"platform_root" required:"1"`

	KubernetesManifestsBasePath string
}
