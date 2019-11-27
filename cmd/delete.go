package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the get command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources resources",
	Long: `lhctl delete [resource].
For example:

# lhctl delete replica`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)

}
