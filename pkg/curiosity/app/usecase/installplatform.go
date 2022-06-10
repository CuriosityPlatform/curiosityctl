package usecase

import (
	"context"
	"fmt"

	"curiosity/pkg/common/app/vcs"
	"curiosity/pkg/common/infrastructure/progress"
)

const (
	//nolint
	platformDefaultPath = "platform"
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

func (c *InstallPlatform) Execute(ctx context.Context, params InstallPlatformExecuteParams) (platformPath string, err error) {
	return platformPath, progress.Run(ctx, func(ctx context.Context) error {
		platformPath = platformDefaultPath
		if params.OutPath != nil {
			platformPath = *params.OutPath
		}

		w := progress.ContextWriter(ctx)
		eventID := fmt.Sprintf("Clone platform to %s", platformPath)
		w.Event(progress.StartedEvent(eventID))

		_, err = c.vcs.WithClonedRepo(params.RepoRemoteURL, params.OutPath)
		if err != nil {
			progress.ErrorEvent(eventID)
			return err
		}

		w.Event(progress.StoppedEvent(eventID))
		return nil
	})
}
