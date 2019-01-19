package eventmanager

import (
	"time"

	"github.com/giongto35/gowog/server/game/playerpkg"
	"github.com/giongto35/gowog/server/game/shootpkg"
)

type EventManager interface {
	AddEvent(clientID int32) playerpkg.Player
	RegisterShoot(player playerpkg.Player, x float32, y float32, dx float32, dy float32, startTime time.Time) shootpkg.Shoot
}
