package container

import (
	"droplet/internal/spec"
)

type mockeContainerCgroupController struct {
	prepareCallFlag    bool
	preapreContainerId string
	prepareeSpec       spec.Spec
	preparePid         int
	prepareErr         error

	// createCgroupDirectory()
	createCgroupDirectoryErr error

	// setMemoryLimit()
	setMemoryLimitErr error

	// setCpuLimit()
	setCpuLimitErr error

	// setProcessToCgroup
	setProcessToCgroupErr error
}

func (m *mockeContainerCgroupController) prepare(containerId string, spec spec.Spec, pid int) error {
	m.prepareCallFlag = true
	m.preapreContainerId = containerId
	m.prepareeSpec = spec
	m.preparePid = pid
	return m.prepareErr
}
