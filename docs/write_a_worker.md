# How to write a dolater.io worker

This document will guide you through the creation of a *dolater.io* worker.

##Code

We have some [templates](https://github.com/dolaterio/templates) in different languages with the simplest structure for your worker.

You can also have a look at some of our internal workers as a reference:
- [Simple image resizer](https://github.com/dolaterio/simple_image_resizer)
- [Asciify](https://github.com/dolaterio/asciify)
- [Parrot](https://github.com/dolaterio/parrot)
- [Dummy worker](https://github.com/dolaterio/dummy_worker)
- [Webhook caller](https://github.com/dolaterio/webhook_caller)

## Configuration via environment variables

Dolater.io lets you set your environment variables in different ways, so in your worker code it's recommended to use them to get configuration values.

Check the API documentation to see how environment variables are defined.

## Send job-specific data via STDIN

All job input will be sent as STDIN stream. Make sure your worker uses it as your main data input.

If you're not sure how to read from the standard input, check our [templates](https://github.com/dolaterio/templates) or drop us an [email](mailto:admin@dolater.io) if you're having trouble and we'll help you out.

## Build a docker container with your worker

You'll need to wrap your worker in a docker image by using a _Dockerfile_. The easiest way to get up to speed on creating a _Dockerfile_ and building a docker image is by reading the [Dockerfile reference](https://docs.docker.com/reference/builder/). There are plenty of Dockerfiles you can use as a reference. If you need help, drop us an [email](mailto:admin@dolater.io).

dolater.io runs the image as it is, so remember to keep in `CMD` the exact command to execute your worker.

## Monitoring STDOUT & STDERR when a job runs

To test if your image, run it in your local docker with `echo YOUR_STDIN | docker run -i -e ENV_1=a -e ENV_2=b your_image` replacing `YOUR_STDIN` with the standard input you want to send to the process and replacing the environment variables with whatever your image needs.

When a job runs on dolater.io, we capture the standard output and standard error, so whatever you see when you run the process will be available on your dolater.io job results.

## Publish your image

At this point dolater.io only supports public images stored at [the public docker hub](https://hub.docker.com/). See [How to work with Docker Hub](http://docs.docker.com/userguide/dockerrepos/) for a guide on publishing your images.
