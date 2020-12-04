package commands

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/artifactory/commands"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/actions"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/http"
	"os"
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
		Name:        "support-case",
		Description: `Creates a Support Bundle and uploads it to JFrog Support "dropbox" service`,
		Aliases:     []string{"c", "case"},
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

func supportBundleCmd(componentContext *components.Context) error {
	_, err := SupportBundleCmd(context.Background(), &cliAdapter{ctx: componentContext})
	return err
}

type cliAdapter struct {
	ctx *components.Context
}

func (p *cliAdapter) GetStringFlagValue(flagName string) string {
	return p.ctx.GetStringFlagValue(flagName)
}
func (p *cliAdapter) GetBoolFlagValue(flagName string) bool {
	return p.ctx.GetBoolFlagValue(flagName)
}
func (p *cliAdapter) GetArguments() []string {
	return p.ctx.Arguments
}
func (p *cliAdapter) GetRtDetails() (*config.ArtifactoryDetails, error) {
	return getRtDetails(p, p)
}
func (p *cliAdapter) GetTargetDetails() (*config.ArtifactoryDetails, error) {
	return getTargetDetails(p, p)
}
func (p *cliAdapter) GetConfig(serverID string, excludeRefreshableTokens bool) (*config.ArtifactoryDetails, error) {
	return commands.GetConfig(serverID, excludeRefreshableTokens)
}
func (p *cliAdapter) CreateInitialRefreshableTokensIfNeeded(artifactoryDetails *config.ArtifactoryDetails) error {
	return config.CreateInitialRefreshableTokensIfNeeded(artifactoryDetails)
}

// CliFacade is a facade for JFrog CLI APIs. Introduced to facilitate testing
type CliFacade interface {
	flagValueProvider
	argumentsProvider
	artifactoryDetailsProvider
}

// SupportBundleCmdResult gives details on what the command has done
type SupportBundleCmdResult struct {
	BundleID      actions.BundleID
	LocalFilePath string
	UploadPath    string
}

// SupportBundleCmd is the core of the command
func SupportBundleCmd(ctx context.Context, cli CliFacade) (*SupportBundleCmdResult, error) {
	caseNumber, err := parseArguments(cli)
	if err != nil {
		return nil, err
	}
	log.Output(fmt.Sprintf("Case number is %s", caseNumber))

	client, err := getRtClient(cli.GetRtDetails)
	if err != nil {
		return nil, err
	}
	log.Debug(fmt.Sprintf("Selected Artifactory: %s", client.GetURL()))

	targetClient, err := getRtClient(cli.GetTargetDetails)
	if err != nil {
		return nil, err
	}
	log.Debug(fmt.Sprintf("Selected \"dropbox\" Artifactory: %s", targetClient.GetURL()))

	result := &SupportBundleCmdResult{}
	// 1. Create Support Bundle
	result.BundleID, err = actions.CreateSupportBundle(client, caseNumber, getPromptOptions(cli))
	if err != nil {
		return result, err
	}

	// 2. Download Support Bundle
	result.LocalFilePath, err = actions.DownloadSupportBundle(ctx, client, getTimeout(cli),
		getRetryInterval(cli), result.BundleID)
	if err != nil {
		return result, err
	}
	if shouldCleanup(cli) {
		defer deleteSupportBundleArchive(result.LocalFilePath)
	}

	// 3. Upload Support Bundle
	result.UploadPath, err = actions.UploadSupportBundle(targetClient, caseNumber, result.LocalFilePath, getTargetRepo(cli), time.Now)
	return result, err
}

func getRtClient(rtDetailsProvider func() (*config.ArtifactoryDetails, error)) (*http.Client, error) {
	rtDetails, err := rtDetailsProvider()
	if err != nil {
		return nil, err
	}
	return &http.Client{RtDetails: rtDetails}, nil
}

func deleteSupportBundleArchive(supportBundleArchivePath string) {
	log.Debug(fmt.Sprintf("Deleting generated support bundle: %s", supportBundleArchivePath))
	err := os.Remove(supportBundleArchivePath)
	if err != nil {
		log.Warn(fmt.Sprintf("Error occurred while deleting the generated support bundle archive: %+v", err))
	}
}

func parseArguments(ctx argumentsProvider) (actions.CaseNumber, error) {
	arguments := ctx.GetArguments()
	if len(arguments) != 1 {
		return "", fmt.Errorf("wrong number of arguments. Expected: 1, Received: %d", len(arguments))
	}
	return actions.CaseNumber(strings.TrimSpace(arguments[0])), nil
}
