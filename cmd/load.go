package cmd

import (
	"github.com/eugene-babichenko/zpm/plugin"

	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getLastUpdateTime() (t time.Time, err error) {
	filename := filepath.Join(rootDir, ".lastupdate")
	data, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		return t, nil
	}
	if err != nil {
		return t, errors.Wrap(err, "while reading last update time")
	}
	t, err = time.Parse(time.RFC3339, string(data))
	if err != nil {
		return t, errors.Wrap(err, "failed to parse time")
	}
	return t, nil
}

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

		ps, err := plugin.MakePluginStorage(rootDir, pluginsSpecs)
		if err != nil {
			log.Fatalf("while reading plugin configurations: %s", err)
		}

		// check if there are downloaded updates
		ps.CheckPluginUpdates(true)

		if installMissing {
			ps.InstallAll()
		}

		for _, name := range ps.LoadOrder {
			pse := ps.Plugins[name]
			fpathPlugin, execPlugin, err := pse.Plugin.Load()
			if err != nil {
				log.Errorf("while loading plugin %s: %s", pse.Name, err)
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
