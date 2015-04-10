FROM google/golang:1.4

# Getting it is very slow, adding it here to make use of caching.
RUN go get "github.com/fsouza/go-dockerclient"
RUN go get "github.com/gorilla/mux"
RUN go get "github.com/dancannon/gorethink"

WORKDIR /gopath/src/github.com/dolaterio/dolaterio
ADD . /gopath/src/github.com/dolaterio/dolaterio

RUN make dep-install
RUN go get github.com/dolaterio/dolaterio

EXPOSE 8080

ENV BIND_IP 0.0.0.0

CMD ["/gopath/bin/dolaterio"]

