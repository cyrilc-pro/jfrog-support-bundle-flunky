package http

import (
	"encoding/json"
	"fmt"
)

// JSONObject is the map representation of a JSON object.
type JSONObject map[string]interface{}

// ParseJSON parses bytes into a JSONObject.
func ParseJSON(bytes []byte) (JSONObject, error) {
	parsedResponse := make(JSONObject)
	err := json.Unmarshal(bytes, &parsedResponse)
	return parsedResponse, err
}

// GetString gets the value of a given JSON property.
func (o JSONObject) GetString(p string) (string, error) {
	v, ok := o[p]
	if !ok {
		return "", fmt.Errorf("property %s not found", p)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("property %s is not a string", p)
	}
	return s, nil
}
