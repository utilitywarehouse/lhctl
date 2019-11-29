package cmd

import (
	"errors"
	"fmt"
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
Accepts one replica name, rest of the arguments will be ignored.
Example:

# lhctl --url=http://10.88.1.3/v1 delete replica pvc-17d9c917-b44e

`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitManagerClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateDeleteReplicaArgs(cmd, args); err != nil {
			eh.ExitOnError(err)
		}
		if err := deleteReplica(); err != nil {
			eh.ExitOnError(err)
		}
		if err := waitForReplicaDeletion(10); err != nil {
			eh.ExitOnError(err)
		}
	},
}

func init() {
	deleteCmd.AddCommand(deletereplicaCmd)
}

func validateDeleteReplicaArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		cmd.Help()
		return errors.New("No replica name specified")
	}

	ReplicaName = args[0]

	volumes, err := mc.ListVolumes()
	if err != nil {
		return err
	}
	for _, volume := range volumes {
		for _, replica := range volume.Replicas {
			if replica.Name == ReplicaName {
				ReplicaVolume = volume
				return nil
			}
		}
	}

	return errors.New(fmt.Sprintf(
		"Replica not found: %s",
		ReplicaName,
	))

}

func deleteReplica() error {

	fmt.Println(fmt.Sprintf(
		"deleting replica %s for vol: %s",
		ReplicaName,
		ReplicaVolume.Name,
	))

	_, err := mc.RemoveReplica(ReplicaVolume, ReplicaName)
	if err != nil {
		return err
	}

	return nil
}

func waitForReplicaDeletion(seconds int) error {

	deadline := time.Now().Add(time.Duration(seconds) * time.Second)
	for {
		if time.Now().After(deadline) {
			err := errors.New(fmt.Sprintf(
				"timeout while deleting %s",
				ReplicaName,
			))
			return err
		}

		vol, err := mc.GetVolume(ReplicaVolume.Name)
		if err != nil {
			return err
		}
		replicaFound := false
		for _, replica := range vol.Replicas {
			if replica.Name == ReplicaName {
				replicaFound = true
			}
		}
		if !replicaFound {
			return nil
		}
	}

	return nil
}
