FROM golang:1.6.2
ADD . /omniscient
RUN go get github.com/constabulary/gb/... \
  && cd /omniscient \
  && gb build
ENTRYPOINT /omniscient/bin/omniscient
EXPOSE 8080