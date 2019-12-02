package cmd

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	lh_client "github.com/longhorn/longhorn-manager/client"
	"github.com/stretchr/testify/assert"
	"github.com/utilitywarehouse/lhctl/util"
	"testing"
)

func TestParseDetachArgs(t *testing.T) {

	// Test: parse with no args returns error
	expectedErr := errors.New("volume name not specified")
	_, err := parseDetachArgs(detachCmd, []string{})

	assert.Equal(t, expectedErr, err)

	// Test: first argument returned as volume name
	vol, err := parseDetachArgs(detachCmd, []string{"volume"})

	assert.Equal(t, nil, err)
	assert.Equal(t, vol, "volume")

	// Test: Rest of arguments ignored
	vol, err = parseDetachArgs(detachCmd, []string{"volume", "other"})

	assert.Equal(t, nil, err)
	assert.Equal(t, vol, "volume")
}

func TestDetachVolume(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	// Test: Error when get volume errors
	volErr := errors.New("volume error")
	mockClient.EXPECT().GetVolume("volume").Times(1).Return(nil, volErr)
	err := detachVolume("volume")

	assert.Equal(t, volErr, err)

	// Test: Error when detach volume errors
	vol := &lh_client.Volume{}
	mockClient.EXPECT().GetVolume("volume").Times(1).Return(vol, nil)

	detachErr := errors.New("detach error")
	mockClient.EXPECT().VolumeDetach(vol).Times(1).Return(nil, detachErr)

	err = detachVolume("volume")
	assert.Equal(t, detachErr, err)

	// Test: nil when no call errors
	mockClient.EXPECT().GetVolume("volume").Times(1).Return(vol, nil)
	mockClient.EXPECT().VolumeDetach(vol).Times(1).Return(vol, nil)

	err = detachVolume("volume")
	assert.Equal(t, nil, err)
}

func TestTimeoutWhileWaitingForDetachedVolume(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	// Test: timeout returns error
	vol := &lh_client.Volume{}
	vol.State = "attached"
	mockClient.EXPECT().GetVolume("volume").AnyTimes().Return(vol, nil)

	timeoutErr := errors.New(fmt.Sprintf(
		"could not detach %s before timeout",
		"volume",
	))

	err := waitForDetachedVol("volume", 1)
	assert.Equal(t, timeoutErr, err)

	// Test: get volume error will not exit until timeout
	volErr := errors.New("volume error")
	mockClient.EXPECT().GetVolume("volume").AnyTimes().
		Return(vol, volErr)

	err = waitForDetachedVol("volume", 1)
	assert.Equal(t, timeoutErr, err)
}

func TestWaitForDetachedVol(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	// Test: node disabled while trying
	detached := &lh_client.Volume{}
	detached.State = "detached"
	attached := &lh_client.Volume{}
	attached.State = "attached"

	gomock.InOrder(
		mockClient.EXPECT().GetVolume("volume").Times(3).
			Return(attached, nil),
		mockClient.EXPECT().GetVolume("volume").Times(1).
			Return(detached, nil),
	)

	err := waitForDetachedVol("volume", 5)
	assert.Equal(t, nil, err)

	// Test: already disabled
	mockClient.EXPECT().GetVolume("volume").Times(1).Return(detached, nil)

	err = waitForDetachedVol("volume", 1)
	assert.Equal(t, nil, err)
}
