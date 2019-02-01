package common

import	"github.com/giongto35/gowog/server/game/ws"

// DestroyPlayerEvent is event sent from objManager to game master
type DestroyPlayerEvent struct {
	ClientID int32
	PlayerID int32
}

// NewPlayerEvent is event sent from objManager to game master
type NewPlayerEvent struct {
	Client ws.Client
}

// ProcessInputEvent is game input event sent from objManager to game master
type ProcessInputEvent struct {
	Message []byte
}
