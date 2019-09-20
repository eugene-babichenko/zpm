package cmd

import (
	"github.com/eugene-babichenko/zpm/plugin"

	"fmt"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	errorState  error
	updateState *string
}

type pluginStorage struct {
	plugins map[string]*pluginStorageEntry
	// the order in which plugins are loaded is important, so we must preserve it
	loadOrder []string
}

func makePluginStorage(
	root string,
	pluginSpecs []string,
) (ps *pluginStorage, err error) {
	ps = &pluginStorage{
		plugins: make(map[string]*pluginStorageEntry),
	}

	root = filepath.Join(root, "Plugins")
	factory := &plugin.Factory{Root: root}

	for _, pluginSpec := range pluginSpecs {
		p, isDependency, err := factory.MakePlugin(pluginSpec)
		if err != nil {
			msg := fmt.Sprintf("while loading plugin %s", pluginSpec)
			return nil, errors.Wrap(err, msg)
		}

		pse := &pluginStorageEntry{
			name:        pluginSpec,
			plugin:      *p,
			state:       pluginConfigLoaded,
			errorState:  nil,
			updateState: nil,
		}

		// Oh My Zsh is required to be inserted in the beginning of the plugin
		// load sequence.
		if !isDependency {
			ps.plugins[pse.name] = pse
			ps.loadOrder = append(ps.loadOrder, pse.name)
		}
	}

	factoryDependencies, factoryDependenciesSpecs := factory.Dependencies()
	if factoryDependencies == nil {
		return ps, nil
	}

	dependencies := make([]*pluginStorageEntry, 0, len(factoryDependencies))

	for i, p := range factoryDependencies {
		pse := &pluginStorageEntry{
			name:       factoryDependenciesSpecs[i],
			plugin:     p,
			state:      pluginConfigLoaded,
			errorState: nil,
		}
		dependencies = append(dependencies, pse)
		ps.plugins[pse.name] = pse
	}

	dependenciesNames := make([]string, len(dependencies))
	for i, pse := range dependencies {
		dependenciesNames[i] = pse.name
	}
	ps.loadOrder = append(dependenciesNames, ps.loadOrder...)

	return ps, nil
}

func (pse *pluginStorageEntry) updateInternal() bool {
	if err := pse.plugin.InstallUpdate(); err != nil {
		log.Errorf("while installing %s: %s", pse.name, err)
		pse.state = pluginCheckError
		errorState := errors.Wrap(err, "while installing %s")
		pse.errorState = errorState
		return false
	}
	return true
}

func (pse *pluginStorageEntry) update() {
	if pse.state != pluginNeedUpdate {
		return
	}

	if pse.updateInternal() {
		log.Infof("installed update for %s: %s", pse.name, *pse.updateState)
		pse.state = pluginInstalled
		pse.updateState = nil
	}
}

func (pse *pluginStorageEntry) install() {
	if pse.state != pluginNeedInstall {
		log.Debugf("this plugin is not required to be installed: %s", pse.name)
		return
	}

	if pse.updateInternal() {
		log.Infof("installed %s", pse.name)
		pse.state = pluginInstalled
	}
}

func (pse *pluginStorageEntry) checkPluginInstall() {
	isInstalled, err := pse.plugin.IsInstalled()
	if err != nil && err != plugin.NotInstalled && err != plugin.NotInstallable {
		log.Errorf("while checking for %s: %s", pse.name, err)
		pse.state = pluginCheckError
		pse.errorState = err
		pse.updateState = nil
	} else if !isInstalled && err != plugin.NotInstallable {
		log.Infof("not installed: %s", pse.name)
		pse.state = pluginNeedInstall
		pse.errorState = nil
		pse.updateState = nil
	}
}

func (pse *pluginStorageEntry) checkPluginUpdate(offline bool) {
	update, err := pse.plugin.CheckUpdate(offline)

	if plugin.IsNotInstalled(err) {
		log.Infof("not installed: %s", pse.name)
		pse.state = pluginNeedInstall
	} else if err != plugin.NotInstallable {
		log.Debugf("plugin %s is not installable", pse.name)
	} else if plugin.IsNotUpgradable(err) {
		log.Debugf("plugin %s is not upgradable", pse.name)
		pse.state = pluginInstalled
	} else if plugin.IsUpToDate(err) {
		log.Debugf("up to date: %s", pse.name)
		pse.state = pluginInstalled
	} else if err != nil {
		log.Errorf("while checking for %s: %s", pse.name, err)
		pse.state = pluginCheckError
		errorState := errors.Wrap(err, fmt.Sprintf("while checking for %s", pse.name))
		pse.errorState = errorState
	} else if update != nil {
		updateLine := fmt.Sprintf("update available for %s: %s", pse.name, *update)
		log.Infof(updateLine)
		pse.state = pluginNeedUpdate
		pse.updateState = &updateLine
	}
}

// checkPluginUpdates checks for both updates and plugins that are not installed
func (ps *pluginStorage) checkPluginUpdates(offline bool) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(ps.plugins))
	for i := range ps.plugins {
		go func(i string) {
			ps.plugins[i].checkPluginUpdate(offline)
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()
}

// checkPluginInstalls checks for plugins that are not installed
func (ps *pluginStorage) checkPluginInstalls() {
	for i := range ps.plugins {
		ps.plugins[i].checkPluginInstall()
	}
}

func (ps *pluginStorage) updateAll() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(ps.plugins))
	for i := range ps.plugins {
		go func(i string) {
			ps.plugins[i].update()
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()
}

// installAll installs all plugins detected by checkPluginInstalls or checkPluginUpdates
func (ps *pluginStorage) installAll() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(ps.plugins))
	for i := range ps.plugins {
		go func(i string) {
			ps.plugins[i].install()
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()
}
