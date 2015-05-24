test:
	godep go test -v ./...

build:
	CGO_ENABLED=0 GOOS=linux go build -a -o dolater ./api && \
	docker build -t dolaterio/dolaterio . &&
	rm dolater

run:
	docker run \
		-it \
		--rm \
		--link dolaterio-rethinkdb:rethinkdb\
		--link dolaterio-redis:redis \
		-e "BINDING=0.0.0.0" \
		-p 7000:7000 \
		dolaterio/dolaterio

dep-install:
	go get github.com/tools/godep && GOPATH=$$(godep path) godep restore

3rd-party-tools:
	docker run --restart always -d -p 8080:8080 -p 28015:28015 -p 29015:29015 --name dolaterio-rethinkdb rethinkdb:2.0
	docker run --restart always -d -p 6380:6379 --name dolaterio-redis redis:2.8
