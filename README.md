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

Then, run dolater.io:

```
docker run \
  -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -p 8080:8080 \
  --link dolaterio-rethinkdb:rethinkdb \
  dolaterio/dolaterio \
  bash -c "\
    RETHINKDB_ADDRESS="\$RETHINKDB_PORT_28015_TCP_ADDR:\$RETHINKDB_PORT_28015_TCP_PORT" /gopath/bin/dolaterio --bind 0.0.0.0\
  "
```

Now, ready to queue a job! Send the following request to your instance:
```
curl http://DOCKERHOST:8080/v1/jobs -H "Content-Type: application/json" -X POST -d '{"docker_image": "dolaterio/parrot", "stdin": "Hello world!"}'
```

It will return a new JSON containing an `id`. You can request dolater.io for the current state of the job:

```
curl http://DOCKERHOST:8080/v1/jobs/ID
```
