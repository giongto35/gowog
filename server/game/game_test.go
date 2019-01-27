package game

import (
	"io/ioutil"
	"log"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/giongto35/gowog/server/Message_proto"
	"github.com/giongto35/gowog/server/game/common"
	"github.com/giongto35/gowog/server/game/gameconst"
	"github.com/giongto35/gowog/server/game/mappkg"
	"github.com/giongto35/gowog/server/game/objmanager"
	"github.com/giongto35/gowog/server/game/ws"
	"github.com/golang/protobuf/proto"
)

func initGame(bufferSize int) *gameImpl {
	var g = gameImpl{}
	hub := ws.NewHub()
	g.hub = hub
	gameconst.BufferSize = bufferSize

	// Setup Object manager
	g.destroyPlayerStream = make(chan common.DestroyPlayerEvent, bufferSize)
	g.newPlayerStream = make(chan common.NewPlayerEvent, bufferSize)
	g.inputStream = make(chan common.ProcessInputEvent, bufferSize)
	gameMap := mappkg.NewMap(gameconst.BlockWidth, gameconst.BlockHeight)
	g.objManager = objmanager.NewObjectManager(g.destroyPlayerStream, gameMap)

	go hub.Run()
	hub.BindGameMaster(&g)
	return &g
}

func moveMessage(clientID int32) []byte {
	msg := &Message_proto.ClientGameMessage{
		Message: &Message_proto.ClientGameMessage_MovePositionPayload{
			MovePositionPayload: &Message_proto.MovePosition{
				Id: clientID,
				Dx: rand.Float32()*2 - 1,
				Dy: rand.Float32()*2 - 1,
			},
		},
	}

	encodedMsg, _ := proto.Marshal(msg)
	return encodedMsg
}

func shootMessage(clientID int32) []byte {
	msg := &Message_proto.ClientGameMessage{
		Message: &Message_proto.ClientGameMessage_ShootPayload{
			ShootPayload: &Message_proto.Shoot{
				Id:       int64(clientID),
				PlayerId: clientID,
				X:        rand.Float32() * 100,
				Y:        rand.Float32() * 100,
				Dx:       rand.Float32()*2 - 1,
				Dy:       rand.Float32()*2 - 1,
			},
		},
	}

	encodedMsg, _ := proto.Marshal(msg)
	return encodedMsg
}

func playerRun(g *gameImpl, clientID int32) {
	g.NewPlayerConnect(clientID)
	for {
		g.ProcessInput(moveMessage(clientID))
		g.ProcessInput(shootMessage(clientID))
	}
}

func benchmarkGame(numcores int, numplayers int, bufferSize int, b *testing.B) {
	runtime.GOMAXPROCS(numcores)
	g := initGame(bufferSize)
	log.SetOutput(ioutil.Discard)

	for np := 0; np < numplayers; np++ {
		go playerRun(g, int32(np))
	}

	ticker := time.NewTicker(gameconst.RefreshRate * time.Millisecond)
	for n := 0; n < b.N; n++ {
		// Update loop

		select {
		case v := <-g.destroyPlayerStream:
			log.Println("Remove player", v)
			g.removePlayer(v.PlayerID, v.ClientID)
			log.Println("Remove player done", v)

		case v := <-g.newPlayerStream:
			log.Println("New player with clientID", v)
			g.newPlayerConnect(v.ClientID)
			log.Println("New player with clientID done", v)

		case v := <-g.inputStream:
			log.Println("Processs Message", v)
			g.processInput(v.Message)
			log.Println("Processs Message done", v)

		case <-ticker.C:
			g.Update()

		default:
		}
	}

}

func BenchmarkGame1Players(b *testing.B)                { benchmarkGame(1, 1, 500, b) }
func BenchmarkGame50Players(b *testing.B)               { benchmarkGame(1, 50, 500, b) }
func BenchmarkGame50PlayersMoreCores(b *testing.B)      { benchmarkGame(8, 50, 500, b) }
func BenchmarkGame1PlayersSize1(b *testing.B)           { benchmarkGame(1, 1, 1, b) }
func BenchmarkGame50PlayersSize1(b *testing.B)          { benchmarkGame(1, 50, 1, b) }
func BenchmarkGame50PlayersMoreCoresSize1(b *testing.B) { benchmarkGame(8, 50, 1, b) }
