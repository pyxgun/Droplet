package status

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerStatus_CREATING(t *testing.T) {
	// == arrange ==

	// == act ==

	// == assert ==
	assert.Equal(t, ContainerStatus(0), CREATING)
	assert.Equal(t, "creating", CREATING.String())
}

func TestContainerStatus_CREATED(t *testing.T) {
	// == arrange ==

	// == act ==

	// == assert ==
	assert.Equal(t, ContainerStatus(1), CREATED)
	assert.Equal(t, "created", CREATED.String())
}

func TestContainerStatus_RUNNING(t *testing.T) {
	// == arrange ==

	// == act ==

	// == assert ==
	assert.Equal(t, ContainerStatus(2), RUNNING)
	assert.Equal(t, "running", RUNNING.String())
}

func TestContainerStatus_STOPPED(t *testing.T) {
	// == arrange ==

	// == act ==

	// == assert ==
	assert.Equal(t, ContainerStatus(3), STOPPED)
	assert.Equal(t, "stopped", STOPPED.String())
}

func TestParseContainerStatus_Success(t *testing.T) {
	// == arrange ==
	creatingStatus := "creating"
	createdStatus := "created"
	runningStatus := "running"
	stoppedStatus := "stopped"

	// == act ==
	creatingStatusGot, creatingStatusErr := ParseContainerStatus(creatingStatus)
	createdStatusGot, createdStatusErr := ParseContainerStatus(createdStatus)
	runningStatusGot, runningStatusErr := ParseContainerStatus(runningStatus)
	stoppedStatusGot, stoppedStatusErr := ParseContainerStatus(stoppedStatus)

	// == assert ==
	// error is nil
	assert.Nil(t, creatingStatusErr)
	assert.Nil(t, createdStatusErr)
	assert.Nil(t, runningStatusErr)
	assert.Nil(t, stoppedStatusErr)

	// return int value
	assert.Equal(t, CREATING, creatingStatusGot)
	assert.Equal(t, CREATED, createdStatusGot)
	assert.Equal(t, RUNNING, runningStatusGot)
	assert.Equal(t, STOPPED, stoppedStatusGot)
}
