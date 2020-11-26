package commands

import (
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"strconv"
	"strings"
)

func GetSupportBundleCommand() components.Command {
	return components.Command{
		Name:        "support-bundle",
		Description: "TBD",
		Aliases:     []string{"sb"},
		Arguments:   getArguments(),
		Flags:       getFlags(),
		EnvVars:     nil,
		Action: func(c *components.Context) error {
			return supportBundleCmd(c)
		},
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
			Name:        "server-id",
			Description: "Artifactory server ID configured using the config command.",
		},
	}
}

type supportBundleCommandConfiguration struct {
	caseNumber string
}

func supportBundleCmd(ctx *components.Context) error {
	conf, err := parseArguments(ctx)
	if err != nil {
		return err
	}

	rtDetails, err := getRtDetails(ctx)
	if err != nil {
		return err
	}
	log.Debug(fmt.Sprintf("Using: %s...", rtDetails.Url))
	log.Output(fmt.Sprintf("Case number is %s", conf.caseNumber))

	client := &HttpClient{rtDetails: rtDetails}

	// 1. Create Support Bundle
	response, err := createSupportBundle(client, conf, Now)
	if err != nil {
		return err
	}
	// 2. Download Support Bundle
	tmpFile, err := downloadSupportBundle(ctx, rtDetails, conf, response)
	if err != nil {
		return err
	}
	// 3. Upload Support Bundle
	return uploadSupportBundle(ctx, conf, tmpFile)
}

func parseArguments(ctx *components.Context) (*supportBundleCommandConfiguration, error) {
	if len(ctx.Arguments) != 1 {
		return nil, errors.New("Wrong number of arguments. Expected: 1, " + "Received: " + strconv.Itoa(len(ctx.Arguments)))
	}
	var conf = new(supportBundleCommandConfiguration)
	conf.caseNumber = strings.TrimSpace(ctx.Arguments[0])
	return conf, nil
}
