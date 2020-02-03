package cmd

import (
	"errors"
	"github.com/golang/mock/gomock"
	lh_client "github.com/longhorn/longhorn-manager/client"
	"github.com/stretchr/testify/assert"
	"github.com/utilitywarehouse/lhctl/util"
	"testing"
)

func TestParseDrainNodeArgs(t *testing.T) {

	// Test empty argument list
	args := []string{}
	node, err := parseDrainNodeArgs(drainNodeCmd, args)

	expectedErr := errors.New("No node name specified")
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, "", node)

	// Test 1 node argument passed
	args = append(args, "node")
	node, err = parseDrainNodeArgs(drainNodeCmd, args)

	assert.Equal(t, nil, err)
	assert.Equal(t, "node", node)

	// Tets more arguments will be ignored
	args = append(args, "other")
	node, err = parseDrainNodeArgs(drainNodeCmd, args)

	assert.Equal(t, nil, err)
	assert.Equal(t, "node", node)
}

func TestValidateDrainNode(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	node := "node"

	// Test no node found
	getNodeErr := errors.New("Get node error")
	mockClient.EXPECT().GetNode(node).Times(1).Return(
		nil,
		getNodeErr,
	)
	err := validateDrainNode(node)
	assert.Equal(t, getNodeErr, err)

	// Test node found
	mockClient.EXPECT().GetNode(node).Times(1).Return(
		&lh_client.Node{},
		nil,
	)
	err = validateDrainNode(node)
	assert.Equal(t, nil, err)
}

func TestFilterVolumes(t *testing.T) {

	// 2 fake nodes
	drain_node := "drain_node"
	other_node := "other_node"

	// 2 Volumes with replicas on both nodes
	vol1 := lh_client.Volume{
		Name: "vol1",
		Replicas: []lh_client.Replica{
			lh_client.Replica{
				Name:   "drain_rep_1",
				HostId: drain_node,
			},
			lh_client.Replica{
				Name:   "other_rep_1",
				HostId: other_node,
			},
		},
	}
	vol2 := lh_client.Volume{
		Name: "vol2",
		Replicas: []lh_client.Replica{
			lh_client.Replica{
				Name:   "drain_rep_2",
				HostId: drain_node,
			},
			lh_client.Replica{
				Name:   "other_rep_1",
				HostId: other_node,
			},
		},
	}
	// One with replicas on the other node
	vol3 := lh_client.Volume{
		Name: "vol3",
		Replicas: []lh_client.Replica{
			lh_client.Replica{
				Name:   "other_rep_3",
				HostId: other_node,
			},
		},
	}
	// One with no replicas
	vol4 := lh_client.Volume{
		Name: "vol4",
	}
	volume_list := []lh_client.Volume{vol1, vol2, vol3, vol4}

	// Test: Called with empty list
	filtered := filterVolumes([]lh_client.Volume{}, drain_node)
	assert.Equal(t, 0, len(filtered))

	// Test: Called with the above list will return 2 volumes
	filtered = filterVolumes(volume_list, drain_node)
	assert.Equal(t, 2, len(filtered))

	// Test: with wrong node name
	filtered = filterVolumes(volume_list, "wrong")
	assert.Equal(t, 0, len(filtered))

}
