package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"strconv"
	"strings"
)

const (
	serverID        = "server-id"
	downloadTimeout = "download-timeout"
	retryInterval   = "retry-interval"
	promptOptions   = "prompt-options"
)

func GetSupportBundleCommand() components.Command {
	return components.Command{
		Name:        "support-bundle",
		Description: "TBD",
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
			Name:        serverID,
			Description: "Artifactory server ID configured using the config command.",
		},
		components.StringFlag{
			Name:        downloadTimeout,
			Description: "The timeout for download.",
		},
		components.StringFlag{
			Name:        retryInterval,
			Description: "The duration to wait between retries.",
		},
		components.BoolFlag{
			Name:        promptOptions,
			Description: "Ask for support bundle options.",
		},
	}
}

type SupportBundleCommandConfiguration struct {
	CaseNumber string
}

func supportBundleCmd(componentContext *components.Context) error {
	ctx := context.Background()
	conf, err := parseArguments(componentContext)
	if err != nil {
		return err
	}

	rtDetails, err := getRtDetails(componentContext)
	if err != nil {
		return err
	}
	log.Debug(fmt.Sprintf("Using: %s...", rtDetails.Url))
	log.Output(fmt.Sprintf("Case number is %s", conf.CaseNumber))

	client := &HTTPClient{RtDetails: rtDetails}

	// 1. Create Support Bundle
	supportBundle, err := CreateSupportBundle(client, conf, getPromptOptions(componentContext))
	if err != nil {
		return err
	}

	// 2. Download Support Bundle
	tmpFile, err := DownloadSupportBundle(ctx, client, getTimeout(componentContext), getRetryInterval(componentContext),
		supportBundle)
	if err != nil {
		return err
	}

	// 3. Upload Support Bundle
	return uploadSupportBundle(componentContext, conf, tmpFile)
}

func parseArguments(ctx *components.Context) (*SupportBundleCommandConfiguration, error) {
	if len(ctx.Arguments) != 1 {
		return nil, errors.New("Wrong number of arguments. Expected: 1, " + "Received: " + strconv.Itoa(len(ctx.Arguments)))
	}
	var conf = new(SupportBundleCommandConfiguration)
	conf.CaseNumber = strings.TrimSpace(ctx.Arguments[0])
	return conf, nil
}
