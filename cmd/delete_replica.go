package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	lh_client "github.com/longhorn/longhorn-manager/client"
)

var (
	ReplicaName   string           // Name of replica to be deleted
	ReplicaVolume lh_client.Volume // Volume the replica belongs to
)

// getreplicaCmd represents the getreplica command
var deletereplicaCmd = &cobra.Command{
	Use:   "replica",
	Short: "Delete a replica",
	Long: `lhctl delete replica [replica-name]
For example:

# lhctl --url=http://10.88.1.3/v1 delete replica pvc-17d9c917-b44e

`,
	Run: func(cmd *cobra.Command, args []string) {
		validateDeleteReplicaArgs(cmd, args)
		deleteReplica()
		waitForReplicaDeletion()
	},
}

func init() {
	deleteCmd.AddCommand(deletereplicaCmd)
}

func validateDeleteReplicaArgs(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("You must specify a replica name")
		fmt.Println("see: `lhctl delete replica -h` for more details")
		os.Exit(1)
	}

	ReplicaName = args[0]

	volumes, err := mc.ListVolumes()
	if err != nil {
		eh.ExitOnError(err)
	}
	for _, volume := range volumes {
		for _, replica := range volume.Replicas {
			if replica.Name == ReplicaName {
				ReplicaVolume = volume
				return
			}
		}
	}

	eh.ExitOnError(errors.New("Replica not found"), ReplicaName)

}

func deleteReplica() {

	fmt.Println(fmt.Sprintf(
		"deleting replica %s for vol: %s",
		ReplicaName,
		ReplicaVolume.Name,
	))

	_, err := mc.RemoveReplica(ReplicaVolume, ReplicaName)
	if err != nil {
		eh.ExitOnError(err, "failed to delete replica")
	}
}

func waitForReplicaDeletion() {

	deadline := time.Now().Add(10 * time.Second)
	for {
		if time.Now().After(deadline) {
			err := errors.New(fmt.Sprintf(
				"could not delete %s",
				ReplicaName,
			))
			eh.ExitOnError(err, "timeout")
		}

		vol, err := mc.GetVolume(ReplicaVolume.Name)
		if err != nil {
			eh.ExitOnError(err)
		}
		for _, replica := range vol.Replicas {
			if replica.Name == ReplicaName {
				continue
			}
		}
		break
	}
}
