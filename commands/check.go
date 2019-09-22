package commands

import (
	"github.com/eugene-babichenko/zpm/plugin"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("checking for updates...")

		ps, err := plugin.MakePluginStorage(rootDir, pluginsSpecs)
		if err != nil {
			log.Fatalf("while reading plugin configurations: %s", err)
		}

		ps.CheckPluginUpdates(false)

		response, err := http.Get(updateAPILink)
		if err != nil {
			log.Errorf("failed to check for zpm update: %s", err)
			return
		}
		if response.StatusCode != 200 {
			log.Errorf("failed to check for zpm update: HTTP response status %d", response.StatusCode)
			return
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Errorf("failed to check for zpm update: %s", err)
			return
		}
		responseJSON := make(map[string]interface{})
		if err := json.Unmarshal(body, &responseJSON); err != nil {
			log.Errorf("failed to check for zpm update: %s", err)
			return
		}
		version, ok := responseJSON["tag_name"]
		if !ok {
			log.Error("failed to check for zpm update: no field named tag_name in the response")
			return
		}
		versionString, ok := version.(string)
		if !ok {
			log.Error("failed to check for zpm update: field tag_name is not a string")
			return
		}
		if err := ioutil.WriteFile(filepath.Join(rootDir, ".github_version"), []byte(versionString[1:]), os.ModePerm); err != nil {
			log.Errorf("failed to save zpm update info: %s", err)
			return
		}

		if versionString[1:] != Version {
			log.Infof("zpm update available: newer version %s, current version %s", versionString[1:], Version)
			log.Infof("to download the update go to %s", updateLink)
		}

		log.Info("update check finished")
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
