package cmd

import (
	"github.com/spf13/cobra"
)

var validArgs = []string{"node", "replica", "volume"}
var argAliases = []string{"nodes", "replicas", "volumes"}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "List resources",
	Long: `Get details for the available resources.
For example:

# lhctl get volume`,
	ValidArgs:  validArgs,
	ArgAliases: argAliases,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
