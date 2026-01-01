package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerDir_Success(t *testing.T) {
	containerId := "123456"

	result := containerDir(containerId)

	expect := "/etc/raind/container/123456"

	assert.Equal(t, expect, result)
}

func TestConfigFilePath_Success(t *testing.T) {
	containerId := "123456"

	result := configFilePath(containerId)

	expect := "/etc/raind/container/123456/config.json"

	assert.Equal(t, expect, result)
}

func TestFifoPath_Success(t *testing.T) {
	containerId := "123456"

	result := fifoPath(containerId)

	expect := "/etc/raind/container/123456/exec.fifo"

	assert.Equal(t, expect, result)
}
