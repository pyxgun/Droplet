package container

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitExecute_Success(t *testing.T) {
	fifo := "/tmp/fifo"
	entrypoint := []string{"/bin/sh", "-c", "Hello World"}

	dummyFifoHandler := &dummyFifoHandler{path: fifo}
	dummyProcessReplacer := &dummySyscallProcessReplacer{argv0: entrypoint[0], argv: entrypoint, envv: os.Environ()}
	dummyContainerInit := &ContainerInit{
		fifoReader:      dummyFifoHandler,
		processReplacer: dummyProcessReplacer,
	}

	result := dummyContainerInit.Execute(InitOption{Fifo: fifo, Entrypoint: entrypoint})

	// assert
	// 1. readFifo() is being called
	expectReadFifoFlag := true
	resultReadFifoFlag := dummyFifoHandler.callFlag
	assert.Equal(t, expectReadFifoFlag, resultReadFifoFlag)

	// 2. Exec() is being called
	expectExecFlag := true
	resultExecFlag := dummyProcessReplacer.callFlag
	assert.Equal(t, expectExecFlag, resultExecFlag)

	// 3. nil is returned
	assert.Equal(t, nil, result)
}

func TestInitExecute_ReadFifoError(t *testing.T) {
	fifo := "/tmp/fifo"
	entrypoint := []string{"/bin/sh", "-c", "Hello World"}

	dummyFifoHandler := &dummyFifoHandler{path: fifo, readErr: fmt.Errorf("read FIFO failed")}
	dummyContainerInit := &ContainerInit{
		fifoReader: dummyFifoHandler,
	}

	result := dummyContainerInit.Execute(InitOption{Fifo: fifo, Entrypoint: entrypoint})

	expect := fmt.Errorf("read FIFO failed")

	assert.Equal(t, expect, result)
}

func TestInitExecute_syscallExecError(t *testing.T) {
	fifo := "/tmp/fifo"
	entrypoint := []string{"/bin/sh", "-c", "Hello World"}

	dummyFifoHandler := &dummyFifoHandler{path: fifo}
	dummyProcessReplacer := &dummySyscallProcessReplacer{
		argv0: entrypoint[0], argv: entrypoint, envv: os.Environ(),
		err: fmt.Errorf("syscall.Exec failed"),
	}
	dummyContainerInit := &ContainerInit{
		fifoReader:      dummyFifoHandler,
		processReplacer: dummyProcessReplacer,
	}

	result := dummyContainerInit.Execute(InitOption{Fifo: fifo, Entrypoint: entrypoint})

	expect := fmt.Errorf("syscall.Exec failed")

	assert.Equal(t, expect, result)
}
