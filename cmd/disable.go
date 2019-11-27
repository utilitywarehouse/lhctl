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
disables scheduling of replicas on the given node.
For example:

# lhctl disable storage-node-0
`,
	Run: func(cmd *cobra.Command, args []string) {
		validateDisableArgs(cmd, args)
		disableNode(cmd, args)
		waitForDisabledNode(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(disableCmd)
}

func validateDisableArgs(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		err := errors.New("node name not secified, see lhctl disable -h")
		eh.ExitOnError(err)
	}

	nodeId := args[0]
	_, err := mc.GetNode(nodeId)
	if err != nil {
		eh.ExitOnError(err, "error getting node info")
	}
}

func disableNode(cmd *cobra.Command, args []string) {

	nodeId := args[0]

	node, _ := mc.GetNode(nodeId)

	update := lh_client.Node{
		AllowScheduling: false,
	}

	mc.UpdateNode(node, update)
}

func waitForDisabledNode(cmd *cobra.Command, args []string) {

	nodeId := args[0]

	deadline := time.Now().Add(10 * time.Second)
	for {
		node, _ := mc.GetNode(nodeId)
		if !node.AllowScheduling {
			fmt.Println("Scheduling disabled for node:", args[0])
			break
		}
		if time.Now().After(deadline) {
			err := errors.New(fmt.Sprintf(
				"could not disable %s",
				args[0],
			))
			eh.ExitOnError(err, "timeout")
		}
	}

}
