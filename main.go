package main

import (
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands"
)

func main() {
	plugins.PluginMain(getApp())
}

func getApp() components.App {
	return components.App{
		Name: "sb-flunky",
		Description: "This plugin dutifully creates a Support Bundle on an Artifactory service and obediently " +
			"uploads it to another Artifactory service.",
		Version:  "v0.1.0",
		Commands: getCommands(),
	}
}

func getCommands() []components.Command {
	return []components.Command{
		commands.GetSupportBundleCommand()}
}
