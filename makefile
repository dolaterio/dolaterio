test:
	godep go test -v ./...

build:
	CGO_ENABLED=0 GOOS=linux go build -a -o dolater ./api && \
	docker build -t dolaterio/dolaterio .

run:
	docker run -it --rm -v /Users/albert/workspace/gocode/src/github.com/dolaterio/dolaterio/config.yml dolaterio/dolaterio

dep-install:
	go get github.com/tools/godep && GOPATH=$$(godep path) godep restore

3rd-party-tools:
	docker run --restart always -d -p 8080:8080 -p 28015:28015 -p 29015:29015 --name dolaterio-rethinkdb rethinkdb:2.0
	docker run --restart always -d -p 6380:6379 --name dolaterio-redis redis:2.8
