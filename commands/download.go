package commands

import (
	"fmt"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

//nolint: unparam, lll
func downloadSupportBundle(ctx *components.Context, rtDetails *config.ArtifactoryDetails, conf *supportBundleCommandConfiguration, response creationResponse) (string, error) {
	tmpFile := "/tmp/foo.tmp"
	log.Debug(fmt.Sprintf("Download Support Bundle %s from %s to %s", response.ID, rtDetails.Url, tmpFile))
	return tmpFile, nil
}
