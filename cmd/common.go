package cmd

import "zpm/plugin"

func checkPluginUpdate(name string, pluginInstance plugin.Plugin) (*string, error) {
	update, err := pluginInstance.CheckUpdate()

	if plugin.IsNotInstalled(err) {
		logger.Info("not installed: ", name)
	} else if err != nil {
		logger.Errorf("while checking for %s: %s", name, err.Error())
	} else if update != nil {
		logger.Infof("update available for %s: %s", name, *update)
	} else {
		logger.Info("up to date: ", name)
	}

	return update, err
}
