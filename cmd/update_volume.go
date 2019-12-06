package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	UpdateVolumeName string
	ReplicaCountFlag string
	ReplicaCount     int64
)

// updateVolumeCmd represents the update volume command
var updateVolumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Updates a volume resource",
	Long: `lhctl update volume [volume-name] [flags]
Used to update a volume resource attributes. At least one of the available 
command flags should be set.
For example:

# lhctl --url=http://10.88.1.3/v1 update volume [volume-name] --replicas=2
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitManagerClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := parseUpdateVolumeArgs(cmd, args); err != nil {
			eh.ExitOnError(err)
		}
		if err := updateVolume(); err != nil {
			eh.ExitOnError(err)
		}
	},
}

func init() {
	updateCmd.AddCommand(updateVolumeCmd)

	updateVolumeCmd.PersistentFlags().StringVar(
		&ReplicaCountFlag,
		"replicas",
		"",
		"Replica count, 0 not allowed as it destroys the volume",
	)
}

func parseUpdateVolumeArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		cmd.Help()
		return errors.New("volume name not specified")
	}

	UpdateVolumeName = args[0]

	// Check that at least on flag is set
	if ReplicaCountFlag == "" {
		cmd.Help()
		return errors.New("at least one flag should be set")
	}

	// ReplicaCount
	count, err := strconv.ParseInt(ReplicaCountFlag, 10, 64)
	ReplicaCount = count
	if err != nil {
		return err
	}
	if ReplicaCount <= 0 {
		cmd.Help()
		return errors.New("--replicas= flag invalid value")
	}

	return nil
}

func updateVolume() error {

	vol, err := mc.GetVolume(UpdateVolumeName)
	if err != nil {
		return err
	}

	// Update replica count
	if ReplicaCount > 0 {
		_, err := mc.UpdateReplicaCount(vol, ReplicaCount)
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf(
			"updated %s replica count: %d",
			UpdateVolumeName,
			ReplicaCount,
		))
	}

	return nil

}
