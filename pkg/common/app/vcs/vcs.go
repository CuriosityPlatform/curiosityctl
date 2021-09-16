package vcs

type VCS interface {
	WithRepoPath(path string) RepoManager
	WithClonedRepo(remoteURL string, to *string) (RepoManager, error)
}
