test: test-stub test-docker

test-stub:
	go test -v github.com/dolaterio/dolaterio

test-docker:
	USE_DOCKER=1 go test -v github.com/dolaterio/dolaterio

dep-install:
	go get "github.com/fsouza/go-dockerclient"
