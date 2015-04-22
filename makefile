test: test-core test-api
test-core:
	godep go test -v github.com/dolaterio/dolaterio/core
test-api:
	godep go test -v github.com/dolaterio/dolaterio/api

dep-install:
	go get github.com/tools/godep && GOPATH=$$(godep path) godep restore

build:
	go get github.com/dolaterio/dolaterio

run: build
	$$GOPATH/bin/dolaterio

3rd-party-tools:
	docker run --restart always -d -p 8080:8080 -p 28015:28015 -p 29015:29015 --name dolaterio-rethinkdb rethinkdb:2.0

build-dist:
	docker build -t dolaterio .

run-dist:
	docker run \
	  -v /var/run/docker.sock:/var/run/docker.sock \
	  --rm \
	  -p 8081:8080 \
	  --link dolaterio-rethinkdb:rethinkdb \
	  dolaterio
