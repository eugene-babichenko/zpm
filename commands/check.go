package commands

import (
	"github.com/eugene-babichenko/zpm/plugin"

	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func showVersionUpdateGuide(currentVersion string) {
	if cmp := strings.Compare(currentVersion, Version); cmp > 0 {
		log.Infof("zpm update available: newer version %s, current version %s", currentVersion, Version)
		log.Infof("to download the update go to %s", updateLink)
	}
	// "-" can appear in generated version numbers or in versions like v0.2.1-beta.1
	if strings.Contains(Version, "-") {
		log.Warnf("Be careful! You are running a possibly unstable version %s", Version)
	}
}

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

		if ps.HasUpdates() {
			log.Info("To install updates, run `zpm update`")
		}

		if ps.HasInstalls() {
			log.Info("To install new plugins, run `zpm install`")
		}

		githubClient := github.NewClient(nil)
		release, _, err := githubClient.Repositories.GetLatestRelease(context.Background(), "eugene-babichenko", "zpm")
		if err != nil {
			log.Fatalf("failed to check for zpm update: %s", err)
		}

		releaseTag := *(release.TagName)
		releaseTag = releaseTag[1:]

		showVersionUpdateGuide(releaseTag)

		if err := ioutil.WriteFile(filepath.Join(rootDir, ".github_version"), []byte(releaseTag), os.ModePerm); err != nil {
			log.Fatalf("failed to write .github_version: %s", err)
		}

		log.Info("update check finished")
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
