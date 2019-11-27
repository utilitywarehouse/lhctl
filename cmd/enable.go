package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	lh_client "github.com/longhorn/longhorn-manager/client"
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable scheduling on node",
	Long: `lhctl enable [node-name]
enables scheduling of replicas on the given node.
For example:

# lhctl enable storage-node-0
`,
	Run: func(cmd *cobra.Command, args []string) {
		validateEnableArgs(cmd, args)
		enableNode(cmd, args)
		waitForEnabledNode(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)
}

func validateEnableArgs(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		err := errors.New("node name not secified, see lhctl enable -h")
		eh.ExitOnError(err)
	}

	nodeId := args[0]
	_, err := mc.GetNode(nodeId)
	if err != nil {
		eh.ExitOnError(err, "error getting node info")
	}
}

func enableNode(cmd *cobra.Command, args []string) {

	nodeId := args[0]

	node, _ := mc.GetNode(nodeId)

	update := lh_client.Node{
		AllowScheduling: true,
	}

	mc.UpdateNode(node, update)
}

func waitForEnabledNode(cmd *cobra.Command, args []string) {

	nodeId := args[0]

	deadline := time.Now().Add(10 * time.Second)
	for {
		node, _ := mc.GetNode(nodeId)
		if node.AllowScheduling {
			fmt.Println("Scheduling enabled for node:", args[0])
			break
		}
		if time.Now().After(deadline) {
			err := errors.New(fmt.Sprintf(
				"could not enable %s",
				args[0],
			))
			eh.ExitOnError(err, "timeout")
		}
	}

}
