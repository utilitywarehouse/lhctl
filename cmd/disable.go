package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	lh_client "github.com/longhorn/longhorn-manager/client"
)

// disableCmd represents the disable command
var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable scheduling on node",
	Long: `lhctl disable [node-name]
disables scheduling of replicas on the given node. If more than one arguments
are passed the rest will be ugnored.
example:

# lhctl disable storage-node-0
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitManagerClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		node, err := parseDisableArgs(cmd, args)
		if err != nil {
			eh.ExitOnError(err)
		}
		if err = disableNode(node); err != nil {
			eh.ExitOnError(err)
		}
		if err = waitForDisabledNode(node, 10); err != nil {
			eh.ExitOnError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(disableCmd)
}

func parseDisableArgs(cmd *cobra.Command, args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("node name not secified")
	}

	return args[0], nil
}

func disableNode(nodeId string) error {

	node, err := mc.GetNode(nodeId)
	if err != nil {
		return err
	}

	update := lh_client.Node{
		AllowScheduling: false,
	}
	_, err = mc.UpdateNode(node, update)
	return err
}

func waitForDisabledNode(nodeId string, seconds int) error {

	deadline := time.Now().Add(time.Duration(seconds) * time.Second)
	for {
		node, _ := mc.GetNode(nodeId)
		if !node.AllowScheduling {
			fmt.Println("Scheduling disabled for node:", nodeId)
			return nil
		}
		if time.Now().After(deadline) {
			return errors.New(fmt.Sprintf(
				"could not disable %s",
				nodeId,
			))
		}
	}

	return nil

}
