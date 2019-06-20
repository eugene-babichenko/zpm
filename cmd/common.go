package cmd

import (
	"encoding/json"
	"github.com/eugene-babichenko/zpm/meta"
	"github.com/eugene-babichenko/zpm/plugin"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
)

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

func MakePluginsFromSpecs(
	root string,
	pluginSpecs []string,
) (names []string, plugins []plugin.Plugin, err error) {
	for _, pluginSpec := range pluginSpecs {
		p, err := plugin.MakePlugin(root, pluginSpec)
		if err != nil {
			return nil, nil, errors.Wrap(err, "while loading plugins")
		}
		if pluginSpec != "oh-my-zsh" {
			plugins = append(plugins, *p)
			names = append(names, pluginSpec)
		}
	}

	if ohMyZsh := plugin.GetOhMyZsh(); ohMyZsh != nil {
		plugins = append([]plugin.Plugin{*ohMyZsh}, plugins...)
		names = append([]string{"oh-my-zsh"}, names...)
	}

	return names, plugins, nil
}

func updateMeta() {
	newMeta := meta.Meta{
		LastUpdateCheck: time.Now().Format(meta.LastUpdateCheckLayout),
	}
	newMetaJSON, err := json.Marshal(newMeta)
	if err != nil {
		logger.Fatal("failed to write down the meta file: ", err.Error())
	}
	if err := ioutil.WriteFile(metaFilePath, []byte(newMetaJSON), os.ModePerm); err != nil {
		logger.Fatal("failed to write down the meta file: ", err.Error())
	}
}
