package cmd

import (
	"github.com/spf13/cobra"
	"striveworks.us/stampede/pkg"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize stampede",
	Long:  `Initialize stampede`,
	Run: func(cmd *cobra.Command, args []string) {
		node := pkg.CreateNode()
		node.Start()
	},
}
