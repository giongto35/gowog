package mappkg

import (
	"github.com/giongto35/gowog/server/Message_proto"
	"github.com/giongto35/gowog/server/game/shape"
)

type Map interface {
	ToProto() *Message_proto.Map
	GetWidth() float32
	GetHeight() float32
	GetNumCols() int
	GetNumRows() int
	IsCollide(x float32, y float32) bool
	GetRectBlocks() []shape.Rect
}
