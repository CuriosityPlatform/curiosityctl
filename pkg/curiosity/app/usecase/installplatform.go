package usecase

import (
	"context"

	"curiosity/pkg/common/app/vcs"
	"curiosity/pkg/common/infrastructure/progress"
)

const (
	//nolint
	platformPathEnvName = "CURIOSITYCTL_PLATFORM_ROOT"
)

func NewInstallPlatform(v vcs.VCS) *InstallPlatform {
	return &InstallPlatform{vcs: v}
}

type InstallPlatform struct {
	vcs vcs.VCS
}

type InstallPlatformExecuteParams struct {
	OutPath       *string
	RepoRemoteURL string
}

func (c *InstallPlatform) Execute(ctx context.Context, params InstallPlatformExecuteParams) error {
	return progress.Run(ctx, func(ctx context.Context) error {
		w := progress.ContextWriter(ctx)
		eventID := "Clone platform"
		w.Event(progress.StartedEvent(eventID))

		_, err := c.vcs.WithClonedRepo(params.RepoRemoteURL, params.OutPath)
		if err != nil {
			progress.ErrorEvent(eventID)
			return err
		}

		w.Event(progress.StoppedEvent(eventID))
		return nil
	})
}
