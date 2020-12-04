package main

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getApp(t *testing.T) {
	expectedApp := components.App{
		Name: "sb-flunky",
		Description: "This plugin dutifully creates a Support Bundle on an Artifactory service and obediently " +
			"uploads it to another Artifactory service.",
		Version: "v0.1.0",
	}
	assert.Empty(t, cmp.Diff(expectedApp, getApp(), cmpopts.IgnoreFields(components.App{}, "Commands")))
}
