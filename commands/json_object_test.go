package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ParseJSON_BadJSON(t *testing.T) {
	_, err := ParseJSON([]byte("not json"))
	require.Error(t, err)
}

func Test_ParseJSON_ValidJSON(t *testing.T) {
	o, err := ParseJSON([]byte(`{"foo": "bar"}`))
	require.NoError(t, err)
	assert.Equal(t, o, JSONObject{
		"foo": "bar",
	})
}

func Test_JSONObjectGetString_NoProperty(t *testing.T) {
	o := JSONObject{
		"foo": "bar",
	}
	_, err := o.GetString("unknown")
	require.EqualError(t, err, "property unknown not found")
}

func Test_JSONObjectGetString_NotString(t *testing.T) {
	o := JSONObject{
		"foo": 1,
	}
	_, err := o.GetString("foo")
	require.EqualError(t, err, "property foo is not a string")
}

func Test_JSONObjectGetString_String(t *testing.T) {
	o := JSONObject{
		"foo": "bar",
	}
	v, err := o.GetString("foo")
	require.NoError(t, err)
	assert.Equal(t, v, "bar")
}
