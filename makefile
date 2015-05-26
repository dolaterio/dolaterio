migrate_db:
	godep go run ./migrate/main.go

test: migrate_db
	godep go test -v ./...

dolater.bin:
	CGO_ENABLED=0 GOOS=linux go build -a -o dolater.bin ./api
migrate.bin:
	CGO_ENABLED=0 GOOS=linux go build -a -o migrate.bin ./api
clean:
	if [ -a dolater.bin ] ; \
	then \
	     rm dolater.bin ; \
	fi; \
	if [ -a migrate.bin ] ; \
	then \
	     rm migrate.bin ; \
	fi;
build: clean dolater.bin migrate.bin
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

run:
	docker run \
		-it \
		--rm \
		--link dolaterio-rethinkdb:rethinkdb \
		--link dolaterio-redis:redis \
		-e "BINDING=0.0.0.0" \
		-p 7000:7000 \
		dolaterio/dolaterio

dep-install:
	go get github.com/tools/godep && GOPATH=$$(godep path) godep restore

3rd-party-tools:
	docker run --restart always -d -p 8080:8080 -p 28015:28015 -p 29015:29015 --name dolaterio-rethinkdb rethinkdb:2.0
	docker run --restart always -d -p 6380:6379 --name dolaterio-redis redis:2.8
