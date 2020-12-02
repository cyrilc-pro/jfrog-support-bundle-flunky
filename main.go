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
	app := components.App{}
	app.Name = "jfrog-support-bundle-flunky"
	app.Description = "This plugin dutifully creates a Support Bundle on an Artifactory service and obediently " +
		"uploads it to another Artifactory service."
	app.Version = "v0.1.0"
	app.Commands = getCommands()
	return app
}

func getCommands() []components.Command {
	return []components.Command{
		commands.GetSupportBundleCommand()}
}
