package commands

import (
	"encoding/json"
	"fmt"
)

type SupportBundleCreationOptions struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Parameters  *SupportBundleParameters `json:"parameters"`
}

type SupportBundleParameters struct {
	Configuration bool                               `json:"configuration"`
	Logs          *SupportBundleParametersLogs       `json:"logs"`
	System        bool                               `json:"system"`
	ThreadDump    *SupportBundleParametersThreadDump `json:"thread_dump"`
}

type SupportBundleParametersLogs struct {
	Include   bool   `json:"include"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type SupportBundleParametersThreadDump struct {
	Count    uint `json:"count"`
	Interval uint `json:"interval"`
}

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
