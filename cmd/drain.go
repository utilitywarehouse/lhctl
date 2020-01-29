package cmd

import (
	"github.com/spf13/cobra"
)

var validDrainArgs = []string{"node"}

// getCmd represents the get command
var drainCmd = &cobra.Command{
	Use:   "drain",
	Short: "Drain replicas",
	Long: `Drain resources out of replicas.
Availbale resource types [node]

example:
# lhctl drain node [node-id]`,
	ValidArgs: validDrainArgs,
}

func init() {
	rootCmd.AddCommand(drainCmd)
}
