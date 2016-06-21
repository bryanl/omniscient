FROM golang:1.6-alpine
ADD . /omniscient
ENV buildNum="dev"
RUN apk --update add git
RUN go get github.com/constabulary/gb/... \
  && cd /omniscient \
  && gb build -ldflags "-X omniscient.version=$buildNum"
ENTRYPOINT /omniscient/bin/omniscient
EXPOSE 8080
