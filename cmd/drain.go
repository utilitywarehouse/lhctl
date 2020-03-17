package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var (
	DrainTargetNode string
)

var drainCmd = &cobra.Command{
	Use:   "drain",
	Short: "Drain replicas from a node",
	Long: `lhctl drain --node=[node-name]
disables scheduling on the node and iterates through its
current replicas deleting those that belong to healthy volumes
and waiting for that are not healthy to become healthy:

# lhctl drain --node=my-node
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitManagerClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := parseDrainArgs(cmd, args); err != nil {
			eh.ExitOnError(err)
		}
		if err := disableNode(DrainTargetNode); err != nil {
			eh.ExitOnError(err)
		}
		if err := deleteNodeReplcias(DrainTargetNode); err != nil {
			eh.ExitOnError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(drainCmd)

	drainCmd.PersistentFlags().StringVar(
		&DrainTargetNode,
		"node",
		"",
		"(required) node to drain replicas from",
	)
}

func parseDrainArgs(cmd *cobra.Command, args []string) error {
	if DrainTargetNode == "" {
		cmd.Help()
		return errors.New("--node= flag must be set")
	}
	return nil
}

func deleteNodeReplcias(nodeId string) error {
	volumes, err := mc.ListVolumes()
	if err != nil {
		return err
	}

	for _, volume := range volumes {
		for _, replica := range volume.Replicas {
			if replica.HostId == nodeId {
				println(fmt.Sprintf("removing replica %s for pvc %s", replica.Name, volume.KubernetesStatus.PvcName))

				err := waitForVolumeHealthy(volume.Id)
				if err != nil {
					return err
				}

				_, err = mc.RemoveReplica(volume, replica.Name)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func waitForVolumeHealthy(volumeId string) error {
	for {
		time.Sleep(time.Second * 5)

		volume, err := mc.GetVolume(volumeId)
		if err != nil {
			return err
		}

		if volume.Robustness == "healthy" {
			return nil
		}

		println(fmt.Sprintf("pvc %s currently %s waiting for it to become healthy", volume.KubernetesStatus.PvcName, volume.Robustness))
	}
}
