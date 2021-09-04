package usecase

import (
	"context"
)

func NewRestart(
	up *Up,
	down *Down,
) *Restart {
	return &Restart{
		up:   up,
		down: down,
	}
}

type Restart struct {
	up   *Up
	down *Down
}

func (c *Restart) Execute(ctx context.Context) (err error) {
	err = c.down.Execute(ctx)
	if err != nil {
		return err
	}

	return c.up.Execute(ctx)
}
