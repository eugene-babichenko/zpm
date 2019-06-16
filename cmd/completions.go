package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionsCmd = &cobra.Command{
	Use:   "completions",
	Short: "Generate zsh command completions",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RootCmd.GenZshCompletion(os.Stdout); err != nil {
			logger.Fatal("cannot generate completions: ", err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(completionsCmd)
}
