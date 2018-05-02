FROM golang:1.9-alpine3.7 as builder

COPY . /go/src/github.com/flaccid/snstxtr

WORKDIR /go/src/github.com/flaccid/snstxtr/cmd/snstxtr

RUN apk add --update --no-cache git gcc musl-dev && \
    go get ./... && \
    CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o /opt/bin/snstxtr .

FROM centurylink/ca-certs

COPY --from=builder /opt/bin/snstxtr /opt/bin/snstxtr

ENTRYPOINT ["/opt/bin/snstxtr"]
