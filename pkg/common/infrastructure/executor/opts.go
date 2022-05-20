package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Opt interface {
	apply(cmd *exec.Cmd)
}

type executorOptFunc func(cmd *exec.Cmd)

func (f executorOptFunc) apply(cmd *exec.Cmd) {
	f(cmd)
}

func WithWorkdir(workdir string) Opt {
	return executorOptFunc(func(cmd *exec.Cmd) {
		cmd.Dir = workdir
	})
}

func WithStdin(reader io.Reader) Opt {
	return executorOptFunc(func(cmd *exec.Cmd) {
		cmd.Stdin = reader
	})
}

func WithStdout(writer io.Writer) Opt {
	return executorOptFunc(func(cmd *exec.Cmd) {
		cmd.Stdout = writer
	})
}

func WithEnv(envMap map[string]string) Opt {
	return executorOptFunc(func(cmd *exec.Cmd) {
		env := make([]string, 0, len(envMap))
		for k, v := range envMap {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = append(os.Environ(), env...)
	})
}
