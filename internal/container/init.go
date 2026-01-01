package container

import (
	"os"
	"syscall"
)

// NewContainerInit returns a ContainerInit wired with the default
// implementations of its dependencies (fifoReader and processReplacer).
// This is the standard entry point for executing the container init phase.
func NewContainerInit() *ContainerInit {
	return &ContainerInit{
		fifoReader:      newContainerFifoHandler(),
		processReplacer: newSyscallProcessReplacer(),
	}
}

// ContainerInit represents the runtime logic executed inside the
// container's init process.
//
// The init process waits for a start signal via FIFO and then
// replaces itself with the container entrypoint using execve-style
// semantics (syscall.Exec).
type ContainerInit struct {
	fifoReader      fifoReader
	processReplacer processReplacer
}

// Execute performs the init sequence for the container.
//
// The sequence is:
//
//  1. Wait for a start signal by reading from the FIFO path
//  2. Replace the current process image with the container entrypoint
//
// On success, this function does not return because the process image
// is replaced. Errors are returned only if the FIFO read fails or
// syscall.Exec cannot be invoked.
func (c *ContainerInit) Execute(opt InitOption) error {
	fifo := opt.Fifo
	entrypoint := opt.Entrypoint

	// read fifo for waiting start signal
	if err := c.fifoReader.readFifo(fifo); err != nil {
		return err
	}

	if err := c.processReplacer.Exec(entrypoint[0], entrypoint, os.Environ()); err != nil {
		return err
	}

	return nil
}

// processReplacer abstracts the operation of replacing the current
// process image with another program.
//
// It is defined as an interface to allow syscall.Exec to be mocked
// in tests and substituted by alternative implementations if needed.
type processReplacer interface {
	Exec(argv0 string, argv []string, envv []string) error
}

// newSyscallProcessReplacer returns a processReplacer that delegates to
// syscall.Exec to replace the current process image.
func newSyscallProcessReplacer() *syscallProcessReplacer {
	return &syscallProcessReplacer{}
}

// syscallProcessReplacer is the default implementation of processReplacer.
//
// It invokes syscall.Exec directly, causing the current process to be
// replaced by the specified executable if successful.
type syscallProcessReplacer struct{}

// Exec calls syscall.Exec with the provided arguments.
//
// On success, this call does not return. Any returned error indicates
// that the process could not be replaced.
func (s *syscallProcessReplacer) Exec(argv0 string, argv []string, envv []string) error {
	return syscall.Exec(argv0, argv, envv)
}
