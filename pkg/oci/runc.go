package oci

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/iximiuz/conman/pkg/container"
)

// implementes oci.Runtime interface
type runcRuntime struct {
	// path to runc executable, eg. /usr/local/bin/runc
	exePath string

	// dir to store container state (on tmpfs), eg. /run/runc/
	rootPath string
}

func NewRuntime(exePath, rootPath string) Runtime {
	return &runcRuntime{
		exePath:  exePath,
		rootPath: rootPath,
	}
}

func (r *runcRuntime) CreateContainer(id container.ID, bundle string) error {
	cmd := exec.Command(
		r.exePath,
		"--root", r.rootPath,
		"create",
		"--bundle", bundle,
		string(id),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func (r *runcRuntime) StartContainer(id container.ID) error {
	attrs := os.ProcAttr{
		Dir: r.rootPath, // TODO: container bundle dir
		Env: os.Environ(),
		Sys: &syscall.SysProcAttr{Noctty: true},
	}
	args := []string{r.exePath, "start", string(id)}
	proc, err := os.StartProcess(r.exePath, args, &attrs)
	if err != nil {
		return err
	}
	return proc.Release()
}

func (r *runcRuntime) KillContainer() {
	panic("not implemented")
}

func (r *runcRuntime) DeleteContainer() {
	panic("not implemented")
}

func (r *runcRuntime) ContainerState() {
	panic("not implemented")
}
