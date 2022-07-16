package cmd

import (
	"os"

	trv "github.com/harakeishi/trv/trv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "trv",
	Short: "trv",
	Long:  `trv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var trv trv.Trv
		trv.Init()
		trv.Draw()
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
