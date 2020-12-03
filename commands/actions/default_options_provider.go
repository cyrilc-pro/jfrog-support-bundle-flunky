package actions

import (
	"fmt"
	flunkyhttp "github.com/jfrog/jfrog-support-bundle-flunky/commands/http"
	"time"
)

// DefaultOptionsProvider provides default options for the creation of a Support Bundle.
type DefaultOptionsProvider struct {
	getDate Clock
}

// NewDefaultOptionsProvider creates a new DefaultOptionsProvider
func NewDefaultOptionsProvider() *DefaultOptionsProvider {
	return &DefaultOptionsProvider{getDate: time.Now}
}

// GetOptions gets the default options.
func (p *DefaultOptionsProvider) GetOptions(caseNumber CaseNumber) (flunkyhttp.SupportBundleCreationOptions, error) {
	return flunkyhttp.SupportBundleCreationOptions{
		Name:        fmt.Sprintf("JFrog Support Case number %s", caseNumber),
		Description: fmt.Sprintf("Generated on %s", formattedString(p.getDate())),
		Parameters:  nil,
	}, nil
}
