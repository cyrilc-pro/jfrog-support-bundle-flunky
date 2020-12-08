![GitHub Workflow Status](https://img.shields.io/github/workflow/status/cyrilc-pro/jfrog-support-bundle-flunky/Go?logo=github)
[![Codecov](https://img.shields.io/codecov/c/github/cyrilc-pro/jfrog-support-bundle-flunky?label=codecov&logo=codecov&logoColor=fff)](https://codecov.io/gh/cyrilc-pro/jfrog-support-bundle-flunky)
[![Go Report](https://goreportcard.com/badge/github.com/cyrilc-pro/jfrog-support-bundle-flunky)](https://goreportcard.com/report/github.com/cyrilc-pro/jfrog-support-bundle-flunky)
[![Codacy grade](https://img.shields.io/codacy/grade/b286b95be72c4aa19de86f8c4a985f34?label=codacy&logo=codacy)](https://www.codacy.com/gh/cyrilc-pro/jfrog-support-bundle-flunky/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=cyrilc-pro/jfrog-support-bundle-flunky&amp;utm_campaign=Badge_Grade)

# JFrog Support Bundle Flunky

## About this plugin

This plugin helps with generating and exchanging a JFrog Support Bundle with JFrog Support.

## Installation with JFrog CLI

Since this plugin is currently not included in [JFrog CLI Plugins Registry](https://github.com/jfrog/jfrog-cli-plugins-reg),
it needs to be built and installed manually. Follow these steps to install and use this plugin with JFrog CLI.

1.  Make sure JFrog CLI is installed on you machine by running ```jfrog```. If it is not installed, [install](https://jfrog.com/getcli/) it.
2.  Create a directory named ```plugins``` under ```~/.jfrog/``` if it does not already exist.
3.  Clone this repository.
4.  CD into the root directory of the cloned project.
5.  Run ```go build -o sb-flunky``` to create the binary in the current directory.
6.  Copy the binary into the ```~/.jfrog/plugins``` directory.

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

```
jfrog sb-flunky support-case 1234
```

or

```
jfrog sb-flunky case 1234
```

or

```
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

## License

[Apache 2.0 License](LICENSE).

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fcyrilc-pro%2Fjfrog-support-bundle-flunky.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fcyrilc-pro%2Fjfrog-support-bundle-flunky?ref=badge_large)