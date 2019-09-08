package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var completionsCmd = &cobra.Command{
	Use:   "completions",
	Short: "Generate zsh command completions",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RootCmd.GenZshCompletion(os.Stdout); err != nil {
			log.Fatalf("cannot generate completions: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(completionsCmd)
}
