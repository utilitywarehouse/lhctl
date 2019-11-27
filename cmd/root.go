package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	//homedir "github.com/mitchellh/go-homedir"
	//"github.com/spf13/viper"

	"github.com/utilitywarehouse/lhctl/util"
)

var cfgFile string
var url string

// root manager client
var mc *util.ManagerClient

// root error handler
var eh *util.ErrorHandler

// print output interface
var pr *util.Printer

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lhctl",
	Short: "CLI for longhorn",
	Long: `CLI to perform basic actions against longhorn api to help with
automated tasks.
For example:

# lhctl get volume

# lhctl get node`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		eh.ExitOnError(err)
	}

}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initManagerClient)

	// Placeholder for config file
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lhctl.yaml)")
	// root persistent flag `url` will be available everywhere in the package
	rootCmd.PersistentFlags().StringVar(&url, "url", "", "longhorn manager url (example: http://10.88.1.3/v1)")
	rootCmd.MarkFlagRequired("url")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// init error handler
	eh = &util.ErrorHandler{}

	pr = &util.Printer{}
}

func initManagerClient() {
	// Get a new manager client
	client, err := util.NewManagerClient(url)
	if err != nil {
		eh.ExitOnError(err)
	}
	mc = client
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	//if cfgFile != "" {
	//	// Use config file from the flag.
	//	viper.SetConfigFile(cfgFile)
	//} else {
	//	// Find home directory.
	//	home, err := homedir.Dir()
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}

	//	// Search config in home directory with name ".lhctl" (without extension).
	//	viper.AddConfigPath(home)
	//	viper.SetConfigName(".lhctl")
	//}
	//viper.AutomaticEnv() // read in environment variables that match

	//// If a config file is found, read it in.
	//if err := viper.ReadInConfig(); err == nil {
	//	fmt.Println("Using config file:", viper.ConfigFileUsed())
	//}

	if url == "" {
		err := errors.New("You need to provide a url using `--url=` flag")
		eh.ExitOnError(err)
	}
}
