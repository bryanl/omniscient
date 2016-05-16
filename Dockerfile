FROM golang:1.6.2
ADD . /omniscient
ENV buildNum="dev"
RUN go get github.com/constabulary/gb/... \
  && cd /omniscient \
    && gb build -ldflags "-X omniscient.version=$buildNum"
ENTRYPOINT /omniscient/bin/omniscient
EXPOSE 8080
