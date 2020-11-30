package commands

import (
	"fmt"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"os"
)

//nolint: unparam
func uploadSupportBundle(ctx *components.Context, conf *SupportBundleCommandConfiguration, file string) error {
	log.Debug(fmt.Sprintf("Upload Support Bundle from %s to https://supportlogs.jfrog.com/", file))
	log.Debug(fmt.Sprintf("Deleting file: %s", file))
	return os.Remove(file)
}
