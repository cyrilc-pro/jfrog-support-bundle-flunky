package commands

import (
	"fmt"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func createSupportBundle(ctx *components.Context, rtDetails *config.ArtifactoryDetails, conf *supportBundleCommandConfiguration) (creationResponse, error) {
	log.Debug(fmt.Sprintf("Create Support Bundle %s on %s", conf.caseNumber, rtDetails.Url))
	return creationResponse{Id: "foo-id"}, nil
}
