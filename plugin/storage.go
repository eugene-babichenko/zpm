package plugin

import (
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
	Name        string
	Plugin      Plugin
	state       pluginState
	errorState  error
	updateState *string
}

type pluginStorage struct {
	Plugins map[string]*pluginStorageEntry
	// the order in which plugins are loaded is important, so we must preserve it
	LoadOrder []string
}

func MakePluginStorage(
	root string,
	pluginSpecs []string,
) (ps *pluginStorage, err error) {
	ps = &pluginStorage{
		Plugins: make(map[string]*pluginStorageEntry),
	}

	root = filepath.Join(root, "Plugins")
	factory := &Factory{Root: root}

	for _, pluginSpec := range pluginSpecs {
		p, isDependency, err := factory.MakePlugin(pluginSpec)
		if err != nil {
			msg := fmt.Sprintf("while loading plugin %s", pluginSpec)
			return nil, errors.Wrap(err, msg)
		}

		pse := &pluginStorageEntry{
			Name:        pluginSpec,
			Plugin:      *p,
			state:       pluginConfigLoaded,
			errorState:  nil,
			updateState: nil,
		}

		// Oh My Zsh is required to be inserted in the beginning of the plugin
		// load sequence.
		if !isDependency {
			ps.Plugins[pse.Name] = pse
			ps.LoadOrder = append(ps.LoadOrder, pse.Name)
		}
	}

	factoryDependencies, factoryDependenciesSpecs := factory.Dependencies()
	if factoryDependencies == nil {
		return ps, nil
	}

	dependencies := make([]*pluginStorageEntry, 0, len(factoryDependencies))

	for i, p := range factoryDependencies {
		pse := &pluginStorageEntry{
			Name:       factoryDependenciesSpecs[i],
			Plugin:     p,
			state:      pluginConfigLoaded,
			errorState: nil,
		}
		dependencies = append(dependencies, pse)
		ps.Plugins[pse.Name] = pse
	}

	dependenciesNames := make([]string, len(dependencies))
	for i, pse := range dependencies {
		dependenciesNames[i] = pse.Name
	}
	ps.LoadOrder = append(dependenciesNames, ps.LoadOrder...)

	return ps, nil
}

func (pse *pluginStorageEntry) updateInternal() bool {
	if err := pse.Plugin.InstallUpdate(); err != nil {
		log.Errorf("while installing %s: %s", pse.Name, err)
		pse.state = pluginCheckError
		errorState := errors.Wrap(err, "while installing %s")
		pse.errorState = errorState
		return false
	}
	return true
}

func (pse *pluginStorageEntry) Update() {
	if pse.state != pluginNeedUpdate {
		return
	}

	if pse.updateInternal() {
		log.Infof("installed update for %s", pse.Name)
		pse.state = pluginInstalled
		pse.updateState = nil
	}
}

func (pse *pluginStorageEntry) install() {
	if pse.state != pluginNeedInstall {
		log.Debugf("this plugin is not required to be installed: %s", pse.Name)
		return
	}

	if pse.updateInternal() {
		log.Infof("installed %s", pse.Name)
		pse.state = pluginInstalled
	}
}

func (pse *pluginStorageEntry) checkPluginInstall() {
	isInstalled, err := pse.Plugin.IsInstalled()
	if err != nil && err != NotInstalled && err != NotInstallable {
		log.Errorf("while checking for %s: %s", pse.Name, err)
		pse.state = pluginCheckError
		pse.errorState = err
		pse.updateState = nil
	} else if !isInstalled && err != NotInstallable {
		log.Infof("not installed: %s", pse.Name)
		pse.state = pluginNeedInstall
		pse.errorState = nil
		pse.updateState = nil
	}
}

func (pse *pluginStorageEntry) CheckPluginUpdate(offline bool) {
	update, err := pse.Plugin.CheckUpdate(offline)

	if IsNotInstalled(err) {
		log.Infof("not installed: %s", pse.Name)
		pse.state = pluginNeedInstall
	} else if err == NotInstallable {
		log.Debugf("plugin %s is not installable", pse.Name)
	} else if IsNotUpgradable(err) {
		log.Debugf("plugin %s is not upgradable", pse.Name)
		pse.state = pluginInstalled
	} else if IsUpToDate(err) {
		log.Debugf("up to date: %s", pse.Name)
		pse.state = pluginInstalled
	} else if err != nil {
		log.Errorf("while checking for %s: %s", pse.Name, err)
		pse.state = pluginCheckError
		errorState := errors.Wrap(err, fmt.Sprintf("while checking for %s", pse.Name))
		pse.errorState = errorState
	} else if update != nil {
		updateLine := fmt.Sprintf("update available for %s: %s", pse.Name, *update)
		log.Infof(updateLine)
		pse.state = pluginNeedUpdate
		pse.updateState = &updateLine
	}
}

// checkPluginUpdates checks for both updates and plugins that are not installed
func (ps *pluginStorage) CheckPluginUpdates(offline bool) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(ps.Plugins))
	for i := range ps.Plugins {
		go func(i string) {
			ps.Plugins[i].CheckPluginUpdate(offline)
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()
}

// checkPluginInstalls checks for plugins that are not installed
func (ps *pluginStorage) CheckPluginInstalls() {
	for i := range ps.Plugins {
		ps.Plugins[i].checkPluginInstall()
	}
}

func (ps *pluginStorage) UpdateAll() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(ps.Plugins))
	for i := range ps.Plugins {
		go func(i string) {
			ps.Plugins[i].Update()
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()
}

// installAll installs all plugins detected by checkPluginInstalls or checkPluginUpdates
func (ps *pluginStorage) InstallAll() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(ps.Plugins))
	for i := range ps.Plugins {
		go func(i string) {
			ps.Plugins[i].install()
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()
}
