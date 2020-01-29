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
