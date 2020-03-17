package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/utilitywarehouse/lhctl/util"
)

// Flags
var (
	cfgFile     string
	urlFlag     string
	userFlag    string
	passFlag    string
	contextFlag string
)

// root manager client
var mc util.ManagerClientInterface

// root error handler
var eh util.ErrorHandlerInterface

// print output interface
var pr util.PrinterInterface

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lhctl",
	Short: "CLI for longhorn",
	Long: `CLI to perform basic actions against longhorn api to help with
automated tasks.
For example:

# lhctl get volume

# lhctl get node`,
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

	// Placeholder for config file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lhctl.yaml)")
	// root persistent flag `url` will be available everywhere in the package
	rootCmd.PersistentFlags().StringVar(&urlFlag, "url", "", "longhorn manager url (example: http://10.88.1.3/v1)")
	rootCmd.PersistentFlags().StringVar(&userFlag, "user", "", "user for http request")
	rootCmd.PersistentFlags().StringVar(&passFlag, "pass", "", "password for http request")
	rootCmd.PersistentFlags().StringVar(&contextFlag, "context", "", "config file context")

	// init error handler
	eh = &util.ErrorHandler{}

	pr = &util.Printer{}
}

type Context struct {
	Name string `mapstructure:"name"`
	Url  string `mapstructure:"url"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
}

type Config struct {
	Contexts       []Context
	DefaultContext string `mapstructure:"default"`
}

func readConfig() (Config, error) {

	var C Config
	err := viper.Unmarshal(&C)
	if err != nil {
		return Config{}, err
	}
	return C, nil

}

func getClientParams() (string, string, string, error) {
	config, err := readConfig()
	if err != nil {
		return "", "", "", err
	}

	// Get active context using default and overriding with flag
	activeContext := config.DefaultContext
	if contextFlag != "" {
		activeContext = contextFlag
	}

	// Get config from context if any
	var clientUrl, clientUser, clientPass string
	if activeContext != "" {
		for _, c := range config.Contexts {
			if c.Name == activeContext {
				clientUrl = c.Url
				clientUser = c.User
				clientPass = c.Pass
				break
			}
		}
	}

	// Override url via flag
	if urlFlag != "" {
		clientUrl = urlFlag
	}

	if clientUrl == "" {
		err := errors.New(
			"You need to provide a url via config or using `--url=` flag",
		)
		return "", "", "", err
	}

	// Override user via flag
	if userFlag != "" {
		clientUser = userFlag
	}

	// Override pass via flag
	if passFlag != "" {
		clientPass = passFlag
	}

	return clientUrl, clientUser, clientPass, nil
}

// InitManagerClient: call it before running commangs that need to interact
// with the longhorn manager api
func InitManagerClient() {
	clientUrl, clientUser, clientPass, err := getClientParams()
	if err != nil {
		eh.ExitOnError(err)
	}

	// Get a new manager client
	client, err := util.NewManagerClient(clientUrl, clientUser, clientPass)
	if err != nil {
		eh.ExitOnError(err)
	}
	mc = client
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetConfigType("yaml")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".lhctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".lhctl")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Print(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	}
}
