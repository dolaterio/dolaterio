test: migrate_db
	godep go test -v ./...

migrate_db:
	godep go run ./migrate/main.go

api.bin:
	CGO_ENABLED=0 GOOS=linux go build -a -o api.bin ./api
worker.bin:
	CGO_ENABLED=0 GOOS=linux go build -a -o worker.bin ./worker
migrate.bin:
	CGO_ENABLED=0 GOOS=linux go build -a -o migrate.bin ./migrate
clean:
	if [ -a api.bin ] ; \
	then \
	     rm api.bin ; \
	fi; \
	if [ -a migrate.bin ] ; \
	then \
	     rm migrate.bin ; \
	fi;
	if [ -a worker.bin ] ; \
	then \
	     rm worker.bin ; \
	fi;
build: clean api.bin migrate.bin worker.bin
	docker build -t dolaterio/dolaterio .

run_migrate:
	docker run \
		-it \
		--rm \
		--link dolaterio-rethinkdb:rethinkdb \
		--link dolaterio-redis:redis \
		-e "BINDING=0.0.0.0" \
		-p 7000:7000 \
		dolaterio/dolaterio
		/migrate

run_api:
	docker run \
		--rm \
		-d \
		--link dolaterio-rethinkdb:rethinkdb \
		--link dolaterio-redis:redis \
		-e "BINDING=0.0.0.0" \
		-p 7000:7000 \
		dolaterio/dolaterio \
		/api

run_worker:
	docker run \
		--rm \
		-d \
		--link dolaterio-rethinkdb:rethinkdb \
		--link dolaterio-redis:redis \
		dolaterio/dolaterio

dep-install:
	go get github.com/tools/godep && GOPATH=$$(godep path) godep restore

3rd-party-tools:
	docker run --restart always -d -p 8080:8080 -p 28015:28015 -p 29015:29015 --name dolaterio-rethinkdb rethinkdb:2.0
	docker run --restart always -d -p 6380:6379 --name dolaterio-redis redis:2.8
