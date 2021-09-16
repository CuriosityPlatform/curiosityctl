package git

import (
	"os/exec"

	"github.com/pkg/errors"

	"curiosity/pkg/common/app/vcs"
)

func NewGitVcs() (vcs.VCS, error) {
	executor, err := NewExecutor()
	if err != nil {
		return nil, err
	}

	return &gitVcs{
		executor: executor,
	}, nil
}

type gitVcs struct {
	executor Executor
}

func (git *gitVcs) WithRepoPath(path string) vcs.RepoManager {
	return newRepoManager(path, git.executor)
}

func (git *gitVcs) WithClonedRepo(remoteURL string, to *string) (vcs.RepoManager, error) {
	projectName, err := projectNameFromUrl(remoteURL)
	if err != nil {
		return nil, err
	}

	args := []string{"clone", remoteURL}

	if to != nil {
		args = append(args, *to)
	}

	_, err = git.executor.Output(args)
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, errors.Wrap(err, string(ee.Stderr))
		}
		return nil, err
	}

	if to != nil {
		return git.WithRepoPath(*to), nil
	}

	return git.WithRepoPath(projectName), nil
}
