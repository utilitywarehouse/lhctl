package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	lh_client "github.com/longhorn/longhorn-manager/client"
)

//var EnableNode lh_client.Node

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable scheduling on node",
	Long: `lhctl enable [node-name]
enables scheduling of replicas on the given node. If more than one arguments
are passed the rest will be ugnored.
example:

# lhctl enable storage-node-0
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitManagerClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		node, err := parseEnableArgs(cmd, args)
		if err != nil {
			eh.ExitOnError(err)
		}
		if err = enableNode(node); err != nil {
			eh.ExitOnError(err)
		}
		if err = waitForEnabledNode(node, 10); err != nil {
			eh.ExitOnError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)
}

func parseEnableArgs(cmd *cobra.Command, args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("node name not secified")
	}

	return args[0], nil
}

func enableNode(nodeId string) error {

	node, err := mc.GetNode(nodeId)
	if err != nil {
		return err
	}

	update := lh_client.Node{
		AllowScheduling: true,
	}

	_, err = mc.UpdateNode(node, update)
	return err
}

func waitForEnabledNode(nodeId string, seconds int) error {

	deadline := time.Now().Add(time.Duration(seconds) * time.Second)
	for {
		node, _ := mc.GetNode(nodeId)
		if node.AllowScheduling {
			fmt.Println("Scheduling enabled for node:", nodeId)
			return nil
		}
		if time.Now().After(deadline) {
			return errors.New(fmt.Sprintf(
				"could not enable %s",
				nodeId,
			))
		}
	}

	return nil

}
