# Kwsk (Knative OpenWhisk)

## Prerequisites

You'll need a recent [`go`](https://golang.org/doc/install) and
[`dep`](https://github.com/golang/dep). You'll also want this repository
checked somewhere under `$GOPATH/src`.

## Running the server

    dep ensure
    go run ./cmd/kwsk-server/main.go --port 8080

## Testing the server

No automated testing yet, but you can hit the thing via curl like:

    curl http://127.0.0.1:8080/api/v1/namespaces/foo/actions

## Implementing the server

The server just contains a but of stubs and unimplemented code right
now. Read through the [go-swagger docs](https://goswagger.io/generate/server.html)
to see how to get started with actual implementations.

## Generating the server

This server was generated using [go-swagger](https://goswagger.io/)
with the command:

    swagger generate server -A Kwsk -P models.Principal -f apiv1swagger.json
