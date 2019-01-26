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

func initGame() *gameImpl {
	var g = gameImpl{}
	hub := ws.NewHub()
	g.hub = hub

	// Setup Object manager
	g.eventStream = make(chan interface{})
	gameMap := mappkg.NewMap(gameconst.BlockWidth, gameconst.BlockHeight)
	g.objManager = objmanager.NewObjectManager(g.eventStream, gameMap)

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

func benchmarkGame(numcores int, numplayers int, b *testing.B) {
	runtime.GOMAXPROCS(numcores)
	g := initGame()
	log.SetOutput(ioutil.Discard)

	for np := 0; np < numplayers; np++ {
		go playerRun(g, int32(np))
	}

	ticker := time.NewTicker(gameconst.RefreshRate * time.Millisecond)
	for n := 0; n < b.N; n++ {
		// Update loop

		select {
		case e := <-g.eventStream:
			switch v := e.(type) {
			case common.DestroyPlayerEvent:
				g.removePlayer(v.PlayerID, v.ClientID)

			case common.NewPlayerEvent:
				g.newPlayerConnect(v.ClientID)

			case common.ProcessInputEvent:
				g.processInput(v.Message)
			}

		case <-ticker.C:
			g.Update()

		default:
		}
	}

}

func BenchmarkGame1Players(b *testing.B)           { benchmarkGame(1, 1, b) }
func BenchmarkGame50Players(b *testing.B)          { benchmarkGame(1, 50, b) }
func BenchmarkGame50PlayersMoreCores(b *testing.B) { benchmarkGame(8, 50, b) }
