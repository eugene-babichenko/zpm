package cmd

import (
	"github.com/eugene-babichenko/zpm/log"
	"github.com/eugene-babichenko/zpm/plugin"

	"fmt"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
)

type pluginState int32

const (
	pluginConfigLoaded pluginState = iota
	pluginInstalled
	pluginNeedInstall
	pluginNeedUpdate
	pluginCheckError
)

type pluginStorageEntry struct {
	name        string
	plugin      plugin.Plugin
	state       pluginState
	errorState  *error
	updateState *string
}

type pluginStorage struct {
	plugins []pluginStorageEntry
}

func makePluginStorage(
	root string,
	pluginSpecs []string,
) (ps *pluginStorage, err error) {
	ps = &pluginStorage{
		plugins: make([]pluginStorageEntry, 0, len(pluginSpecs)),
	}

	root = filepath.Join(root, "Plugins")
	factory := &plugin.Factory{Root: root}

	for _, pluginSpec := range pluginSpecs {
		p, isDependency, err := factory.MakePlugin(pluginSpec)
		if err != nil {
			msg := fmt.Sprintf("while loading plugin %s", pluginSpec)
			return nil, errors.Wrap(err, msg)
		}

		pse := pluginStorageEntry{
			name:        pluginSpec,
			plugin:      *p,
			state:       pluginConfigLoaded,
			errorState:  nil,
			updateState: nil,
		}

		// Oh My Zsh is required to be inserted in the beginning of the plugin
		// load sequence.
		if !isDependency {
			ps.plugins = append(ps.plugins, pse)
		}
	}

	factoryDependencies, factoryDependenciesSpecs := factory.Dependencies()
	if factoryDependencies == nil {
		return ps, nil
	}

	dependencies := make([]pluginStorageEntry, 0, len(factoryDependencies))

	for i, p := range factoryDependencies {
		pse := pluginStorageEntry{
			name:       factoryDependenciesSpecs[i],
			plugin:     p,
			state:      pluginConfigLoaded,
			errorState: nil,
		}
		dependencies = append(dependencies, pse)
	}

	ps.plugins = append(dependencies, ps.plugins...)
	return ps, nil
}

// checkPluginUpdates checks for both updates and plugins that are not installed
func (ps *pluginStorage) checkPluginUpdates() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(ps.plugins))

	for i, pse := range ps.plugins {
		go func(i int, pse pluginStorageEntry) {
			update, err := pse.plugin.CheckUpdate()

			if plugin.IsNotInstalled(err) {
				log.Info("not installed: %s", pse.name)
				ps.plugins[i].state = pluginNeedInstall
			} else if plugin.IsNotUpgradable(err) {
				log.Debug("plugin %s is not upgradable", pse.name)
				ps.plugins[i].state = pluginInstalled
			} else if plugin.IsUpToDate(err) {
				log.Debug("up to date: %s", pse.name)
				ps.plugins[i].state = pluginInstalled
			} else if err != nil {
				log.Error("while checking for %s: %s", pse.name, err)
				ps.plugins[i].state = pluginCheckError
				errorState := errors.Wrap(err, fmt.Sprintf("while checking for %s", pse.name))
				ps.plugins[i].errorState = &errorState
			} else if update != nil {
				updateLine := fmt.Sprintf("update available for %s: %s", pse.name, *update)
				log.Info(updateLine)
				ps.plugins[i].state = pluginNeedUpdate
				ps.plugins[i].updateState = &updateLine
			}

			waitGroup.Done()
		}(i, pse)
	}

	waitGroup.Wait()
}

// checkPluginInstalls checks for plugins that are not installed
func (ps *pluginStorage) checkPluginInstalls() {
	for i, pse := range ps.plugins {
		isInstalled, err := pse.plugin.IsInstalled()
		if err != nil && err != plugin.NotInstalled {
			log.Error("while checking for %s: %s", pse.name, err)
			ps.plugins[i].state = pluginCheckError
			ps.plugins[i].errorState = &err
			ps.plugins[i].updateState = nil
		} else if !isInstalled {
			log.Info("not installed: %s", pse.name)
			ps.plugins[i].state = pluginNeedInstall
			ps.plugins[i].errorState = nil
			ps.plugins[i].updateState = nil
		}
	}

	log.Info("plugins: %v", ps.plugins)
}

func (ps *pluginStorage) updateAll() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(ps.plugins))

	for i, pse := range ps.plugins {
		go func(i int, pse pluginStorageEntry) {
			defer waitGroup.Done()

			if pse.state != pluginNeedUpdate {
				return
			}

			if err := pse.plugin.InstallUpdate(); err != nil {
				log.Error("while updating %s: %s", pse.name, err)
				ps.plugins[i].state = pluginCheckError
				errorState := errors.Wrap(err, "while updating %s")
				ps.plugins[i].errorState = &errorState
				return
			}

			log.Info("installed update for %s: %s", pse.name, *pse.updateState)
			ps.plugins[i].state = pluginInstalled
			ps.plugins[i].updateState = nil
		}(i, pse)
	}

	waitGroup.Wait()
}

// installAll installs all plugins detected by checkPluginInstalls or checkPluginUpdates
func (ps *pluginStorage) installAll() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(ps.plugins))

	for i, pse := range ps.plugins {
		go func(i int, pse pluginStorageEntry) {
			defer waitGroup.Done()

			if pse.state != pluginNeedInstall {
				log.Info("this plugin is not required to be installed: %s", pse.name)
				return
			}

			if err := pse.plugin.InstallUpdate(); err != nil {
				log.Error("while installing %s: %s", pse.name, err)
				ps.plugins[i].state = pluginCheckError
				errorState := errors.Wrap(err, "while installing %s")
				ps.plugins[i].errorState = &errorState
				return
			}

			log.Info("installed %s", pse.name)
			ps.plugins[i].state = pluginInstalled
		}(i, pse)
	}

	waitGroup.Wait()
}
