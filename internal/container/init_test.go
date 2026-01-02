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
	dummySpecLoader := &dummyFileSpecLoader{spec: dummySpec()}
	dummyContainerRootEnvPreparer := &dummyContainerRootEnvPreparer{spec: dummySpec()}
	dummySyscallHandler := &dummySyscallHandler{argv0: entrypoint[0], argv: entrypoint, envv: os.Environ()}
	dummyContainerInit := &ContainerInit{
		fifoReader:           dummyFifoHandler,
		specLoader:           dummySpecLoader,
		containerEnvPreparer: dummyContainerRootEnvPreparer,
		syscallHandler:       dummySyscallHandler,
	}

	result := dummyContainerInit.Execute(InitOption{Fifo: fifo, Entrypoint: entrypoint})

	// assert
	// 1. readFifo() is being called
	expectReadFifoFlag := true
	resultReadFifoFlag := dummyFifoHandler.callFlag
	assert.Equal(t, expectReadFifoFlag, resultReadFifoFlag)

	// 2. loadFile() is being called
	expectLoadFileFlag := true
	resultLoadFileFlag := dummySpecLoader.callFlag
	assert.Equal(t, expectLoadFileFlag, resultLoadFileFlag)

	// 3. prepare() is being called
	expectPrepareFlag := true
	resultPrepareFlag := dummyContainerRootEnvPreparer.prepareCallFlag
	assert.Equal(t, expectPrepareFlag, resultPrepareFlag)

	// 4. Exec() is being called
	expectExecFlag := true
	resultExecFlag := dummySyscallHandler.execCallFlag
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

func TestInitExecute_LoadFileError(t *testing.T) {
	fifo := "/tmp/fifo"
	entrypoint := []string{"/bin/sh", "-c", "Hello World"}

	dummyFifoHandler := &dummyFifoHandler{path: fifo}
	dummySpecLoader := &dummyFileSpecLoader{spec: dummySpec(), err: fmt.Errorf("failed to load spec file")}
	dummyContainerInit := &ContainerInit{
		fifoReader: dummyFifoHandler,
		specLoader: dummySpecLoader,
	}

	result := dummyContainerInit.Execute(InitOption{Fifo: fifo, Entrypoint: entrypoint})

	expect := fmt.Errorf("failed to load spec file")

	assert.Equal(t, expect, result)
}

func TestInitExecute_PrepareError(t *testing.T) {
	fifo := "/tmp/fifo"
	entrypoint := []string{"/bin/sh", "-c", "Hello World"}

	dummyFifoHandler := &dummyFifoHandler{path: fifo}
	dummySpecLoader := &dummyFileSpecLoader{spec: dummySpec()}
	dummyContainerRootEnvPreparer := &dummyContainerRootEnvPreparer{spec: dummySpec(), err: fmt.Errorf("failed to prepare environment")}
	dummyContainerInit := &ContainerInit{
		fifoReader:           dummyFifoHandler,
		specLoader:           dummySpecLoader,
		containerEnvPreparer: dummyContainerRootEnvPreparer,
	}

	result := dummyContainerInit.Execute(InitOption{Fifo: fifo, Entrypoint: entrypoint})

	expect := fmt.Errorf("failed to prepare environment")

	assert.Equal(t, expect, result)
}

func TestInitExecute_syscallExecError(t *testing.T) {
	fifo := "/tmp/fifo"
	entrypoint := []string{"/bin/sh", "-c", "Hello World"}

	dummyFifoHandler := &dummyFifoHandler{path: fifo}
	dummySpecLoader := &dummyFileSpecLoader{spec: dummySpec()}
	dummyContainerRootEnvPreparer := &dummyContainerRootEnvPreparer{spec: dummySpec()}
	dummySyscallHandler := &dummySyscallHandler{
		argv0: entrypoint[0], argv: entrypoint, envv: os.Environ(),
		err: fmt.Errorf("syscall.Exec failed"),
	}
	dummyContainerInit := &ContainerInit{
		fifoReader:           dummyFifoHandler,
		specLoader:           dummySpecLoader,
		containerEnvPreparer: dummyContainerRootEnvPreparer,
		syscallHandler:       dummySyscallHandler,
	}

	result := dummyContainerInit.Execute(InitOption{Fifo: fifo, Entrypoint: entrypoint})

	expect := fmt.Errorf("syscall.Exec failed")

	assert.Equal(t, expect, result)
}

func TestRootContainerENvPreparer_Success(t *testing.T) {
	dummySyscallHandler := &dummySyscallHandler{}
	dummyRootContainerEnvPrepare := &rootContainerEnvPreparer{
		syscallHandler: dummySyscallHandler,
	}

	result := dummyRootContainerEnvPrepare.prepare(dummySpec())

	// Setresgid() is being called
	expectSetresgidFlag := true
	resultSetresgidFlag := dummySyscallHandler.setresgidCallFlag
	assert.Equal(t, expectSetresgidFlag, resultSetresgidFlag)

	// Setresuid() is being called
	expectSetresuidFlag := true
	resultSetresuidFlag := dummySyscallHandler.setresuidCallFlag
	assert.Equal(t, expectSetresuidFlag, resultSetresuidFlag)

	// Sethostname() is being called
	expectSethostnameFlag := true
	resultSethostnameFlag := dummySyscallHandler.sethostnameCallFlag
	assert.Equal(t, expectSethostnameFlag, resultSethostnameFlag)

	// nil is returned
	assert.Equal(t, nil, result)
}
