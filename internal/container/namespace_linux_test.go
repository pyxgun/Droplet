package container

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildNamespaceConfig_Success(t *testing.T) {
	// dummySpec
	// namespace target: mount, uts
	dummySpec := dummySpec()

	nsCfg := buildNamespaceConfig(dummySpec)

	// assert
	// mount = true
	assert.Equal(t, true, nsCfg.mount)
	// network = false
	assert.Equal(t, false, nsCfg.network)
	// uts = true
	assert.Equal(t, true, nsCfg.uts)
	// pid = false
	assert.Equal(t, false, nsCfg.pid)
	// ipc = false
	assert.Equal(t, false, nsCfg.ipc)
	// user = false
	assert.Equal(t, false, nsCfg.user)
	// cgroup = false
	assert.Equal(t, false, nsCfg.cgroup)
}

func TestBuildNamespaceAttr_Success(t *testing.T) {
	// dummySpec
	// namespace target: mount, uts
	dummySpec := dummySpec()

	nsCfg := buildNamespaceConfig(dummySpec)

	nsAttr := buildNamespaceAttr(nsCfg)

	var attr uintptr
	attr |= (syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS)

	assert.Equal(t, attr, nsAttr.Cloneflags)
}
