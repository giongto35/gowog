package playerpkg

import (
	"time"

	"github.com/giongto35/gowog/server/Message_proto"
	"github.com/giongto35/gowog/server/game/common"
	"github.com/giongto35/gowog/server/game/shape"
)

type Player interface {
	Move(dx float32, dy float32)
	Shoot(x float32, y float32, dx float32, dy float32)

	GetPlayerProto() *Message_proto.Player

	GetPosition() common.Point
	SetPosition(common.Point)
	GetName() string
	GetID() int32
	GetClientID() int32
	GetNextReload() time.Time
	SetHealth(health float32)
	GetHealth() float32
	AddScore()

	IsEnable() bool
	SetEnable(enable bool)

	// Box interface
	GetTop() float32
	GetBottom() float32
	GetLeft() float32
	GetRight() float32

	// Client related Data
	SetCurrentInputNumber(int32)

	// Body interface
	GetCircle() shape.Circle
}
