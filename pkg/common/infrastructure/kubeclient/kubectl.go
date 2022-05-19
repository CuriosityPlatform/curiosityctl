package kubeclient

import (
	"os"

	"github.com/pkg/errors"

	"curiosity/pkg/common/infrastructure/executor"
)

const (
	kubectlExecutable = "kubectl"
)

type kubectl interface {
	ApplyPath(path string) error
}

type k struct {
	kExecutor executor.Executor
}

func newK() (*k, error) {
	kubectl, err := executor.New(kubectlExecutable)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &k{kExecutor: kubectl}, nil
}

func (k *k) ApplyPath(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return errors.Wrapf(err, "failed to apply")
	}
	return k.kExecutor.Run([]string{"apply", "-f", path}, executor.WithStdout(os.Stdout))
}
