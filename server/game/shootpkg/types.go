package shootpkg

import "github.com/giongto35/gowog/server/game/shape"
import "time"

type Shoot interface {
	GetShootAtTime(CurrentTime time.Time) Shoot
	GetShootObject() *ShootObject
	GetID() int64
	GetPlayerID() int32
	GetX() float32
	GetY() float32
	GetDX() float32
	GetDY() float32

	// Body interface
	GetPoint() shape.Point
}
