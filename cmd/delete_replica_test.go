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

func TestValidateDeleteReplicaArgs(t *testing.T) {

	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	// Test: validate with no args raise error
	expectedErr := errors.New("No replica name specified")
	_, _, err := validateDeleteReplicaArgs(deletereplicaCmd, []string{})

	assert.Equal(t, expectedErr, err)

	// Test: validate errors when list volumes errors
	listingErr := errors.New("No replica name specified")
	mockClient.EXPECT().ListVolumes().Times(1).Return(nil, listingErr)
	_, _, err = validateDeleteReplicaArgs(deletereplicaCmd, []string{"replica"})

	assert.Equal(t, listingErr, err)

	// Test: Replica not found returns error
	volumes := []lh_client.Volume{
		lh_client.Volume{
			Name: "volume",
		},
	}
	volumes[0].Replicas = []lh_client.Replica{}
	mockClient.EXPECT().ListVolumes().Times(1).Return(volumes, nil)
	_, _, err = validateDeleteReplicaArgs(
		deletereplicaCmd,
		[]string{"replica"},
	)

	expectedErr = errors.New(fmt.Sprintf(
		"Replica not found: %s",
		"replica",
	))
	assert.Equal(t, expectedErr, err)

	// Test: Replica found sets vars and return nil
	volumes[0].Replicas = append(
		volumes[0].Replicas,
		lh_client.Replica{Name: "replica"},
	)
	mockClient.EXPECT().ListVolumes().Times(1).Return(volumes, nil)
	r, v, err := validateDeleteReplicaArgs(deletereplicaCmd, []string{"replica"})

	assert.Equal(t, r, "replica")
	assert.Equal(t, v, "volume")
	assert.Equal(t, nil, err)

	// Test: Additional arguments ignored
	mockClient.EXPECT().ListVolumes().Times(1).Return(volumes, nil)
	r, v, err = validateDeleteReplicaArgs(
		deletereplicaCmd,
		[]string{"replica", "other"},
	)

	assert.Equal(t, r, "replica")
	assert.Equal(t, v, "volume")
	assert.Equal(t, nil, err)
}

func TestDeleteReplica(t *testing.T) {

	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	replicaName := "replica"
	volumeName := "volume"
	retVolume := &lh_client.Volume{
		Name: volumeName,
	}

	// Test: GetVolume failed
	getVolumeErr := errors.New("Get volume error")
	mockClient.EXPECT().
		GetVolume(volumeName).Times(1).Return(nil, getVolumeErr)

	err := deleteReplica(replicaName, volumeName)
	assert.Equal(t, getVolumeErr, err)

	// Test: Delete failed
	deleteErr := errors.New("Delete error")
	mockClient.EXPECT().
		GetVolume(volumeName).Times(1).Return(retVolume, nil)
	mockClient.EXPECT().
		RemoveReplica(*retVolume, replicaName).
		Times(1).Return(nil, deleteErr)

	err = deleteReplica(replicaName, volumeName)
	assert.Equal(t, deleteErr, err)

	// Test: Delete passed no err
	mockClient.EXPECT().
		GetVolume(volumeName).Times(1).Return(retVolume, nil)
	mockClient.EXPECT().
		RemoveReplica(*retVolume, replicaName).
		Times(1).Return(nil, nil)

	err = deleteReplica(replicaName, volumeName)
	assert.Equal(t, nil, err)
}

func TestTimeoutWhileWaitingForDeletion(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	replicaName := "replica"
	volumeName := "volume"
	retVolume := &lh_client.Volume{
		Name: volumeName,
	}
	retVolume.Replicas = []lh_client.Replica{
		lh_client.Replica{Name: replicaName},
	}

	// Test: timeout and error
	mockClient.EXPECT().
		GetVolume(volumeName).
		AnyTimes().
		Return(retVolume, nil)

	err := waitForReplicaDeletion(replicaName, volumeName, 1)

	timeoutErr := errors.New(fmt.Sprintf(
		"timeout while deleting %s",
		replicaName,
	))
	assert.Equal(t, timeoutErr, err)
}

func TestWaitForReplicaDeletion(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	replicaName := "replica"
	volumeName := "volume"
	retVolume := &lh_client.Volume{
		Name: volumeName,
	}
	retVolume.Replicas = []lh_client.Replica{
		lh_client.Replica{Name: replicaName},
	}

	// Test: error on getting volumes
	volumeErr := errors.New("volume error")
	mockClient.EXPECT().
		GetVolume(volumeName).
		Times(1).
		Return(retVolume, volumeErr)

	err := waitForReplicaDeletion(replicaName, volumeName, 1)
	assert.Equal(t, volumeErr, err)

	// Test: replica deleted while trying

	gomock.InOrder(
		mockClient.EXPECT().
			GetVolume(volumeName).
			Times(3).
			Return(retVolume, nil),
		mockClient.EXPECT().
			GetVolume(volumeName).
			Times(1).
			Return(&lh_client.Volume{}, nil),
	)
	err = waitForReplicaDeletion(replicaName, volumeName, 5)
	assert.Equal(t, nil, err)

	// Test: replica already gone
	mockClient.EXPECT().
		GetVolume(volumeName).
		Times(1).
		Return(&lh_client.Volume{}, nil)

	err = waitForReplicaDeletion(replicaName, volumeName, 5)
	assert.Equal(t, nil, err)
}
