package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var getnodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Returns nodes information",
	Long: `View the nodes' longhorn attributes

example:
# lhctl get node
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitManagerClient()
	},
	Run: getNodes,
}

func init() {
	getCmd.AddCommand(getnodeCmd)

}

func getNodes(cmd *cobra.Command, args []string) {

	out := [][]string{}
	nodes, err := mc.ListNodes()
	if err != nil {
		eh.ExitOnError(err, "Failed to list nodes")
	}
	for _, node := range nodes {
		readiness := node.Conditions["Ready"].(map[string]interface{})
		diskPaths := []string{}
		for _, disk := range node.Disks {
			disk := disk.(map[string]interface{})
			diskPaths = append(diskPaths, disk["path"].(string))
		}
		out = append(out, []string{
			node.Name,
			strconv.FormatBool(node.AllowScheduling),
			readiness["status"].(string),
			strings.Join(diskPaths, ","),
		})
	}
	if len(out) == 0 {
		fmt.Println("No resources found")
		return
	}

	pr.PrintWithColumns(
		out,
		[]string{"Name", "AllowScheduling", "Ready", "Disks"},
	)

}
