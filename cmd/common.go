package cmd

import (
	"github.com/eugene-babichenko/zpm/log"
	"github.com/eugene-babichenko/zpm/meta"
	"github.com/eugene-babichenko/zpm/plugin"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

var Version string

type checkResult int32

const (
	checkDone checkResult = iota
	checkFailed
	checkNeedInstall
	checkNeedUpdate
)

func installUpdate(name string, pluginInstance plugin.Plugin) {
	log.Info("installing: %s", name)
	if err := pluginInstance.InstallUpdate(); err != nil {
		log.Error("while installing %s: %s", name, err.Error())
		return
	}
	log.Info("installed: %s", name)
}

func checkPluginUpdate(
	name string,
	pluginInstance plugin.Plugin,
	installUpdates bool,
	installMissing bool,
) checkResult {
	update, err := pluginInstance.CheckUpdate()

	if plugin.IsNotInstalled(err) {
		log.Info("not installed: %s", name)
		if installMissing {
			installUpdate(name, pluginInstance)
		} else {
			return checkNeedInstall
		}
	} else if plugin.IsNotUpgradable(err) {
		log.Debug("plugin %s is not upgradable", name)
	} else if plugin.IsUpToDate(err) {
		log.Debug("up to date: %s", name)
	} else if err != nil {
		log.Error("while checking for %s: %s", name, err)
		return checkFailed
	} else if update != nil {
		log.Info("update available for %s: %s", name, *update)
		if installUpdates {
			installUpdate(name, pluginInstance)
		} else {
			return checkNeedUpdate
		}
	}

	return checkDone
}

func checkAndInstallUpdates(
	names []string,
	plugins []plugin.Plugin,
	installUpdates bool,
	installMissing bool,
) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(plugins))

	var updatesAvailable int32
	var installationsAvailable int32

	allChecksSuccessful := true

	for idx, pluginInstance := range plugins {
		go func(name string, pluginInstance plugin.Plugin) {
			switch checkPluginUpdate(name, pluginInstance, installUpdates, installMissing) {
			case checkFailed:
				allChecksSuccessful = false
			case checkNeedInstall:
				atomic.AddInt32(&installationsAvailable, 1)
			case checkNeedUpdate:
				atomic.AddInt32(&updatesAvailable, 1)
			}

			waitGroup.Done()
		}(names[idx], pluginInstance)
	}

	waitGroup.Wait()

	if updatesAvailable > 0 || installationsAvailable > 0 {
		log.Info(
			"%d updates available and %d plugins need to be installed.",
			updatesAvailable,
			installationsAvailable,
		)
		log.Info("You can run the update using `zpm update`.")
	} else {
		log.Info("Everything is up to date.")
	}

	if allChecksSuccessful {
		meta := meta.Meta{
			LastUpdateCheck:       time.Now(),
			UpdatesAvailable:      updatesAvailable,
			InstallationsRequired: installationsAvailable,
		}
		newMetaJSON, err := meta.Marshal()
		if err != nil {
			log.Fatal("failed to encode the meta file: %s", err)
		}
		if err := ioutil.WriteFile(metaPath(), []byte(newMetaJSON), os.ModePerm); err != nil {
			log.Fatal("failed to write down the meta file: %s", err)
		}
	}
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
			msg := fmt.Sprintf("while loading plugin %s", pluginSpec)
			return nil, nil, errors.Wrap(err, msg)
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
