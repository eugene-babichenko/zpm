package cmd

import (
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
		logger.Info("not installed: ", name)
	} else if err != nil {
		logger.Errorf("while checking for %s: %s", name, err.Error())
	} else if update != nil {
		logger.Infof("update available for %s: %s", name, *update)
	} else {
		logger.Debug("up to date: ", name)
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
	for _, pluginSpec := range pluginSpecs {
		p, err := plugin.MakePlugin(root, pluginSpec)
		if err != nil {
			return nil, nil, errors.Wrap(err, "while loading plugins")
		}
		// Oh My Zsh is required to be inserted in the beginning of the plugin
		// load sequence.
		if pluginSpec != "oh-my-zsh" {
			plugins = append(plugins, *p)
			names = append(names, pluginSpec)
		}
	}

	// Oh My Zsh is required to be inserted in the beginning of the plugin load
	// sequence.
	if ohMyZsh := plugin.GetOhMyZsh(); ohMyZsh != nil {
		plugins = append([]plugin.Plugin{*ohMyZsh}, plugins...)
		names = append([]string{"oh-my-zsh"}, names...)
	}

	return names, plugins, nil
}

func cachePath() string {
	return filepath.Join(appConfig.Root, "cache-"+Version+".zsh")
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
		logger.Fatal("failed to write down the meta file: ", err.Error())
	}
	if err := ioutil.WriteFile(metaPath(), []byte(newMetaJSON), os.ModePerm); err != nil {
		logger.Fatal("failed to write down the meta file: ", err.Error())
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
