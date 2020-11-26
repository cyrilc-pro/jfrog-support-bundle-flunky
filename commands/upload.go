package commands

import (
	"fmt"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

//nolint: unparam
func uploadSupportBundle(ctx *components.Context, conf *supportBundleCommandConfiguration, file string) error {
	log.Debug(fmt.Sprintf("Upload Support Bundle from %s to https://supportlogs.jfrog.com/", file))
	return nil
}
