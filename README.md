# Dolater.io

Dolater.io lets you execute background jobs on a remote docker server.

# How to run it

You'll need [docker-compose](https://docs.docker.com/compose/) to run the services and its dependencies easily. To run all services at once use, run:

```
docker-compose up -d --no-recreate
```

This will run the API server as well as one dolater.io worker. You can always scale the amount of workers by using `docker-compose scale worker=N` command.

Now it's all ready to use.

# Simple Example

Since dolater.io is running in docker, you'll need to know your docker host IP address to access it. If you use boot2docker, run `boot2docker ip` to find out.

Create a worker using our parrot docker image:

```
curl http://DOCKERHOST:7000/v1/workers -H "Content-Type: application/json" -X POST -d '{"docker_image": "dolaterio/parrot"}'
```

You'll get a JSON response back with the information of the worker you just created. Use its `id` to create jobs on it:

```
curl http://DOCKERHOST:7000/v1/jobs -H "Content-Type: application/json" -X POST -d '{"worker_id": WORKER_ID, "stdin": "Hello world!"}'
```

It will return a new JSON containing, between others, the `id` of the job. You can request dolater.io for the current state of the job:

```
curl http://DOCKERHOST:7000/v1/jobs/JOB_ID
```
