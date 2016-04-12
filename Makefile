test: dev_migrate_db
	godep go test -v ./...

run:
	docker-compose up -d --no-recreate
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
	rm -f *.bin

build_image: clean api.bin migrate.bin worker.bin
	docker build -t dolaterio/dolaterio .

dep-install:
	go get github.com/tools/godep && GOPATH=$$(godep path) godep restore

run_dev_dependencies:
	docker run --restart always -d -p 8080:8080 -p 28015:28015 -p 29015:29015 --name dolaterio-rethinkdb rethinkdb:2.0
	docker run --restart always -d -p 6380:6379 --name dolaterio-redis redis:2.8
