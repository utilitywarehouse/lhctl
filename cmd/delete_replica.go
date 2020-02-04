package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// deletereplicaCmd is the command to delete a replica
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
		replica, volume, err := validateDeleteReplicaArgs(cmd, args)
		if err != nil {
			eh.ExitOnError(err)
		}
		if err := deleteReplica(replica, volume); err != nil {
			eh.ExitOnError(err)
		}
		if err := waitForReplicaDeletion(replica, volume, 10); err != nil {
			eh.ExitOnError(err)
		}
	},
}

func init() {
	deleteCmd.AddCommand(deletereplicaCmd)
}

func validateDeleteReplicaArgs(cmd *cobra.Command, args []string) (string, string, error) {
	if len(args) < 1 {
		cmd.Help()
		return "", "", errors.New("No replica name specified")
	}

	replicaName := args[0]

	volumes, err := mc.ListVolumes()
	if err != nil {
		return "", "", err
	}
	for _, volume := range volumes {
		for _, replica := range volume.Replicas {
			if replica.Name == replicaName {
				return replicaName, volume.Name, nil
			}
		}
	}

	return "", "", errors.New(fmt.Sprintf(
		"Replica not found: %s",
		replicaName,
	))

}

func deleteReplica(replicaName, volumeName string) error {

	fmt.Println(fmt.Sprintf(
		"deleting replica %s for vol: %s",
		replicaName,
		volumeName,
	))

	volume, err := mc.GetVolume(volumeName)
	if err != nil {
		return err
	}

	_, err = mc.RemoveReplica(*volume, replicaName)
	if err != nil {
		return err
	}

	return nil
}

func waitForReplicaDeletion(replicaName, volumeName string, seconds int) error {

	deadline := time.Now().Add(time.Duration(seconds) * time.Second)
	for {
		if time.Now().After(deadline) {
			err := errors.New(fmt.Sprintf(
				"timeout while deleting %s",
				replicaName,
			))
			return err
		}

		vol, err := mc.GetVolume(volumeName)
		if err != nil {
			return err
		}
		replicaFound := false
		for _, replica := range vol.Replicas {
			if replica.Name == replicaName {
				replicaFound = true
			}
		}
		if !replicaFound {
			return nil
		}
		// sleep a second to avoid hammering cpu
		time.Sleep(1 * time.Second)
	}

	return nil
}
