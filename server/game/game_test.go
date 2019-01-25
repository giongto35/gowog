package game

import (
	"fmt"
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
	g.eventStream = make(chan interface{}, 500)
	gameMap := mappkg.NewMap(gameconst.BlockWidth, gameconst.BlockHeight)
	g.objManager = objmanager.NewObjectManager(g.eventStream, gameMap)

	go hub.Run()
	hub.BindGameMaster(&g)
	return &g
}

func benchmarkGame(b *testing.B) {
	//gi := NewGame(hub)
	//g := gi.(*gameImpl)
	g := initGame()

	g.NewPlayerConnect(0)
	fmt.Println("1")
	go func() {
		for i := 0; i <= 1000; i++ {
			msg := &Message_proto.ClientGameMessage{
				Message: &Message_proto.ClientGameMessage_MovePositionPayload{
					MovePositionPayload: &Message_proto.MovePosition{
						Id: 0,
						Dx: 1,
						Dy: 1,
					},
				},
			}

			encodedMsg, _ := proto.Marshal(msg)
			g.ProcessInput(encodedMsg)
		}
	}()

	fmt.Println("2")
	ticker := time.NewTicker(gameconst.RefreshRate * time.Millisecond)
	for n := 0; n < b.N; n++ {
		// Update loop

		select {
		case e := <-g.eventStream:
			switch v := e.(type) {
			case common.DestroyPlayerEvent:
				fmt.Println("Remove player", v)
				g.removePlayer(v.PlayerID, v.ClientID)

			case common.NewPlayerEvent:
				fmt.Println("New player with clientID", v)
				g.newPlayerConnect(v.ClientID)

			case common.ProcessInputEvent:
				fmt.Println("Processs Message", v)
				g.processInput(v.Message)
			}

		case <-ticker.C:
			g.Update()

		default:
		}
	}

}

func BenchmarkGame(b *testing.B) { benchmarkGame(b) }
