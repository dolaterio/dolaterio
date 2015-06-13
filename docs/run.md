# Run dolater.io

Dolater.io uses Redis and Rethinkdb. Run them in docker using the following commands:

```bash
docker run \
    --restart always \
    -d \
    -p 8080:8080 -p 28015:28015 -p 29015:29015 \
    --name dolaterio-rethinkdb \
    rethinkdb:2.0
docker run \
    --restart always \
    -d -p 6380:6379 \
    --name dolaterio-redis \
    redis:2.8
```

Once those dependencies are up and running, you'll need to run the migration tool:

```bash
docker run \
    --rm \
    --link dolaterio-rethinkdb:rethinkdb \
    --link dolaterio-redis:redis \
    dolaterio/dolaterio \
    /migrate
```

And now dolater.io is ready to run. You'll need to run the API and at least one worker.
To run the API do the following:

```bash
docker run \
    -d \
    --restart always \
    --link dolaterio-rethinkdb:rethinkdb \
    --link dolaterio-redis:redis \
    -e "BINDING=0.0.0.0" \
    -p 7000:7000 \
    --name dolaterio \
    dolaterio/dolaterio \
    /api
```

To run a worker do the following:

```bash
docker run \
    -d \
    --restart always \
    --link dolaterio-rethinkdb:rethinkdb \
    --link dolaterio-redis:redis \
    -v /var/run/docker.sock:/var/run/docker.sock \
    dolaterio/dolaterio \
    /worker
```
