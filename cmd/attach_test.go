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

func TestParseAttachArgs(t *testing.T) {

	// Test: parse with no args returns error
	expectedErr := errors.New("volume name not specified")
	err := parseAttachArgs(attachCmd, []string{})

	assert.Equal(t, expectedErr, err)

	// Test: flag node must be set
	expectedErr = errors.New("--node= flag must be set")
	AttachTargetNode = ""
	err = parseAttachArgs(attachCmd, []string{"volume"})

	assert.Equal(t, expectedErr, err)

	// Test: Store first argument as node
	AttachTargetNode = "node"
	err = parseAttachArgs(attachCmd, []string{"volume"})

	assert.Equal(t, nil, err)
	assert.Equal(t, AttachVolume, "volume")

	// Test: Rest of arguments ignored
	AttachTargetNode = "node"
	err = parseAttachArgs(attachCmd, []string{"volume", "other"})

	assert.Equal(t, nil, err)
	assert.Equal(t, AttachVolume, "volume")
}

func TestAttachVolume(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	// Input after parse
	AttachTargetNode = "node"
	AttachVolume = "volume"

	// Test: Error when get volume errors
	volErr := errors.New("volume error")
	mockClient.EXPECT().GetVolume("volume").Times(1).Return(nil, volErr)
	err := attachVolume()

	assert.Equal(t, volErr, err)

	// Test: Error when attach volume errors
	vol := &lh_client.Volume{}
	mockClient.EXPECT().GetVolume("volume").Times(1).Return(vol, nil)

	attachErr := errors.New("attach error")
	mockClient.EXPECT().VolumeAttach(vol, AttachTargetNode).
		Times(1).Return(nil, attachErr)

	err = attachVolume()
	assert.Equal(t, attachErr, err)

	// Test: nil when no call errors
	mockClient.EXPECT().GetVolume("volume").Times(1).Return(vol, nil)
	mockClient.EXPECT().VolumeAttach(vol, AttachTargetNode).
		Times(1).Return(vol, nil)

	err = attachVolume()
	assert.Equal(t, nil, err)
}

func TestTimeoutWhileWaitingForAttachedVolume(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	// Input after parse
	AttachTargetNode = "node"
	AttachVolume = "volume"

	// Test: timeout returns error
	vol := &lh_client.Volume{}
	vol.State = "detached"
	mockClient.EXPECT().GetVolume("volume").AnyTimes().Return(vol, nil)

	timeoutErr := errors.New(fmt.Sprintf(
		"could not attach %s before timeout",
		"volume",
	))

	err := waitForAttachedVol(1)
	assert.Equal(t, timeoutErr, err)

	// Test: get volume error will not exit until timeout
	volErr := errors.New("volume error")
	mockClient.EXPECT().GetVolume("volume").AnyTimes().
		Return(vol, volErr)

	err = waitForAttachedVol(1)
	assert.Equal(t, timeoutErr, err)
}

func TestWaitForAttachedVol(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	// Input after parse
	AttachTargetNode = "node"
	AttachVolume = "volume"

	// Test: node disabled while trying
	detached := &lh_client.Volume{}
	detached.State = "detached"
	attached := &lh_client.Volume{}
	attached.State = "attached"

	gomock.InOrder(
		mockClient.EXPECT().GetVolume("volume").Times(3).
			Return(detached, nil),
		mockClient.EXPECT().GetVolume("volume").Times(1).
			Return(attached, nil),
	)

	err := waitForAttachedVol(5)
	assert.Equal(t, nil, err)

	// Test: already attached
	mockClient.EXPECT().GetVolume("volume").Times(1).Return(attached, nil)

	err = waitForAttachedVol(1)
	assert.Equal(t, nil, err)
}
