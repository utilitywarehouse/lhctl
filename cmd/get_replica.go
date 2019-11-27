package cmd

import (
	"fmt"
	"strconv"

	lh_client "github.com/longhorn/longhorn-manager/client"
	"github.com/spf13/cobra"
)

var volume string

// getreplicaCmd represents the getreplica command
var getreplicaCmd = &cobra.Command{
	Use:   "replica",
	Short: "Returns a list of replicas",
	Long: `Gives a list of replicas, can filter via volume.
For example:

# lhctl --url=http://10.88.1.3/v1 get replica

# lhctl --url=http://10.88.1.3/v1 get replica --volume=pvc-17d9c917-b44e-4cac-9f99-6f071ecfb8b9
                         NAME                         | MODE | RUNNING |               HOST ID               |                                  DATA PATH                                   
------------------------------------------------------ ------ --------- ------------------------------------- ------------------------------------------------------------------------------
  pvc-17d9c917-b44e-4cac-9f99-6f071ecfb8b9-r-1a61738b | RW   | true    | storage-node-0.dev.merit.uw.systems | /var/lib/storage/replicas/pvc-17d9c917-b44e-4cac-9f99-6f071ecfb8b9-ac508aad  
  pvc-17d9c917-b44e-4cac-9f99-6f071ecfb8b9-r-296102e6 | RW   | true    | storage-node-2.dev.merit.uw.systems | /var/lib/storage/replicas/pvc-17d9c917-b44e-4cac-9f99-6f071ecfb8b9-1f957794  
  pvc-17d9c917-b44e-4cac-9f99-6f071ecfb8b9-r-9a599398 | RW   | true    | storage-node-1.dev.merit.uw.systems | /var/lib/storage/replicas/pvc-17d9c917-b44e-4cac-9f99-6f071ecfb8b9-7dc909a4 
`,
	Run: getReplicas,
}

func init() {
	getCmd.AddCommand(getreplicaCmd)

	getreplicaCmd.Flags().StringVar(&volume, "volume", "", "Volume that the replica belongs to")
}

func getReplicas(cmd *cobra.Command, args []string) {

	replicas := []lh_client.Replica{}
	// if volume specified get its replicas
	if volume != "" {
		vol, err := mc.GetVolume(volume)
		if err != nil {
			eh.ExitOnError(err)
		}
		replicas = vol.Replicas

	} else {
		// else iterate through all and get replicas
		vols, err := mc.ListVolumes()
		if err != nil {
			eh.ExitOnError(err)
		}
		for _, vol := range vols {
			for _, rep := range vol.Replicas {
				replicas = append(replicas, rep)
			}
		}
	}

	out := [][]string{}
	for _, replica := range replicas {
		out = append(out, []string{
			replica.Name,
			replica.Mode,
			strconv.FormatBool(replica.Running),
			replica.HostId,
			replica.DataPath,
		})
	}

	if len(out) == 0 {
		fmt.Println("No resources found")
		return
	}

	pr.PrintWithColumns(
		out,
		[]string{"Name", "Mode", "Running", "Host Id", "Data Path"},
	)
}
