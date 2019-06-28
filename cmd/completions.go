package cmd

import (
	"github.com/eugene-babichenko/zpm/log"

	"os"

	"github.com/spf13/cobra"
)

var completionsCmd = &cobra.Command{
	Use:   "completions",
	Short: "Generate zsh command completions",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RootCmd.GenZshCompletion(os.Stdout); err != nil {
			log.Fatal("cannot generate completions: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(completionsCmd)
}
