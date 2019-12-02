package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// detachCmd represents the detach command
var detachCmd = &cobra.Command{
	Use:   "detach",
	Short: "Detach volume from node",
	Long: `lhctl detach [volume-name]
issues a detach command for the requested volume. If more than one arguments
are passed the rest will be ignored.
example:

# lhctl detach my-volume
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitManagerClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		vol, err := parseDetachArgs(cmd, args)
		if err != nil {
			eh.ExitOnError(err)
		}
		if err = detachVolume(vol); err != nil {
			eh.ExitOnError(err)
		}
		if err = waitForDetachedVol(vol, 60); err != nil {
			eh.ExitOnError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(detachCmd)
}

func parseDetachArgs(cmd *cobra.Command, args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("volume name not specified")
	}

	return args[0], nil
}

func detachVolume(volume string) error {

	vol, err := mc.GetVolume(volume)
	if err != nil {
		return err
	}

	_, err = mc.VolumeDetach(vol)
	return err
}

func waitForDetachedVol(volume string, seconds int) error {

	deadline := time.Now().Add(time.Duration(seconds) * time.Second)
	for {
		vol, _ := mc.GetVolume(volume)
		if vol.State == "detached" {
			fmt.Println("Successfully detached:", volume)
			return nil
		}
		if time.Now().After(deadline) {
			return errors.New(fmt.Sprintf(
				"could not detach %s before timeout",
				volume,
			))
		}
	}

	return nil

}
