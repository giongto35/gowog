package mappkg

import (
	"github.com/giongto35/gowog/server/Message_proto"
	"github.com/giongto35/gowog/server/game/common"
	"github.com/giongto35/gowog/server/game/shape"
)

type Map interface {
	ToProto() *Message_proto.Map
	GetWidth() float32
	GetHeight() float32
	GetNumCols() int
	GetNumRows() int
	GetStartPoint() common.Point
	GetEndPoint() common.Point
	IsCollide(x float32, y float32) bool
	GetRectBlocks() []shape.Rect
}

type Point struct {
	X float32
	Y float32
}
