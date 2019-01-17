package shootpkg

import (
	"time"

	"github.com/giongto35/gowog/server/game/shape"
)

type ShootObject struct {
	ID        int64
	PlayerID  int32
	X         float32
	Y         float32
	DX        float32
	DY        float32
	StartTime time.Time
}

// NewShoot returns a new shoot object
func NewShoot(ID int64, playerID int32, x float32, y float32, dx float32, dy float32, startTime time.Time) Shoot {
	return &ShootObject{
		ID:        ID,
		PlayerID:  playerID,
		X:         x,
		Y:         y,
		DX:        dx,
		DY:        dy,
		StartTime: startTime,
	}
}

// GetShootAtTime to get the new position of the shoot at given time
func (s *ShootObject) GetShootAtTime(CurrentTime time.Time) Shoot {
	distance := float32(CurrentTime.Sub(s.StartTime).Seconds())
	return &ShootObject{
		X:         s.X + 1000*s.DX*distance,
		Y:         s.Y + 1000*s.DY*distance,
		DX:        s.DX,
		DY:        s.DY,
		StartTime: s.StartTime,
	}
}

func (s *ShootObject) GetPlayerID() int32 {
	return s.PlayerID
}

func (s *ShootObject) GetID() int64 {
	return s.ID
}

func (s *ShootObject) GetX() float32 {
	return s.X
}

func (s *ShootObject) GetY() float32 {
	return s.Y
}

func (s *ShootObject) GetDX() float32 {
	return s.DX
}

func (s *ShootObject) GetDY() float32 {
	return s.DY
}

// Consider return the whole object or just partial
func (s *ShootObject) GetShootObject() *ShootObject {
	return s
}

// Body shape
func (s *ShootObject) GetPoint() shape.Point {
	return shape.Point{
		X: s.X,
		Y: s.Y,
	}
}

// Note shootwidth
