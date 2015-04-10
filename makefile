test: test-stub test-docker

test-stub:
	go test -v github.com/dolaterio/dolaterio

test-docker:
	USE_DOCKER=1 go test -v github.com/dolaterio/dolaterio

dep-install:
	go get "github.com/fsouza/go-dockerclient"
	go get "github.com/gorilla/mux"
	go get "github.com/dancannon/gorethink"

run:
	go build && ./dolaterio

3rd-party-tools:
	docker run --restart always -d -p 8080:8080 -p 28015:28015 -p 29015:29015 --name dolaterio-rethinkdb rethinkdb:1.16
