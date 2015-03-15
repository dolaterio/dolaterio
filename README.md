# Dolater.io

[![Circle CI](https://circleci.com/gh/dolaterio/dolaterio.svg?style=svg)](https://circleci.com/gh/dolaterio/dolaterio)

Dolater.io is a project to execute background jobs in a very scalable way.

# Architecture

Jobs are mostly docker images. Every job to run it's a docker container. The job specifies the docker image to run to process the job. A very simple docker image used for testing can be found here: [dolaterio/dummy_worker](https://github.com/dolaterio/dummy_worker)

At the end of the stack of _dolater.io_ there's a job runner. Each runner runs up to N jobs simultaneously. The runner runs the container, waits for it to finish and gets its results.

In the stack, right before the job runner there's queue with the jobs pending to execute. Multiple runners can be consuming the queue. The queue system will make sure that only one runner gets one specific job.

Once the runner finishes processing a queue message, it'll write the results to a different queue.

At this point the job runners are running behind queues. Multiple processes can queue jobs and multiple process can consume job results.
