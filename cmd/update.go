package cmd

import (
	"github.com/spf13/cobra"
)

var validUpdateArgs = []string{"volume"}
var argUpdateAliases = []string{"volumes"}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update resources",
	Long: `Perform updates on the available resources.
For example:

# lhctl update volume [volume-name] --replicas=2`,
	ValidArgs:  validUpdateArgs,
	ArgAliases: argUpdateAliases,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
