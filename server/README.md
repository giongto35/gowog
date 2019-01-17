# GOWOG Server

## Summary

Golang based multipler game server. The game can tollerate high concurrency throughput.

# Background
Not like webserver when the each requests don't share the state and can have some latency acceptance point. Game server involves modifying sharing memory state to achieve smooth performance.

Golang provides a very elegant solution to solve high concurrency problem by goroutine and channel while still maintaining good running performance.

# Setup

## 1. Install golang

## 2. Pull this repo and put to src folder

```
go get github.com/giongto35/gowog -u
``` 

## 2. Install dependencies:

```
go get github.com/gorilla/websocket
go get github.com/golang/protobuf/protoc-gen-go
go get github.com/pkg/profile
``` 

Run:

```go run server/cmd/server/*```

This will run web server in the terminal, which listens to port 8080

# Credits
The server websocket design is based on Gorila websocket chat example
https://github.com/gorilla/websocke/blob/master/examples/chat
