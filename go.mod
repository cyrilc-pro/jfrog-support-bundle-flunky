module github.com/jfrog/jfrog-support-bundle-flunky

go 1.14

require (
	github.com/AlecAivazis/survey/v2 v2.2.3
	github.com/docker/go-connections v0.4.0
	github.com/google/go-cmp v0.5.4
	github.com/jfrog/jfrog-cli-core v0.0.1
	github.com/jfrog/jfrog-client-go v0.16.0
	github.com/stretchr/testify v1.6.1
	github.com/testcontainers/testcontainers-go v0.9.0
)

replace github.com/jfrog/jfrog-cli-core => github.com/jfrog/jfrog-cli-core v1.1.2
