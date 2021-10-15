FROM golang:alpine3.13 AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src
COPY . /go/src

RUN GO111MODULE=on go get ./...
RUN GO111MODULE=on GOOS=linux go build -o bin/sitemap github.com/skhlv/sitemap


FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/bin /go/bin
ENV PORT = 3000
EXPOSE $PORT

ENTRYPOINT /go/bin/sitemap
