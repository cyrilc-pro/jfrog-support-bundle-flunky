package commands

import (
	"encoding/json"
	"fmt"
)

// SupportBundleCreationOptions defines options for the creation of a Support Bundle.
type SupportBundleCreationOptions struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Parameters  *SupportBundleParameters `json:"parameters"`
}

// SupportBundleParameters defines the content of a Support Bundle.
type SupportBundleParameters struct {
	Configuration bool                               `json:"configuration"`
	Logs          *SupportBundleParametersLogs       `json:"logs,omitempty"`
	System        bool                               `json:"system"`
	ThreadDump    *SupportBundleParametersThreadDump `json:"thread_dump"`
}

// SupportBundleParametersLogs defines which logs are included in a Support Bundle.
type SupportBundleParametersLogs struct {
	Include   bool   `json:"include"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// SupportBundleParametersThreadDump defines which thread dumps are included in a Support Bundle.
type SupportBundleParametersThreadDump struct {
	Count    uint `json:"count"`
	Interval uint `json:"interval"`
}

// MarshalJSON serializes a SupportBundleCreationOptions to JSON.
func (p SupportBundleCreationOptions) MarshalJSON() ([]byte, error) {
	params := "{}"
	if p.Parameters != nil {
		paramsAsBytes, err := json.Marshal(p.Parameters)
		if err != nil {
			return nil, err
		}
		params = string(paramsAsBytes)
	}
	asJSON := fmt.Sprintf(`{"name":"%s","description":"%s","parameters":%s}`, p.Name, p.Description, params)
	return []byte(asJSON), nil
}
