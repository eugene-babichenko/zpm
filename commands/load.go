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

// Note the part I took from Oh My Zsh
//Copyright (c) 2009-2019 Robby Russell and contributors
//See the full list at https://github.com/robbyrussell/oh-my-zsh/contributors
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.

const loadScriptTemplate = `
### TAKEN FROM OH MY ZSH

# Figure out the SHORT hostname
if [[ "$OSTYPE" = darwin* ]]; then
	# macOS's $HOST changes with dhcp, etc. Use ComputerName if possible.
	SHORT_HOST=$(scutil --get ComputerName 2>/dev/null) || SHORT_HOST=${HOST/.*/}
else
	SHORT_HOST=${HOST/.*/}
fi
ZSH_COMPDUMP=${ZDOTDIR:-${HOME}}/.zcompdump-${SHORT_HOST}-${ZSH_VERSION}

### TAKEN FROM OH MY ZSH

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

if [ -z "$ZPM_BINARY" ]; then
	ZPM_BINARY=$(which zpm)
fi
zpm () {
	$ZPM_BINARY $@
	if [ "$1" = "update" ] || [ "$1" = "install" ]; then
		echo "zpm: Loading updates..."
		source <($ZPM_BINARY load)
	fi
}
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
			showVersionUpdateGuide(string(currentVersion))
		}
		if installMissing {
			ps.InstallAll()
		}

		if ps.HasUpdates() {
			log.Info("To install updates, run `zpm update`")
		}

		if ps.HasInstalls() {
			log.Info("To install new plugins, run `zpm install`")
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
