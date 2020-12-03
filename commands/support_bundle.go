package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/artifactory/commands"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/actions"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	serverIDFlag        = "server-id"
	targetServerIDFlag  = "target-server-id"
	downloadTimeoutFlag = "download-timeout"
	retryIntervalFlag   = "retry-interval"
	promptOptionsFlag   = "prompt-options"
	cleanupFlag         = "cleanup"
	targetRepoFlag      = "target-repo"
)

// GetSupportBundleCommand returns the description of the "support-bundle" command.
func GetSupportBundleCommand() components.Command {
	return components.Command{
		Name:        "support-bundle",
		Description: `Creates a Support Bundle and uploads it to JFrog Support "dropbox" service`,
		Aliases:     []string{"sb"},
		Arguments:   getArguments(),
		Flags:       getFlags(),
		EnvVars:     nil,
		Action:      supportBundleCmd,
	}
}

func getArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "case",
			Description: "JFrog Support case number.",
		},
	}
}

func getFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name: serverIDFlag,
			Description: "Artifactory server ID configured using the config command. " +
				"If not provided the default configuration will be used.",
		},
		components.StringFlag{
			Name: targetServerIDFlag,
			Description: "Artifactory server ID configured using the config command to be used as the target for " +
				"uploading the generated Support Bundle. If not provided JFrog support logs will be used.",
		},
		components.StringFlag{
			Name:         downloadTimeoutFlag,
			Description:  "The timeout for download.",
			DefaultValue: "10m",
		},
		components.StringFlag{
			Name:         retryIntervalFlag,
			Description:  "The duration to wait between retries.",
			DefaultValue: "5s",
		},
		components.BoolFlag{
			Name:        promptOptionsFlag,
			Description: "Ask for support bundle options or use Artifactory default options.",
		},
		components.BoolFlag{
			Name:         cleanupFlag,
			Description:  "Delete the support bundle local temp file after upload.",
			DefaultValue: true,
		},
		components.StringFlag{
			Name:         targetRepoFlag,
			Description:  "The target repository key where the support bundle will be uploaded to.",
			DefaultValue: "logs",
		},
	}
}

type artifactoryServiceHelper struct{}

func (cw *artifactoryServiceHelper) GetConfig(serverID string, excludeRefreshableTokens bool) (*config.ArtifactoryDetails, error) {
	return commands.GetConfig(serverID, excludeRefreshableTokens)
}

func (cw *artifactoryServiceHelper) CreateInitialRefreshableTokensIfNeeded(artifactoryDetails *config.ArtifactoryDetails) error {
	return config.CreateInitialRefreshableTokensIfNeeded(artifactoryDetails)
}

func supportBundleCmd(componentContext *components.Context) error {
	ctx := context.Background()
	caseNumber, err := parseArguments(componentContext)
	if err != nil {
		return err
	}

	artifactoryConfigHelper := &artifactoryServiceHelper{}
	rtDetails, err := getRtDetails(componentContext, artifactoryConfigHelper)
	if err != nil {
		return err
	}
	log.Debug(fmt.Sprintf("Using: %s...", rtDetails.Url))
	log.Output(fmt.Sprintf("Case number is %s", caseNumber))

	targetRtDetails, err := getTargetDetails(componentContext, artifactoryConfigHelper)
	if err != nil {
		return err
	}

	client := &http.Client{RtDetails: rtDetails}
	targetClient := &http.Client{RtDetails: targetRtDetails}

	// 1. Create Support Bundle
	supportBundle, err := actions.CreateSupportBundle(client, caseNumber, getPromptOptions(componentContext))
	if err != nil {
		return err
	}

	// 2. Download Support Bundle
	supportBundleArchivePath, err := actions.DownloadSupportBundle(ctx, client, getTimeout(componentContext),
		getRetryInterval(componentContext), supportBundle)
	if err != nil {
		return err
	}
	if shouldCleanup(componentContext) {
		defer deleteSupportBundleArchive(supportBundleArchivePath)
	}

	// 3. Upload Support Bundle
	return actions.UploadSupportBundle(targetClient, caseNumber, supportBundleArchivePath, getTargetRepo(componentContext), time.Now)
}

func deleteSupportBundleArchive(supportBundleArchivePath string) {
	log.Debug(fmt.Sprintf("Deleting generated support bundle: %s", supportBundleArchivePath))
	err := os.Remove(supportBundleArchivePath)
	if err != nil {
		log.Warn(fmt.Sprintf("Error occurred while deleting the generated support bundle archive: %+v", err))
	}
}

func parseArguments(ctx *components.Context) (actions.CaseNumber, error) {
	if len(ctx.Arguments) != 1 {
		return "", errors.New("Wrong number of arguments. Expected: 1, " + "Received: " + strconv.Itoa(len(ctx.Arguments)))
	}
	return actions.CaseNumber(strings.TrimSpace(ctx.Arguments[0])), nil
}
