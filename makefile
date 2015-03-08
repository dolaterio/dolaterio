default:
	go test -v github.com/dolaterio/dolaterio/runner

dep-install:
  go get github.com/fsouza/go-dockerclient
