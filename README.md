# JFrog Support Bundle Flunky

## About this plugin
This plugin helps with generating and exchanging a JFrog Support Bundle with JFrog Support.

## Installation with JFrog CLI
Installing the latest version:

`$ jfrog plugin install jfrog-support-bundle-flunky`

Installing a specific version:

`$ jfrog plugin install jfrog-support-bundle-flunky@version`

Uninstalling the plugin

`$ jfrog plugin uninstall jfrog-support-bundle-flunky`

## Usage

This plugin has a unique command `support-bundle` that:
1. Creates a Support Bundle on the target Artifactory service
2. Downloads the Support Bundle locally in a temporary file
3. Uploads the Support Bundle on JFrog "dropbox" service

**Arguments**
- `case` - The JFrog Support case number (required).

**Example**
```
$ jfrog support-bundle 1234
```

**Optional flags**
- `server-id` - The ID of the target Artifactory service in JFrog CLI configuration. Example: `--server-id=my-jfrog-service`
- `download-timeout` - Timeout of the Support Bundle download. Example: `--download-timeout=10s`
- `retry-interval` - Waiting time between a failed download attempt and the next attempt. Example: `--retry-interval=3s`
- `prompt-options` - Ask what is to be included in the created Support Bundle. Example: `--prompt-options` 

### Environment variables
None.

## Additional info
None.

## Release Notes
The release notes are available [here](RELEASE.md).
