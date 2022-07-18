package cmd

import (
	"os"

	curver "github.com/harakeishi/curver"
	trv "github.com/harakeishi/trv/trv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "trv",
	Short: "`trv` is a remote viewer for tbls",
	Long:  "`trv` is a remote viewer for tbls",
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := cmd.Flags().GetBool("version")
		if err != nil {
			return err
		}
		if version {
			curver.EchoVersion()
			return nil
		}
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
	rootCmd.Flags().BoolP("version", "V", false, "echo version")
}
