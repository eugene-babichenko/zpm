package cmd

import (
	"github.com/eugene-babichenko/zpm/plugin"

	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func setLastUpdateTime(t time.Time) error {
	s := t.Format(time.RFC3339)
	filename := filepath.Join(rootDir, ".lastupdate")
	err := ioutil.WriteFile(filename, []byte(s), os.ModePerm)
	return err
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Install updates",
	Run: func(cmd *cobra.Command, args []string) {
		pluginToUpdate, _ := cmd.Flags().GetString("plugin")

		log.Info("updating installed plugins...")
		log.Info("not installing new plugins! Run `zpm install` to do it.")

		ps, err := plugin.MakePluginStorage(rootDir, pluginsSpecs)
		if err != nil {
			log.Fatalf("while reading plugin configurations: %s", err)
		}

		if pluginToUpdate == "" {
			ps.CheckPluginUpdates(false)
			ps.UpdateAll()
			if err := setLastUpdateTime(time.Now()); err != nil {
				log.Errorf("failed to write last update time: %s", err)
				log.Error("note that this will result in extra update checks on zsh load")
			}
			log.Info("update finished")
			return
		}

		pse, ok := ps.Plugins[pluginToUpdate]
		if !ok {
			log.Fatalf("plugin %s not listed in the configuration file", pluginToUpdate)
		}
		pse.CheckPluginUpdate(false)
		pse.Update()

		log.Info("update finished")
	},
}

func init() {
	updateCmd.Flags().String(
		"plugin",
		"",
		"Update only the specified plugin",
	)

	RootCmd.AddCommand(updateCmd)
}
