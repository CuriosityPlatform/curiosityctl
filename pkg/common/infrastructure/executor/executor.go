//nolint
package executor

import (
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
)

type Executor interface {
	Output(args []string, opts ...Opt) ([]byte, error)
	Run(args []string, opts ...Opt) error
	PTY(ctx context.Context, args []string, opts ...Opt) error
}

func New(executable string) (Executor, error) {
	_, err := exec.LookPath(executable)
	if err != nil {
		return nil, err
	}

	return &executor{executable: executable}, nil
}

type executor struct {
	executable string
}

func (e *executor) Run(args []string, opts ...Opt) error {
	cmd := exec.Command(e.executable, args...)
	for _, opt := range opts {
		opt.apply(cmd)
	}
	return cmd.Run()
}

func (e *executor) Output(args []string, opts ...Opt) ([]byte, error) {
	cmd := exec.Command(e.executable, args...)
	for _, opt := range opts {
		opt.apply(cmd)
	}
	return cmd.Output()
}

func (e *executor) PTY(ctx context.Context, args []string, opts ...Opt) error {
	cmd := exec.Command(e.executable, args...)
	for _, opt := range opts {
		opt.apply(cmd)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)

	g, ctx := errgroup.WithContext(ctx)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	defer func() { _ = ptmx.Close() }()

	g.Go(func() error {
		for {
			select {
			case <-ch:
				if err2 := pty.InheritSize(os.Stdin, ptmx); err2 != nil {
					return errors.Wrap(err2, "error resizing pty")
				}
			case <-ctx.Done():
				return nil
			}
		}
	})
	ch <- syscall.SIGWINCH // Initial resize
	defer func() {
		// Stop notify subscription
		signal.Stop(ch)
		close(ch)
	}()

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()

	_, _ = io.Copy(os.Stdout, ptmx)

	return nil
}
