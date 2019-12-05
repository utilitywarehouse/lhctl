package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	AttachTargetNode string
	AttachVolume     string
)

// attachCmd represents the attach command
var attachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach volume from node",
	Long: `lhctl attach [volume-name] --node=[node-name]
issues an attach command for the requested volume on the specified node. If 
more than one arguments are passed the rest will be ignored.
example:

# lhctl attach my-volume --node=my-node
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitManagerClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := parseAttachArgs(cmd, args); err != nil {
			eh.ExitOnError(err)
		}
		if err := attachVolume(); err != nil {
			eh.ExitOnError(err)
		}
		if err := waitForAttachedVol(60); err != nil {
			eh.ExitOnError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(attachCmd)

	attachCmd.PersistentFlags().StringVar(
		&AttachTargetNode,
		"node",
		"",
		"(required) node to attach the volume to",
	)
}

func parseAttachArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		cmd.Help()
		return errors.New("volume name not specified")
	}

	AttachVolume = args[0]

	if AttachTargetNode == "" {
		cmd.Help()
		return errors.New("--node= flag must be set")
	}
	return nil
}

func attachVolume() error {

	vol, err := mc.GetVolume(AttachVolume)
	if err != nil {
		return err
	}

	_, err = mc.VolumeAttach(vol, AttachTargetNode)
	return err
}

func waitForAttachedVol(seconds int) error {

	deadline := time.Now().Add(time.Duration(seconds) * time.Second)
	for {
		vol, err := mc.GetVolume(AttachVolume)
		if err != nil {
			fmt.Println(err)
		}
		if vol.State == "attached" {
			fmt.Println("Successfully attached:", AttachVolume)
			return nil
		}
		if time.Now().After(deadline) {
			return errors.New(fmt.Sprintf(
				"could not attach %s before timeout",
				AttachVolume,
			))
		}
		// sleep a second to avoid hammering cpu
		time.Sleep(1 * time.Second)
	}

	return nil

}
