From golang:1.11

RUN mkdir -p /go/src/github.com/giongto35/gowog/server/
RUN mkdir -p /go/gowog/game/config

COPY ./server /go/src/github.com/giongto35/gowog/server/
COPY ./server/game/config /go/server/game/config

# Install server dependencies
RUN go get github.com/gorilla/websocket
RUN go get github.com/golang/protobuf/protoc-gen-go
RUN go get github.com/pkg/profile

RUN go install github.com/giongto35/gowog/server/cmd/server/

# Need argument, default is localhost
ARG HOSTNAME=localhost

# build client
COPY ./client ./client

RUN curl -SL https://deb.nodesource.com/setup_10.x | bash
RUN apt-get install nodejs
RUN npm --prefix ./client install ./client
WORKDIR client
RUN HOST_IP=$HOSTNAME npm run deploy
WORKDIR ./..

EXPOSE 8080
