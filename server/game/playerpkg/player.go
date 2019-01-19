package playerpkg

import (
	"math"
	"math/rand"
	"time"

	"github.com/giongto35/gowog/server/Message_proto"
	"github.com/giongto35/gowog/server/game/gameconst"
	"github.com/giongto35/gowog/server/game/shape"
	"github.com/golang/protobuf/ptypes"
)

type PlayerImpl struct {
	player   *Message_proto.Player
	name     string
	isEnable bool
	clientID int32
}

var dirX = map[Message_proto.Direction]float32{
	Message_proto.Direction_UP:    0,
	Message_proto.Direction_DOWN:  0,
	Message_proto.Direction_LEFT:  -1,
	Message_proto.Direction_RIGHT: 1,
}

var dirY = map[Message_proto.Direction]float32{
	Message_proto.Direction_UP:    -1,
	Message_proto.Direction_DOWN:  1,
	Message_proto.Direction_LEFT:  0,
	Message_proto.Direction_RIGHT: 0,
}

// NewPlayer returns new player
func NewPlayer(playerID int32, clientID int32, name string) Player {
	playerProto := &Message_proto.Player{}
	playerProto.Id = playerID
	// Reduce the starting size to 1/10 of the map
	// TODO: change gameconst.MapWidth, gameconst.MapHeight to map.GetWidth and map.GetHeight
	playerProto.X = rand.Float32() * gameconst.BlockWidth * gameconst.MapWidth
	playerProto.Y = rand.Float32() * gameconst.BlockHeight * gameconst.MapHeight

	playerProto.Size = 30
	playerProto.Health = 100
	playerProto.CurrentInputNumber = 0

	return &PlayerImpl{
		player:   playerProto,
		name:     name,
		clientID: clientID,
	}
}

// Move in dx and dy, dx and dy have to be calculated according to timeElapsed beforehand
func (p *PlayerImpl) Move(dx float32, dy float32) {
	p.player.X += dx
	p.player.Y += dy
}

// Shoot fires a bullet at x, y with direction dx, dy
func (p *PlayerImpl) Shoot(x float32, y float32, dx float32, dy float32) {
	// Normalize dx, dy
	dl := float32(math.Sqrt(float64(x*x + y*y)))
	dx = dx / dl
	dy = dy / dl

	// Update reload
	r, _ := ptypes.TimestampProto(time.Now().Add(time.Millisecond * gameconst.ReloadTime))
	p.player.NextReload = r
}

// GetPlayerProto returns player proto
func (p *PlayerImpl) GetPlayerProto() *Message_proto.Player {
	return p.player
}

// GetPlayerProto ...
func (p *PlayerImpl) GetPosition() Position {
	return Position{X: p.player.GetX(), Y: p.player.GetY()}
}

// GetName ...
func (p *PlayerImpl) GetName() string {
	return p.name
}

// GetClientID ...
func (p *PlayerImpl) GetClientID() int32 {
	return p.clientID
}

// GetID ...
func (p *PlayerImpl) GetID() int32 {
	return p.player.GetId()
}

// GetHealth ...
func (p *PlayerImpl) GetHealth() float32 {
	return p.player.GetHealth()
}

func (p *PlayerImpl) SetHealth(health float32) {
	p.player.Health = health
}

func (p *PlayerImpl) GetNextReload() time.Time {
	t, _ := ptypes.Timestamp(p.player.GetNextReload())
	return t
}

func (p *PlayerImpl) IsEnable() bool {
	return p.isEnable
}

func (p *PlayerImpl) SetEnable(enable bool) {
	p.isEnable = enable
}

func (p *PlayerImpl) SetCurrentInputNumber(curNumber int32) {
	// Skip older sequence number
	if p.player.CurrentInputNumber > curNumber {
		return
	}

	p.player.CurrentInputNumber = curNumber
}

// Box interface
func (p *PlayerImpl) GetLeft() float32 {
	return p.player.GetX() - gameconst.PlayerWidth/2
}

func (p *PlayerImpl) GetRight() float32 {
	return p.player.GetX() + gameconst.PlayerWidth/2
}

func (p *PlayerImpl) GetTop() float32 {
	return p.player.GetY() - gameconst.PlayerHeight/2
}

func (p *PlayerImpl) GetBottom() float32 {
	return p.player.GetY() + gameconst.PlayerHeight/2
}

func (p *PlayerImpl) GetSize() float32 {
	return p.player.Size
}

func (p *PlayerImpl) GetCircle() shape.Circle {
	return shape.Circle{X: p.player.GetX(), Y: p.player.GetY(), Radius: p.player.GetSize() / 2}
}

// GetRect returns a boundary rectangle of player
func (p *PlayerImpl) GetRect() shape.Rect {
	return shape.Rect{
		X1: p.player.GetX() - p.player.GetSize()/2,
		Y1: p.player.GetY() - p.player.GetSize()/2,
		X2: p.player.GetX() + p.player.GetSize()/2,
		Y2: p.player.GetY() + p.player.GetSize()/2,
	}
}
