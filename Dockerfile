FROM golang:1.6
ADD . /go/src/github.com/bryanl/omniscient
ENTRYPOINT /go/bin/omniscient
EXPOSE 8080