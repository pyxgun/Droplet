package container

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartExecute_Success(t *testing.T) {
	fifo := "/tmp/fifo"
	containerId := "123456"

	dummyFifoHandler := &dummyFifoHandler{path: fifo}
	dummyContainerStart := &ContainerStart{
		fifoHandler: dummyFifoHandler,
	}

	result := dummyContainerStart.Execute(StartOption{ContainerId: containerId})

	// assert
	// 1. writeFifo() is being called
	expectWriteFifoFlag := true
	resultWriteFifoFlag := dummyFifoHandler.callFlag
	assert.Equal(t, expectWriteFifoFlag, resultWriteFifoFlag)

	// 2. removeFifo() is being called
	expectRemoveFifoFlag := true
	resultRemoveFifoFlag := dummyFifoHandler.callFlag
	assert.Equal(t, expectRemoveFifoFlag, resultRemoveFifoFlag)

	// 3. nil is returned
	assert.Equal(t, nil, result)
}

func TestStartExecute_WriteFifoError(t *testing.T) {
	fifo := "/tmp/fifo"
	containerId := "123456"

	dummyFifoHandler := &dummyFifoHandler{path: fifo, writeErr: fmt.Errorf("failed to write FIFO")}
	dummyContainerStart := &ContainerStart{
		fifoHandler: dummyFifoHandler,
	}

	result := dummyContainerStart.Execute(StartOption{ContainerId: containerId})

	expect := fmt.Errorf("failed to write FIFO")

	assert.Equal(t, expect, result)
}

func TestStartExecute_RemoveFifoError(t *testing.T) {
	fifo := "/tmp/fifo"
	containerId := "123456"

	dummyFifoHandler := &dummyFifoHandler{path: fifo, removeErr: fmt.Errorf("failed to write FIFO")}
	dummyContainerStart := &ContainerStart{
		fifoHandler: dummyFifoHandler,
	}

	result := dummyContainerStart.Execute(StartOption{ContainerId: containerId})

	expect := fmt.Errorf("failed to write FIFO")

	assert.Equal(t, expect, result)
}
