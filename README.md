![GitHub Workflow Status](https://img.shields.io/github/workflow/status/cyrilc-pro/jfrog-support-bundle-flunky/Go?style=plastic)
[![Codecov](https://img.shields.io/codecov/c/github/cyrilc-pro/jfrog-support-bundle-flunky?style=plastic&label=codecov)](https://codecov.io/gh/cyrilc-pro/jfrog-support-bundle-flunky)
[![Go Report](https://goreportcard.com/badge/github.com/cyrilc-pro/jfrog-support-bundle-flunky?style=plastic)](https://goreportcard.com/report/github.com/cyrilc-pro/jfrog-support-bundle-flunky)
[![Codacy grade](https://img.shields.io/codacy/grade/b286b95be72c4aa19de86f8c4a985f34?label=codacy&style=plastic)](https://www.codacy.com/gh/cyrilc-pro/jfrog-support-bundle-flunky/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=cyrilc-pro/jfrog-support-bundle-flunky&amp;utm_campaign=Badge_Grade)
[![Libraries.io dependency status for GitHub repo](https://img.shields.io/librariesio/github/cyrilc-pro/jfrog-support-bundle-flunky?label=libraries.io&style=plastic)](https://libraries.io/github/cyrilc-pro/jfrog-support-bundle-flunky)

# JFrog Support Bundle Flunky

## About this plugin

This plugin helps with generating and exchanging a JFrog Support Bundle with JFrog Support.

## Installation with JFrog CLI

Installing the latest version:

``` bash
jfrog plugin install sb-flunky
```

Installing a specific version:

``` bash
jfrog plugin install sb-flunky@version`
```

Uninstalling the plugin

``` bash
jfrog plugin uninstall sb-flunky`
```

## Usage

This plugin has a unique command `support-case` that:

1.  Creates a Support Bundle on the target Artifactory service

2.  Downloads the Support Bundle locally to a temporary file

3.  Uploads the Support Bundle on JFrog "dropbox" service or to any Artifactory service registered in JFrog CLI 
    configuration

### Arguments

-   `support-case` - The JFrog Support case number (required).

### Aliases

-   `case`
-   `c`

### Examples

``` bash
jfrog sb-flunky support-case 1234
```

or

``` bash
jfrog sb-flunky case 1234
```

or

``` bash
jfrog sb-flunky c 1234
```

### Optional flags

-   `server-id`: The ID of the target Artifactory service in JFrog CLI configuration (default: use default service). 
    Example: `--server-id=my-jfrog-service`.

-   `download-timeout`: Timeout of the Support Bundle download (default: 10 min). Example: `--download-timeout=15m`.

-   `retry-interval`: Waiting time between a failed download attempt and the next attempt (default: 5 sec). Example: 
    `--retry-interval=10s`.

-   `prompt-options`: Specify what is to be included in the created Support Bundle (default: use default Support Bundle 
    configuration). Example: `--prompt-options`.

-   `target-server-id`: The ID of the Artifactory service to which the Support Bundle will be uploaded (default: JFrog 
    "dropbox" service).

### Environment variables

None.

## Additional info

None.

## Release Notes

The release notes are available [here](RELEASE.md).
