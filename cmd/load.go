package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func runUpdateCheck() error {
	if err := exec.Command("zpm", "check").Start(); err != nil {
		return errors.Wrap(err, "while running background update check")
	}

	return nil
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load configured plugins into the current shell",
	Run: func(cmd *cobra.Command, args []string) {
		updateCheck, _ := cmd.Flags().GetBool("update-check")
		installMissing, _ := cmd.Flags().GetBool("install-missing")

		// TODO print out meta data with hints about updates

		fpath := make([]string, 0)
		exec := make([]string, 0)

		lines := []string{
			// Use different compdump for different zsh versions (kindly
			// borrowed from Oh My Zsh).
			"ZSH_COMPDUMP=\"${ZDOTDIR:-${HOME}}/.zcompdump-${SHORT_HOST}-${ZSH_VERSION}\"",
			"autoload -U compaudit compinit",
		}

		ps, err := makePluginStorage(rootDir, pluginsSpecs)
		if err != nil {
			log.Fatalf("while reading plugin configurations: %s", err)
			return
		}

		if installMissing {
			ps.checkPluginInstalls()
			ps.installAll()
		}

		for _, pse := range ps.plugins {
			fpathPlugin, execPlugin, err := pse.plugin.Load()
			if err != nil {
				log.Errorf("while loading plugin %s: %s", pse.name, err)
				continue
			}
			fpath = append(fpath, fpathPlugin...)
			exec = append(exec, execPlugin...)
		}

		var fpathBuilder strings.Builder
		_, _ = fmt.Fprint(&fpathBuilder, "fpath=(")
		for _, fpathEntry := range fpath {
			_, _ = fmt.Fprintf(&fpathBuilder, "%s ", fpathEntry)
		}
		_, _ = fmt.Fprint(&fpathBuilder, "$fpath)")

		// Load completions from `fpath`.
		lines = append(lines, fpathBuilder.String(), "compinit -u -C -d \"${ZSH_COMPDUMP}\"")
		lines = append(lines, exec...)

		for _, line := range lines {
			fmt.Println(line)
		}

		if updateCheck {
			if err := runUpdateCheck(); err != nil {
				log.Errorf("%s", err)
			}
		}
	},
}

func init() {
	loadCmd.Flags().Bool(
		"update-check",
		false,
		"Check for updates once in a period defined in the settings (default: 24h).",
	)

	loadCmd.Flags().Bool(
		"install-missing",
		false,
		"Install plugins that are listed in the configuration but are not installed.",
	)

	RootCmd.AddCommand(loadCmd)
}
