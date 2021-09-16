package vcs

type RepoManager interface {
	Checkout(branch string) error
	ForceCheckout(branch string) error
	Fetch() error
	FetchAll() error
	RemoteBranches() ([]string, error)

	// ListChangedFiles returns slice of changed files
	ListChangedFiles() ([]string, error)
}
