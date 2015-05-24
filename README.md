# Dolater.io

[![Circle CI](https://circleci.com/gh/dolaterio/dolaterio.svg?style=svg)](https://circleci.com/gh/dolaterio/dolaterio)

Dolater.io lets you execute background jobs on a remote docker server.

# Quick start

Dolater.io runs your jobs as docker images. Check out our example docker images like [dolaterio/dummy_worker](https://github.com/dolaterio/dummy_worker), [dolaterio/asciify](https://github.com/dolaterio/asciify) or [dolaterio/parrot](https://github.com/dolaterio/parrot).

Run rethinkdb, it's a dependency for dolater.io:
```
docker run \
  --restart always \
  -d \
  -p 8080:8080 \
  -p 28015:28015 \
  -p 29015:29015 \
  --name dolaterio-rethinkdb \
  rethinkdb:2.0
```

Do the same with redis, another dependency:
```
docker run \
  --restart always \
  -d \
  -p 6380:6379 \
  --name dolaterio-redis \
  redis:2.8
```

Then, run dolater.io:

```
docker run \
  -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -p 8080:8080 \
  -e "BINDING=0.0.0.0" \
  --link dolaterio-rethinkdb:rethinkdb \
  --link dolaterio-redis:redis \
  dolaterio/dolaterio
```

Now create a worker using our parrot docker image:

```
curl http://DOCKERHOST:8080/v1/workers -H "Content-Type: application/json" -X POST -d '{"docker_image": "dolaterio/parrot"}'
```

You'll get the worker json back in the response of that request. Use its `id` to create jobs:

```
curl http://DOCKERHOST:8080/v1/jobs -H "Content-Type: application/json" -X POST -d '{"worker_id": WORKER_ID, "stdin": "Hello world!"}'
```

It will return a new JSON containing an `id`. You can request dolater.io for the current state of the job:

```
curl http://DOCKERHOST:8080/v1/jobs/ID
```
