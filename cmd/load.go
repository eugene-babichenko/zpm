package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		updateCheck := viper.GetBool("OnLoad.CheckForUpdates")
		installMissing := viper.GetBool("OnLoad.InstallMissingPlugins")

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
		}

		// check if there are downloaded updates
		ps.checkPluginUpdates(true)

		if installMissing {
			ps.installAll()
		}

		for _, name := range ps.loadOrder {
			pse := ps.plugins[name]
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

		if !updateCheck {
			return
		}
		t, err := getLastUpdateTime()
		if err != nil {
			log.Errorf("failed to read last update time: %s", err)
			log.Error("note that this will result in extra update checks on zsh load")
		}
		checkAfter := t.Add(updateCheckPeriod)
		if t.Before(checkAfter) {
			log.Debugf("update check should be performed after %s", checkAfter.Format(time.RFC1123))
			return
		}
		if err := runUpdateCheck(); err != nil {
			log.Errorf("%s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(loadCmd)
}
