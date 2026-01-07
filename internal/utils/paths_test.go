package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRootDir_DefaultReturn_Success(t *testing.T) {
	// == arrange ==

	// == act ==
	got := DefaultRootDir()

	// == assert ==
	assert.Equal(t, "/etc/raind/container", got)
}

func TestDefaultRootDir_EnvSet_Success(t *testing.T) {
	// == arrange ==
	t.Setenv("RAIND_ROOT_DIR", "/path/to/root")

	// == act ==
	got := DefaultRootDir()

	// == assert ==
	assert.Equal(t, "/path/to/root", got)
}

func TestContainerDir_Success(t *testing.T) {
	// == arrange ==
	containerId := "12345"

	// == act ==
	got := ContainerDir(containerId)

	// == assert ==
	assert.Equal(t, "/etc/raind/container/12345", got)
}

func TestConfigFilePath_Success(t *testing.T) {
	// == arrange ==
	containerId := "12345"

	// == act ==
	got := ConfigFilePath(containerId)

	// == assert ==
	assert.Equal(t, "/etc/raind/container/12345/config.json", got)
}

func TestFifoPath_Success(t *testing.T) {
	// == arrange ==
	containerId := "12345"

	// == act ==
	got := FifoPath(containerId)

	// == assert ==
	assert.Equal(t, "/etc/raind/container/12345/exec.fifo", got)
}
