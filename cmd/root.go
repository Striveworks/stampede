package cmd

import (
	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "stampede",
		Short: "Quickly bootstrap k8s clusters",
		Long:  `Quickly bootstrap k8s clusters over multicast`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stampede.yaml)")
	rootCmd.PersistentFlags().StringP("cluster-type", "t", "microk8s", "type of kubernetes cluster")
	rootCmd.PersistentFlags().StringP("advertise-address", "a", "", "address for API server")
	viper.BindPFlag("cluster-type", rootCmd.PersistentFlags().Lookup("cluster-type"))
	viper.BindPFlag("advertise-address", rootCmd.PersistentFlags().Lookup("advertise-address"))

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".stampede" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".stampede")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
