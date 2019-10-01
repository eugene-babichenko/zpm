package commands

import (
	"github.com/eugene-babichenko/zpm/plugin"

	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type loadScriptArgs struct {
	FpathEntries []string
	LoadFiles    []string
}

const loadScriptTemplate = `
# Use different compdump for different zsh versions (kindly borrowed from Oh My Zsh).
ZSH_COMPDUMP=${ZDOTDIR:-${HOME}}/.zcompdump-${SHORT_HOST}-${ZSH_VERSION}
# initialize zsh completion system
autoload -U compaudit compinit
# load autocompletion scripts
{{if .FpathEntries}}
fpath=( {{range .FpathEntries}}{{.}} {{end}}$fpath )
{{end}}
compinit -u -C -d ${ZSH_COMPDUMP}
# initialize plugins
{{range .LoadFiles}}
{{.}}{{end}}
`

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
	// This process is forked and run in background even after `zpm load` is
	// finished. Such approach makes on load update checks look much faster:
	// up to 20 ms (forked) vs 1.5-2 secs (synchronous) on my setup.
	if err := exec.Command("zpm", "check").Start(); err != nil {
		return errors.Wrap(err, "while running background update check")
	}

	return nil
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load configured plugins into the current shell",
	Run: func(cmd *cobra.Command, args []string) {
		updateCheck := viper.GetBool(configKeyOnLoadCheckForUpdates)
		installMissing := viper.GetBool(configKeyOnLoadInstallMissingPlugins)

		// Use different compdump for different zsh versions (kindly borrowed from Oh My Zsh).
		fmt.Println("ZSH_COMPDUMP=\"${ZDOTDIR:-${HOME}}/.zcompdump-${SHORT_HOST}-${ZSH_VERSION}\"")
		// initialize zsh completion system
		fmt.Println("autoload -U compaudit compinit")

		ps, err := plugin.MakePluginStorage(rootDir, pluginsSpecs)
		if err != nil {
			log.Fatalf("while reading plugin configurations: %s", err)
		}

		// check if there are downloaded updates
		ps.CheckPluginUpdates(true)
		currentVersion, err := ioutil.ReadFile(filepath.Join(rootDir, ".github_version"))
		if err != nil && !os.IsNotExist(err) {
			log.Errorf("failed to read .github_version: %s", err)
		} else if err == nil {
			if string(currentVersion) != Version {
				log.Infof("zpm update available: newer version %s, current version %s", currentVersion, Version)
				log.Infof("to download the update go to %s", updateLink)
			}
		}
		if installMissing {
			ps.InstallAll()
		}

		pluginLoadData := loadScriptArgs{}

		// plugin load order must be preserved because of dependencies between them
		for _, name := range ps.LoadOrder {
			pse := ps.Plugins[name]
			fpathPlugin, execPlugin, err := pse.Plugin.Load()
			if err != nil {
				log.Errorf("while loading plugin %s: %s", pse.Name, err)
				continue
			}
			pluginLoadData.FpathEntries = append(pluginLoadData.FpathEntries, fpathPlugin...)
			pluginLoadData.LoadFiles = append(pluginLoadData.LoadFiles, execPlugin...)
		}

		tmpl, err := template.New("load").Parse(loadScriptTemplate)
		if err != nil {
			log.Fatalf("failed to parse the loader template: %s", err)
		}
		if err := tmpl.Execute(os.Stdout, pluginLoadData); err != nil {
			log.Fatalf("failed to execute the loader template: %s", err)
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
