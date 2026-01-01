package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"droplet/internal/spec"
)

// NewContainerCreator constructs a ContainerCreator with the default
// implementations of its dependencies (SpecLoader, FifoCreator, ProcessExecutor).
// This acts as the main entry point for the container creation workflow.
func NewContainerCreator() *ContainerCreator {
	return &ContainerCreator{
		specLoader:      newFileSpecLoader(),
		fifoCreator:     newContainerFifoHandler(),
		processExecutor: newContainerInitExecutor(),
	}
}

// ContainerCreator orchestrates the container creation flow.
//
// The flow currently consists of:
//
//  1. Loading the OCI spec (config.json)
//  2. Creating the FIFO used for init synchronization
//  3. Launching the init process via the init subcommand
//
// Each step is delegated to an interface to allow testing and substitution.
type ContainerCreator struct {
	specLoader      specLoader
	fifoCreator     fifoCreator
	processExecutor processExecutor
}

// Create executes the container creation pipeline for the given container ID.
// This method performs no low-level work itself — it coordinates collaborators.
func (c *ContainerCreator) Create(opt CreateOption) error {
	// load config.json
	spec, err := c.specLoader.loadFile(opt.ContainerId)
	if err != nil {
		return err
	}

	// create fifo
	fifo := fifoPath(opt.ContainerId)
	if err := c.fifoCreator.createFifo(fifo); err != nil {
		return err
	}

	// execute init subcommand
	initPid, err := c.processExecutor.executeInit(spec, fifo)
	if err != nil {
		return err
	}

	fmt.Printf("init process has been created. pid: %d\n", initPid)

	return nil
}

// newContainerInitExecutor constructs a containerInitExecutor with the default implementations.
// This acts as the main entry point for spawning the container init process workflow.
func newContainerInitExecutor() *containerInitExecutor {
	return &containerInitExecutor{
		commandFactory: &execCommandFactory{},
	}
}

// processExecutor defines the behavior for spawning the container init process.
//
// It is an interface so that the behavior can be mocked in tests and
// replaced by alternative implementations if needed.
type processExecutor interface {
	executeInit(spec spec.Spec, fifo string) (int, error)
}

// containerInitExecutor is the default implementation of processExecutor.
//
// It invokes this binary with the `init` subcommand and the FIFO path,
// passing the spec's process args as the container entrypoint.
type containerInitExecutor struct {
	commandFactory commandFactory
}

// executeInit starts the init process and returns its PID.
//
// The init process is started as a child of the current runtime binary.
// The FIFO path is passed as an argument so that the init process can
// synchronize with the runtime.
func (c *containerInitExecutor) executeInit(spec spec.Spec, fifo string) (int, error) {
	// retrieve entrypoint from spec
	entrypoint := spec.Process.Args

	// prepare init subcommand
	initArgs := append([]string{"init", fifo}, entrypoint...)
	cmd := c.commandFactory.Command(os.Args[0], initArgs...)
	// set stdout/stderr
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	// execute init subcommand
	if err := cmd.Start(); err != nil {
		return -1, err
	}

	return cmd.Pid(), nil
}

// commandFactory creates commandExecutor instances.
//
// The factory abstracts process creation so that callers do not depend
// directly on exec.Command. This makes the behavior testable by replacing
// the factory with a mock implementation.
type commandFactory interface {
	Command(name string, args ...string) commandExecutor
}

// execCommandFactory is the default implementation of commandFactory.
//
// It creates commandExecutor values backed by *exec.Cmd and launches
// real OS processes.
type execCommandFactory struct{}

// Command returns a commandExecutor that executes the given command
// using exec.Cmd.
func (e *execCommandFactory) Command(name string, args ...string) commandExecutor {
	return &execCmd{cmd: exec.Command(name, args...)}
}

// commandExecutor represents a process that can be started.
//
// It provides a minimal surface over exec.Cmd so that command execution
// can be substituted or mocked in tests.
type commandExecutor interface {
	Start() error
	Pid() int
	SetStdout(w io.Writer)
	SetStderr(w io.Writer)
}

// execCmd is the concrete commandExecutor backed by exec.Cmd.
//
// It delegates all operations to the underlying exec.Cmd instance.
type execCmd struct {
	cmd *exec.Cmd
}

// Start starts the underlying process.
//
// It mirrors (*exec.Cmd).Start.
func (e *execCmd) Start() error {
	return e.cmd.Start()
}

// Pid returns the PID of the started process.
//
// If the process has not been started, -1 is returned.
func (e *execCmd) Pid() int {
	if e.cmd.Process == nil {
		return -1
	}
	return e.cmd.Process.Pid
}

// SetStdout sets the stdout writer for the underlying command.
func (e *execCmd) SetStdout(w io.Writer) {
	e.cmd.Stdout = w
}

// SetStderr sets the stderr writer for the underlying command.
func (e *execCmd) SetStderr(w io.Writer) {
	e.cmd.Stderr = w
}
