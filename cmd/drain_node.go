package cmd

import (
	"errors"
	"fmt"
	"sync"
	"time"

	lh_client "github.com/longhorn/longhorn-manager/client"
	"github.com/spf13/cobra"
)

var drainNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Drain replicas from a node",
	Long: `Delete all replicas from a given node. In order for a replica to
be removed from a node the following requirements should be met:
- volume robustness field is "healthy".

Drain node steps:
1. Disable scheduling on the node.
2. Iterate through the replicas that live on the node and delete in case the
above conditions are met, else wait and retry.

In case an error occurs while deleting a replica the command logs the error and
proceeds, ignoring if the replica is still on the node or not. If this command
is used to perform upgrades, longhorn will clean any failed leftover replicas
after a while.

There is no timeout in the command, so if the conditions never allow deleting a
replica the command will loop forever trying.

example:
# lhctl drain node [node-id]
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitManagerClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		node, err := parseDrainNodeArgs(cmd, args)
		if err != nil {
			eh.ExitOnError(err)
		}

		if err := validateDrainNode(node); err != nil {
			eh.ExitOnError(err)
		}

		// Call disable node from the respective command disableCmd
		if err := disableNode(node); err != nil {
			eh.ExitOnError(err)
		}
		if err := waitForDisabledNode(node, 10); err != nil {
			eh.ExitOnError(err)
		}

		// Drain the node from replicas
		if err := drainNode(node); err != nil {
			eh.ExitOnError(err)
		}
	},
}

func init() {
	drainCmd.AddCommand(drainNodeCmd)
}

func parseDrainNodeArgs(cmd *cobra.Command, args []string) (string, error) {
	if len(args) < 1 {
		cmd.Help()
		return "", errors.New("No node name specified")
	}

	node := args[0]

	if len(args) > 1 {
		fmt.Println(fmt.Sprintf(
			"ignoring arguments: %v",
			args[1:],
		))
	}

	return node, nil
}

func validateDrainNode(node string) error {
	_, err := mc.GetNode(node)
	return err
}

func drainNode(nodeName string) error {

	volumes, err := mc.ListVolumes()
	if err != nil {
		return err
	}

	filteredVolumes := filterVolumes(volumes, nodeName)
	if len(filteredVolumes) == 0 {
		fmt.Println("Nothing to drain")
		return nil
	}

	var wg sync.WaitGroup

	// Iterate through replica list and drain on best effort. If we receive
	// an error while issuing the remove command or waiting to verify
	// deletion just log and proceed.
	for _, volume := range filteredVolumes {
		for _, replica := range volume.Replicas {
			if replica.HostId == nodeName {
				wg.Add(1)
				go func(replicaName, volumeName string) {
					defer wg.Done()
					for {
						err := drainReplicaFromVolume(
							replicaName,
							volumeName,
						)
						if err == nil {
							fmt.Println(fmt.Sprintf(
								"%v deleted",
								replicaName,
							))
							return
						}
						fmt.Println(fmt.Sprintf(
							"error %v deleting replica: %s, retrying",
							err,
							replicaName,
						))
						// sleep before retrying
						time.Sleep(time.Duration(1 * time.Second))
					}
				}(replica.Name, volume.Name)
			}
		}
	}

	wg.Wait()
	return nil

}

func drainReplicaFromVolume(replicaName, volumeName string) error {
	// waitForVolumeHealthy could run forever until volume robustness become
	// healthy
	waitForVolumeHealthy(volumeName)

	if err := deleteReplica(replicaName, volumeName); err != nil {
		return err
	}

	if err := waitForReplicaDeletion(replicaName, volumeName, 10); err != nil {
		return err
	}
	return nil
}

func filterVolumes(volumes []lh_client.Volume, nodeName string) []lh_client.Volume {
	var filteredVolumes []lh_client.Volume

	for _, volume := range volumes {
		for _, replica := range volume.Replicas {
			if replica.HostId == nodeName {
				filteredVolumes = append(filteredVolumes, volume)
			}
		}
	}

	return filteredVolumes
}

func waitForVolumeHealthy(volumeName string) {
	for {
		volume, err := mc.GetVolume(volumeName)
		if err != nil {
			fmt.Println(fmt.Sprintf("%v", err))
		} else {
			if volume.Robustness == "healthy" {
				return
			}
		}
		fmt.Println(fmt.Sprintf(
			"waiting for volume %s to become healthy, sleeping 10 seconds",
			volumeName,
		))
		time.Sleep(time.Duration(10 * time.Second))
	}
}
