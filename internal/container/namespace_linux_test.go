package container

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildNamespaceConfig_Success(t *testing.T) {
	// dummySpec
	// namespace target: mount, uts, user
	dummySpec := dummySpec()

	nsCfg := buildNamespaceConfig(dummySpec)

	// assert
	// mount = true
	assert.Equal(t, true, nsCfg.mount)
	// network = true
	assert.Equal(t, true, nsCfg.network)
	// uts = true
	assert.Equal(t, true, nsCfg.uts)
	// pid = true
	assert.Equal(t, true, nsCfg.pid)
	// ipc = true
	assert.Equal(t, true, nsCfg.ipc)
	// user = true
	assert.Equal(t, true, nsCfg.user)
	// cgroup = true
	assert.Equal(t, true, nsCfg.cgroup)
}

func TestBuildCloneFlags_Success(t *testing.T) {
	// dummySpec
	dummySpec := dummySpec()
	nsCfg := buildNamespaceConfig(dummySpec)
	cloneFlags := buildCloneFlags(nsCfg)

	var expect uintptr
	expect |= (syscall.CLONE_NEWNS |
		syscall.CLONE_NEWUTS |
		syscall.CLONE_NEWUSER |
		syscall.CLONE_NEWNET |
		syscall.CLONE_NEWPID |
		syscall.CLONE_NEWIPC |
		syscall.CLONE_NEWCGROUP)

	// assert
	assert.Equal(t, expect, cloneFlags)
}

func TestBuildRootUserNamespaceIDMap_Success(t *testing.T) {
	// dummySpec
	// namespace target: mount, uts
	dummySpec := dummySpec()
	nsCfg := buildNamespaceConfig(dummySpec)

	uidMap, gidMap := buildRootUserNamespaceIDMap(nsCfg)

	expectUidMap := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      0,
			Size:        65535,
		},
	}
	expectGidMap := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      0,
			Size:        65535,
		},
	}

	// assert
	assert.Equal(t, expectGidMap, gidMap)
	assert.Equal(t, expectUidMap, uidMap)
}

func TestBuildSysProcAttr_Success(t *testing.T) {
	// dummySpec
	// namespace target: mount, uts, user
	dummySpec := dummySpec()
	nsCfg := buildNamespaceConfig(dummySpec)
	cloneFlags := buildCloneFlags(nsCfg)
	uidMap, gidMap := buildRootUserNamespaceIDMap(nsCfg)

	procAttrStruct := procAttr{
		cloneFlags:    cloneFlags,
		uidMap:        uidMap,
		gidMap:        gidMap,
		setGroupsFlag: true,
	}

	procAttr := buildSysProcAttr(procAttrStruct)

	expect := syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWCGROUP,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        65535,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        65535,
			},
		},
		GidMappingsEnableSetgroups: true,
	}

	assert.Equal(t, &expect, procAttr)
}
