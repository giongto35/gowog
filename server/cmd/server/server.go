package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/giongto35/gowog/server/game"
	"github.com/giongto35/gowog/server/game/ws"
	"github.com/gorilla/websocket"
	"github.com/pkg/profile"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")
var logfile = flag.String("logfile", "", "Log file")
var cpuprofile = flag.Bool("cpuprofile", false, "Enable CPUProfile")
var memprofile = flag.Bool("memprofile", false, "Enable MemProfile")
var disablelog = flag.Bool("disablelog", false, "Disable log")
var clientBuild = flag.String("prod", "", "is production")

var upgrader = websocket.Upgrader{} // use default options
var hub = ws.NewHub()
var gameMaster = game.NewGame(hub)

// serveWs handles websocket requests from the peer.
func connect(w http.ResponseWriter, r *http.Request) {
	clientID := ws.NewClient(upgrader, hub, w, r)
	gameMaster.NewPlayerConnect(clientID)
}

func main() {
	// Running on one core only
	runtime.GOMAXPROCS(2)
	flag.Parse()
	if *disablelog {
		log.SetOutput(ioutil.Discard)
	}
	// CPU profile
	if *cpuprofile {
		log.Println("Profiling CPU")
		defer profile.Start().Stop()
	}
	// Memory profile
	if *memprofile {
		log.Println("Profiling MemProfile")
		defer profile.Start(profile.MemProfile).Stop()
	}

	// If there is clientBuild flag, we return the client build for index
	if *clientBuild != "" {
		log.Println("loading file from ", *clientBuild)
		http.Handle("/", http.FileServer(http.Dir(*clientBuild)))

	}

	// Log to file
	if *logfile != "" {
		log.Println("Write to logfile", *logfile)
		f, err := os.OpenFile(*logfile, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		log.SetOutput(f)
		defer f.Close() // HTTP setup
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	log.SetFlags(0)
	// Websocket endpoint
	http.HandleFunc("/game/", connect)

	fmt.Println("Listening to ", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
	fmt.Println("Stop Listening to ", addr)
}
