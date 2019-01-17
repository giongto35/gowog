# GOWOG Server

## Summary

Golang based multipler game server. The game can tollerate high concurrency throughput.

# Background
Not like webserver when the each requests don't share the state and can have some latency acceptance point. Game server involves modifying sharing memory state to achieve smooth performance.

Golang provides a very elegant solution to solve high concurrency problem by goroutine and channel while still maintaining good running performance.

# Installation
[**Main**](..)

This will run web server in the terminal, which listens to port 8080

# Codebase
```
├── server
│   ├── buildwall.js
│   ├── cmd
│   │   └── server
│   │       └── server.go: Entrypoint running server
│   ├── game
│   │   ├── common
│   │   ├── config
│   │   │   └── 1.map: Map represented 0 and 1
│   │   ├── eventmanager
│   │   ├── gameconst
│   │   ├── game.go
│   │   ├── mappkg
│   │   ├── objmanager
│   │   ├── playerpkg
│   │   ├── shape
│   │   ├── shootpkg
│   │   ├── types.go
│   │   └── ws
│   │       ├── wsclient.go
│   │       └── wshub.go
│   ├── generate.sh: Generate protobuf for server + client + AI environment
│   ├── message.proto
│   └── Message_proto
│       └── message.pb.go
├── Dockerfile
└── run_local.sh
```

# Credits
The server websocket design is based on Gorila websocket chat example
https://github.com/gorilla/websocke/blob/master/examples/chat
