package cmd

import (
	"github.com/spf13/cobra"
)

var validGetArgs = []string{"node", "replica", "volume"}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "List resources",
	Long: `Get details for the available resources.
For example:

# lhctl get volume`,
	ValidArgs: validGetArgs,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
