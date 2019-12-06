package cmd

import (
	"errors"
	"github.com/golang/mock/gomock"
	lh_client "github.com/longhorn/longhorn-manager/client"
	"github.com/stretchr/testify/assert"
	"github.com/utilitywarehouse/lhctl/util"
	"testing"
)

func TestParseUpdateVolumeArgs(t *testing.T) {

	// Test: parse with no args returns error
	expectedErr := errors.New("volume name not specified")
	err := parseUpdateVolumeArgs(updateVolumeCmd, []string{})

	assert.Equal(t, expectedErr, err)

	// Test: No flag returns an error
	expectedErr = errors.New("at least one flag should be set")
	err = parseUpdateVolumeArgs(updateVolumeCmd, []string{"volume"})

	assert.Equal(t, expectedErr, err)

	// Test: Invalid replicas value
	expectedErr = errors.New("--replicas= flag invalid value")
	ReplicaCountFlag = "0"
	err = parseUpdateVolumeArgs(updateVolumeCmd, []string{"volume"})

	assert.Equal(t, expectedErr, err)

	// Test: Replicas flag and update volume vars are set
	ReplicaCountFlag = "1"
	err = parseUpdateVolumeArgs(updateVolumeCmd, []string{"volume"})

	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1), ReplicaCount)
	assert.Equal(t, "volume", UpdateVolumeName)
}

func TestUpdateVolumeError(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	volErr := errors.New("volume error")
	mockClient.EXPECT().GetVolume("volume").Times(1).Return(nil, volErr)
	err := updateVolume()

	assert.Equal(t, volErr, err)

}

func TestUpdateVolumeReplicas(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	// Test: Update count errors
	ReplicaCountFlag = "1"
	err := parseUpdateVolumeArgs(updateVolumeCmd, []string{"volume"})
	assert.Equal(t, nil, err)

	vol := &lh_client.Volume{}
	mockClient.EXPECT().GetVolume("volume").Times(1).Return(vol, err)
	updReplicaErr := errors.New("upd replicas error")
	mockClient.EXPECT().UpdateReplicaCount(vol, int64(1)).
		Times(1).Return(vol, updReplicaErr)
	err = updateVolume()

	assert.Equal(t, updReplicaErr, err)

	// Test: Update replicas
	ReplicaCountFlag = "1"
	err = parseUpdateVolumeArgs(updateVolumeCmd, []string{"volume"})
	assert.Equal(t, nil, err)

	vol = &lh_client.Volume{}
	mockClient.EXPECT().GetVolume("volume").Times(1).Return(vol, err)
	mockClient.EXPECT().UpdateReplicaCount(vol, int64(1)).
		Times(1).Return(vol, nil)
	err = updateVolume()

	assert.Equal(t, nil, err)

}
