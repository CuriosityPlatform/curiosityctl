package git

import "curiosity/pkg/common/infrastructure/executor"

const (
	gitPath = "git"
)

type Executor interface {
	executor.Executor
}

func NewExecutor() (Executor, error) {
	return executor.New(gitPath)
}
