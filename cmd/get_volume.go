package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var node string

// getvolumeCmd represents the getvolume command
var getvolumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Returns a list of volumes",
	Long: `Used to obtain the list of volumes currently managed by longhorn
For example:

# lhctl --url=http://10.88.1.3/v1 get volume

# lhctl --url=http://10.88.1.3/v1 get volume --node=worker-0

                   NAME                   |  STATE   | ROBUSTNESS |    SIZE     | REPLICAS |          ATTACHED TO           
------------------------------------------- ---------- ------------ ------------- ---------- --------------------------------
  pvc-343c4e3c-0090-4453-94a0-0b4a62a979c6 | attached | healthy    | 21474836480 |        3 |                  worker-0
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitManagerClient()
	},
	Run: getVolumes,
}

func init() {
	getCmd.AddCommand(getvolumeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	getvolumeCmd.PersistentFlags().StringVar(&node, "node", "", "Node that the volume is attached to")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getvolumeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getVolumes(cmd *cobra.Command, args []string) {

	out := [][]string{}
	volumes, err := mc.ListVolumes()
	if err != nil {
		eh.ExitOnError(err, "error listing volumes")
	}
	for _, volume := range volumes {
		if node == volume.Controllers[0].HostId || node == "" {
			out = append(out, []string{
				volume.Name,
				volume.State,
				volume.Robustness,
				volume.Size,
				strconv.Itoa(len(volume.Replicas)),
				volume.Controllers[0].HostId,
			})
		}

	}
	if len(out) == 0 {
		fmt.Println("No resources found")
		return
	}

	pr.PrintWithColumns(
		out,
		[]string{"Name", "State", "Robustness", "Size", "Replicas", "Attached to"},
	)

}
