package container

import "droplet/internal/status"

type mockStatusHandler struct {
	// CreateStatusFile()
	createStatusFileCallFlag    bool
	createStatusFilePath        string
	createStatusFileContainerId string
	createStatusFilePid         int
	createStatusFileStatus      status.ContainerStatus
	createStatusFileBundle      string
	createStatusFileErr         error

	// UpdateStatus()
	updateStatusCallFlag    bool
	updateStatusPath        string
	updateStatusContainerId string
	updateStatusStatus      status.ContainerStatus
	updateStatusPid         int
	updateStatusErr         error
}

func (m *mockStatusHandler) CreateStatusFile(path string, containerId string, pid int, status status.ContainerStatus, bundle string) error {
	m.createStatusFileCallFlag = true
	m.createStatusFilePath = path
	m.createStatusFileContainerId = containerId
	m.createStatusFilePid = pid
	m.createStatusFileStatus = status
	m.createStatusFileBundle = bundle
	return m.createStatusFileErr
}

func (m *mockStatusHandler) UpdateStatus(path string, containerId string, status status.ContainerStatus, pid int) error {
	m.updateStatusCallFlag = true
	m.updateStatusPath = path
	m.updateStatusContainerId = containerId
	m.updateStatusStatus = status
	m.updateStatusPid = pid
	return m.updateStatusErr
}
