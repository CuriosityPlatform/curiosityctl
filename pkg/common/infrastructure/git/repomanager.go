package git

import (
	"math"
	"regexp"
	"strings"

	"curiosity/pkg/common/app/vcs"
	"curiosity/pkg/common/infrastructure/executor"
)

func newRepoManager(repoPath string, gitExecutor Executor) vcs.RepoManager {
	return &repoManager{
		repoPath: repoPath,
		executor: gitExecutor,
	}
}

type repoManager struct {
	repoPath string
	executor Executor
}

func (repo *repoManager) Checkout(branch string) error {
	return repo.run("checkout", branch)
}

func (repo *repoManager) ForceCheckout(branch string) error {
	return repo.run("checkout", "-f", branch)
}

func (repo *repoManager) Fetch() error {
	return repo.run("fetch")
}

func (repo *repoManager) FetchAll() error {
	return repo.run("fetch", "--all")
}

func (repo *repoManager) RemoteBranches() ([]string, error) {
	output, err := repo.output("remote", "-v")
	if err != nil {
		return nil, err
	}

	reg := regexp.MustCompile(`(^.+?)\s`)

	//nolint:prealloc
	var branches []string

	for i, s := range strings.Split(string(output), "\n") {
		if math.Mod(float64(i), 2) != 0 {
			continue
		}
		if s == "" {
			continue
		}
		branches = append(branches, strings.TrimSpace(reg.FindString(s)))
	}
	return branches, nil
}

func (repo *repoManager) ListChangedFiles() ([]string, error) {
	output, err := repo.output("status", "-s")
	if err != nil {
		return nil, err
	}

	const modifiedPrefix = "M"

	var result []string

	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, modifiedPrefix) {
			result = append(result, line)
		}
	}

	return result, nil
}

func (repo *repoManager) run(args ...string) error {
	return repo.executor.Run(args, executor.WithWorkdir(repo.repoPath))
}

func (repo *repoManager) output(args ...string) ([]byte, error) {
	return repo.executor.Output(args, executor.WithWorkdir(repo.repoPath))
}
