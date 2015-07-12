test: migrate_db
	godep go test -v ./...

run_api:
	docker-compose up -d --no-recreate api
run_worker:
	docker-compose up -d --no-recreate worker

migrate_db:
	docker-compose run --rm worker /migrate
dev_migrate_db:
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
build_image: clean api.bin migrate.bin worker.bin
	docker build -t dolaterio/dolaterio .

dep-install:
	go get github.com/tools/godep && GOPATH=$$(godep path) godep restore
