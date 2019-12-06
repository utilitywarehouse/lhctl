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
	err := validateDeleteReplicaArgs(deletereplicaCmd, []string{})

	assert.Equal(t, expectedErr, err)

	// Test: validate errors when list volumes errors
	listingErr := errors.New("No replica name specified")
	mockClient.EXPECT().ListVolumes().Times(1).Return(nil, listingErr)
	err = validateDeleteReplicaArgs(deletereplicaCmd, []string{"replica"})

	assert.Equal(t, listingErr, err)

	// Test: Replica not found returns error
	volumes := []lh_client.Volume{
		lh_client.Volume{},
	}
	volumes[0].Replicas = []lh_client.Replica{}
	mockClient.EXPECT().ListVolumes().Times(1).Return(volumes, nil)
	err = validateDeleteReplicaArgs(deletereplicaCmd, []string{"replica"})

	expectedErr = errors.New(fmt.Sprintf(
		"Replica not found: %s",
		ReplicaName,
	))
	assert.Equal(t, expectedErr, err)

	// Test: Replica found sets vars and return nil
	volumes[0].Replicas = append(
		volumes[0].Replicas,
		lh_client.Replica{Name: "replica"},
	)
	mockClient.EXPECT().ListVolumes().Times(1).Return(volumes, nil)
	err = validateDeleteReplicaArgs(deletereplicaCmd, []string{"replica"})

	assert.Equal(t, ReplicaName, "replica")
	assert.Equal(t, ReplicaVolume, volumes[0])
	assert.Equal(t, nil, err)

	// Test: Additional arguments ignored
	mockClient.EXPECT().ListVolumes().Times(1).Return(volumes, nil)
	err = validateDeleteReplicaArgs(
		deletereplicaCmd,
		[]string{"replica", "other"},
	)

	assert.Equal(t, ReplicaName, "replica")
	assert.Equal(t, ReplicaVolume, volumes[0])
	assert.Equal(t, nil, err)
}

func TestDeleteReplica(t *testing.T) {

	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	ReplicaName = "replica"
	ReplicaVolume = lh_client.Volume{}

	// Test: Delete failed
	deleteErr := errors.New("Delete error")
	mockClient.EXPECT().
		RemoveReplica(ReplicaVolume, ReplicaName).
		Times(1).Return(nil, deleteErr)

	err := deleteReplica()
	assert.Equal(t, deleteErr, err)

	// Test: Delete passed no err
	mockClient.EXPECT().
		RemoveReplica(ReplicaVolume, ReplicaName).
		Times(1).Return(nil, nil)

	err = deleteReplica()
	assert.Equal(t, nil, err)
}

func TestTimeoutWhileWaitingForDeletion(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	ReplicaVolume = lh_client.Volume{}
	ReplicaVolume.Name = "volume"
	ReplicaName = "replica"
	ReplicaVolume.Replicas = []lh_client.Replica{
		lh_client.Replica{Name: ReplicaName},
	}

	// Test: timeout and error
	mockClient.EXPECT().
		GetVolume(ReplicaVolume.Name).
		AnyTimes().
		Return(&ReplicaVolume, nil)

	err := waitForReplicaDeletion(1)

	timeoutErr := errors.New(fmt.Sprintf(
		"timeout while deleting %s",
		ReplicaName,
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

	ReplicaVolume = lh_client.Volume{}
	ReplicaVolume.Name = "volume"
	ReplicaName = "replica"
	ReplicaVolume.Replicas = []lh_client.Replica{
		lh_client.Replica{Name: ReplicaName},
	}

	// Test: error on getting volumes
	volumeErr := errors.New("volume error")
	mockClient.EXPECT().
		GetVolume(ReplicaVolume.Name).
		Times(1).
		Return(&ReplicaVolume, volumeErr)

	err := waitForReplicaDeletion(1)
	assert.Equal(t, volumeErr, err)

	// Test: replica deleted while trying

	gomock.InOrder(
		mockClient.EXPECT().
			GetVolume(ReplicaVolume.Name).
			Times(3).
			Return(&ReplicaVolume, nil),
		mockClient.EXPECT().
			GetVolume(ReplicaVolume.Name).
			Times(1).
			Return(&lh_client.Volume{}, nil),
	)
	err = waitForReplicaDeletion(5)
	assert.Equal(t, nil, err)

	// Test: replica already gone
	mockClient.EXPECT().
		GetVolume(ReplicaVolume.Name).
		Times(1).
		Return(&lh_client.Volume{}, nil)

	err = waitForReplicaDeletion(5)
	assert.Equal(t, nil, err)
}
