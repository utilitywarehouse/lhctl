package util

import (
	"github.com/longhorn/longhorn-manager/client"
)

type ManagerClient struct {
	rancherClient *client.RancherClient
}

type ManagerClientInterface interface {
	ListVolumes() ([]client.Volume, error)
	GetVolume(string) (*client.Volume, error)
	ListNodes() ([]client.Node, error)
	GetNode(string) (*client.Node, error)
	UpdateNode(*client.Node, interface{}) (*client.Node, error)
	RemoveReplica(client.Volume, string) (*client.Volume, error)
	VolumeDetach(*client.Volume) (*client.Volume, error)
	VolumeAttach(*client.Volume, string) (*client.Volume, error)
	UpdateReplicaCount(*client.Volume, int64) (*client.Volume, error)
}

func NewManagerClient(url, user, pass string) (*ManagerClient, error) {
	rc, err := client.NewRancherClient(&client.ClientOpts{
		Url:       url,
		AccessKey: user,
		SecretKey: pass,
	})
	if err != nil {
		return nil, err
	}

	return &ManagerClient{rancherClient: rc}, nil
}

func (mc ManagerClient) ListVolumes() ([]client.Volume, error) {
	collection, err := mc.rancherClient.Volume.List(client.NewListOpts())
	return collection.Data, err
}

func (mc ManagerClient) GetVolume(id string) (*client.Volume, error) {
	return mc.rancherClient.Volume.ById(id)
}

func (mc ManagerClient) ListNodes() ([]client.Node, error) {
	collection, err := mc.rancherClient.Node.List(client.NewListOpts())
	return collection.Data, err
}

func (mc ManagerClient) GetNode(nodeId string) (*client.Node, error) {
	return mc.rancherClient.Node.ById(nodeId)
}

func (mc ManagerClient) UpdateNode(node *client.Node, update interface{}) (*client.Node, error) {
	return mc.rancherClient.Node.Update(node, update)
}

func (mc ManagerClient) RemoveReplica(volume client.Volume, replicaName string) (*client.Volume, error) {
	return mc.rancherClient.Volume.ActionReplicaRemove(
		&client.Volume{
			Resource: volume.Resource,
		},
		&client.ReplicaRemoveInput{
			Name: replicaName,
		},
	)
}

func (mc ManagerClient) VolumeDetach(volume *client.Volume) (*client.Volume, error) {
	return mc.rancherClient.Volume.ActionDetach(volume)
}

func (mc ManagerClient) VolumeAttach(volume *client.Volume, hostId string) (*client.Volume, error) {
	input := &client.AttachInput{
		HostId: hostId,
	}
	return mc.rancherClient.Volume.ActionAttach(volume, input)
}

func (mc ManagerClient) UpdateReplicaCount(volume *client.Volume, replicaCount int64) (*client.Volume, error) {
	// ActionUpdateReplicaCount is missing from generated client volume
	// files, so let's rewrite it here.
	input := &client.UpdateReplicaCountInput{
		ReplicaCount: replicaCount,
	}
	resp := &client.Volume{}
	VOLUME_TYPE := "volume"
	err := mc.rancherClient.Action(VOLUME_TYPE, "updateReplicaCount", &volume.Resource, input, resp)

	return resp, err
}
