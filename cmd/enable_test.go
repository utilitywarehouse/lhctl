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

func TestParseEnableArgs(t *testing.T) {

	// Test: parse with no args returns error
	expectedErr := errors.New("node name not secified")
	_, err := parseEnableArgs(enableCmd, []string{})

	assert.Equal(t, expectedErr, err)

	// Test: first argument returned as node id
	nodeId, err := parseEnableArgs(enableCmd, []string{"node"})

	assert.Equal(t, nil, err)
	assert.Equal(t, nodeId, "node")

	// Test: Rest of arguments ignored
	nodeId, err = parseEnableArgs(enableCmd, []string{"node", "other"})

	assert.Equal(t, nil, err)
	assert.Equal(t, nodeId, "node")
}

func TestEnableNode(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	// Test: Error when get node errors
	nodeErr := errors.New("node error")
	mockClient.EXPECT().GetNode("node").Times(1).Return(nil, nodeErr)
	err := enableNode("node")

	assert.Equal(t, nodeErr, err)

	// Test: Error when update node errors
	node := &lh_client.Node{}
	mockClient.EXPECT().GetNode("node").Times(1).Return(node, nil)

	update := lh_client.Node{
		AllowScheduling: true,
	}
	updErr := errors.New("update error")
	mockClient.EXPECT().UpdateNode(node, update).
		Times(1).Return(nil, updErr)

	err = enableNode("node")
	assert.Equal(t, updErr, err)

	// Test: nil when no call errors
	mockClient.EXPECT().GetNode("node").Times(1).Return(node, nil)
	mockClient.EXPECT().UpdateNode(node, update).
		Times(1).Return(node, nil)

	err = enableNode("node")
	assert.Equal(t, nil, err)
}

func TestTimeoutWhileWaitingForEnabledNode(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	// Test: timeout returns error
	node := &lh_client.Node{}
	node.AllowScheduling = false
	mockClient.EXPECT().GetNode("node").AnyTimes().Return(node, nil)

	timeoutErr := errors.New(fmt.Sprintf(
		"could not enable %s",
		"node",
	))

	err := waitForEnabledNode("node", 1)
	assert.Equal(t, timeoutErr, err)

	// Test: get node error will not exit until timeout
	node = &lh_client.Node{}
	nodeErr := errors.New("node error")
	mockClient.EXPECT().GetNode("node").AnyTimes().Return(node, nodeErr)

	err = waitForEnabledNode("node", 1)
	assert.Equal(t, timeoutErr, err)
}

func TestWaitForEnabledNode(t *testing.T) {
	// Mock the manager client
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := util.NewMockManagerClientInterface(mockCtrl)
	// override the package scoped manager client
	mc = mockClient

	// Test: node enabled while trying
	disabled := &lh_client.Node{}
	disabled.AllowScheduling = false
	enabled := &lh_client.Node{}
	enabled.AllowScheduling = true

	gomock.InOrder(
		mockClient.EXPECT().GetNode("node").Times(3).
			Return(disabled, nil),
		mockClient.EXPECT().GetNode("node").Times(1).
			Return(enabled, nil),
	)

	err := waitForEnabledNode("node", 5)
	assert.Equal(t, nil, err)

	// Test: already enabled
	mockClient.EXPECT().GetNode("node").Times(1).Return(enabled, nil)

	err = waitForEnabledNode("node", 5)
	assert.Equal(t, nil, err)
}
