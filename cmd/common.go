package cmd

import (
	"github.com/eugene-babichenko/zpm/log"
	"github.com/eugene-babichenko/zpm/meta"
	"github.com/eugene-babichenko/zpm/plugin"

	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

var Version string

func checkPluginUpdate(name string, pluginInstance plugin.Plugin) (*string, error) {
	update, err := pluginInstance.CheckUpdate()

	if plugin.IsNotInstalled(err) {
		log.Info("not installed: %s", name)
	} else if plugin.IsNotUpgradable(err) {
		log.Debug("plugin %s is not upgradable", name)
	} else if plugin.IsUpToDate(err) {
		log.Debug("up to date: %s", name)
	} else if err != nil {
		log.Error("while checking for %s: %s", name, err)
	} else if update != nil {
		log.Info("update available for %s: %s", name, *update)
	}

	return update, err
}

// Make plugins from the specification values. `names` are required for display
// in logs.
func MakePluginsFromSpecs(
	root string,
	pluginSpecs []string,
) (names []string, plugins []plugin.Plugin, err error) {
	root = filepath.Join(root, "Plugins")
	factory := &plugin.Factory{Root: root}

	for _, pluginSpec := range pluginSpecs {
		p, isDependency, err := factory.MakePlugin(pluginSpec)
		if err != nil {
			return nil, nil, errors.Wrap(err, "while loading plugins")
		}
		// Oh My Zsh is required to be inserted in the beginning of the plugin
		// load sequence.
		if !isDependency {
			plugins = append(plugins, *p)
			names = append(names, pluginSpec)
		}
	}

	// Oh My Zsh is required to be inserted in the beginning of the plugin load
	// sequence.
	if dependencies, dependenciesNames := factory.Dependencies(); dependencies != nil {
		plugins = append(dependencies, plugins...)
		names = append(dependenciesNames, names...)
	}

	return names, plugins, nil
}

func metaPath() string {
	return filepath.Join(appConfig.Root, "meta.json")
}

// Update the last time of check for updates.
func updateLastUpdateCheckTime() {
	newMeta := meta.Meta{
		LastUpdateCheck: time.Now().Format(meta.LastUpdateCheckLayout),
	}
	newMetaJSON, err := json.Marshal(newMeta)
	if err != nil {
		log.Fatal("failed to write down the meta file: %s", err)
	}
	if err := ioutil.WriteFile(metaPath(), []byte(newMetaJSON), os.ModePerm); err != nil {
		log.Fatal("failed to write down the meta file: %s", err)
	}
}

func readLastUpdateCheckTime() time.Time {
	var lastUpdateCheckTime time.Time
	metaFile, err := ioutil.ReadFile(metaPath())
	if err != nil {
		return lastUpdateCheckTime
	}
	var metaData meta.Meta
	if err := json.Unmarshal(metaFile, &metaData); err != nil {
		return lastUpdateCheckTime
	}
	lastUpdateCheckTime, _ = time.Parse(meta.LastUpdateCheckLayout, metaData.LastUpdateCheck)
	return lastUpdateCheckTime
}
