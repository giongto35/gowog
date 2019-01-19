# GOWOG, Open source GOlang Web-based Online Game

GOWOG is a multiplayer web game written in Golang. The game can serve high number of players concurrently by following Golang concurrency model.

# Try the game

## Remote

http://game.giongto35.com

## Local Docker environment

You can try the game on local by running `./run_local.sh`. It will build a docker environment and run the game on "localhost:8080".
You can continue the development by exec into the docker.

# Development

The game contains two part: Server and Client. Server uses Golang and Client uses Node.JS.

## Script

### Linux
  * ./install_ubuntu.sh

### MACOS
  * ./install_mac.sh

## Manual Install

### Server

Install Golang https://golang.org/doc/install
Install dependencies
  * go get github.com/gorilla/websocket
  * go get github.com/golang/protobuf/protoc-gen-go
  * go get github.com/pkg/profile
Run the server. The server will listen at port 8080
  * go run cs2dserver/cmd/server/* 
 
### Client

Install NodeJS https://nodejs.org/en/download/
  * npm install
  * go get github.com/golang/protobuf/protoc-gen-go
  * go get github.com/pkg/profile
Run the client. The client will listen at port 3000. env.HOST_IP is the host of server
  * npm run dev -- --env.HOST_IP=localhost:8080
  * open the browser "localhost:3000"
 
### Communication convention

Communication packages between server and client is based on protobuf. Install protoc for protobuf generate.
  * http://google.github.io/proto-lens/installing-protoc.html

Everytime you change package singature in message.proto. Run 
  * ./server/generate.sh

# Architecture
![Techstack](document/images/techstack.jpg)

## Code structure
[**Frontend**](client)

[**Backend**](server)

# FAQ

### Why we need Golang for multiplayer game?

Building a massively multiplayer game is very difficult and it's currently overlooked. You have to ensure the latency is acceptable, handle shared states concurrently and allow it to scale vertically. Golang provides a very elegant concept to handle concurrency with goroutine and channel.

### Why the gameplay is so simple and frontend codebase is so unorganized?

The gameplay is mainly for demonstration purpose. My goal is to keep the game simple as current and scale number of players vertically while maintaining good latency (< 100ms). I welcome all of your ideas to make the game more scalable.

However, I'm still open to Graphic improvement and client codebase refactor. I would love to see some particles burst or glow, motion effects. Remember that different people have different perspective on how the game should look, so we need full consensus before making ui change.

### Why the game only runs on single core?

The game indeed can run well on multi-core parallelly. After some comparision, running on multi-core showed the slower performance due to high channel and lock contention.

We need a better design to reduce context switch and contention.

### If the game runs on single core, why needs to you channel? Why don't fully go with NodeJS for server and callback model?

Remember, concurrency is not parallelism. Context switch can happen everytime. GoRoutine and GoChannel is very elegant solution to deal with concurrency. And it's easier and more intuitive than with Callback model (a.k.a Callback hell).

### Why protobuf?

To optimize package size, we need to compress it into binary. Protobuf offers fast language-neutral serialization and deserilization, so Golang server can communicate with JS client in an optimal way.

We can consider faster serilization format like Cap'n Proto or FlatBuffers.

### Why the multiplayer game is not room-based (less than 20 players per game)

Room-based is less challenging because each room is totally separated from each other and we can scale the game easier by adding more instances. This repo is to illustrate how many players the game can serve on one multi-core single machine.
And of course my ultimate goal is to make the non room-based game scale vertically by adding more instances.

# Credits

https://github.com/gorilla/websocke/blob/master/examples/chat
https://github.com/RenaudROHLINGER/phaser-es6-webpack
https://github.com/huytd/agar.io-clone

# Contributor

Nguyen Huu Thanh  
