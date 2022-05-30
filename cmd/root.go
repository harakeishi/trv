package cmd

import (
	"os"

	butler "github.com/harakeishi/butler/butler"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "butler",
	Short: "butler",
	Long:  `butler`,
	RunE: func(cmd *cobra.Command, args []string) error {
		butler.Viewer()
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
