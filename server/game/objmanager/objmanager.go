package objmanager

import (
	"log"
	"math"
	"time"

	"github.com/giongto35/gowog/server/game/common"
	"github.com/giongto35/gowog/server/game/gameconst"
	"github.com/giongto35/gowog/server/game/mappkg"
	"github.com/giongto35/gowog/server/game/playerpkg"
	"github.com/giongto35/gowog/server/game/shape"
	"github.com/giongto35/gowog/server/game/shootpkg"
)

type objManager struct {
	players   map[int32]playerpkg.Player
	shoots    map[playerpkg.Player][]shootpkg.Shoot
	gameMap   mappkg.Map
	numPlayer int32
	numShoot  int64

	// Stream for all game event
	destroyPlayerStream chan common.DestroyPlayerEvent
}

func NewObjectManager(eventStream chan common.DestroyPlayerEvent, gameMap mappkg.Map) ObjectManager {
	objManager := objManager{}
	objManager.players = map[int32]playerpkg.Player{}
	objManager.shoots = map[playerpkg.Player][]shootpkg.Shoot{}
	objManager.destroyPlayerStream = eventStream
	objManager.gameMap = gameMap

	return &objManager
}

func (m *objManager) RegisterPlayer(clientID int32, name string) playerpkg.Player {
	// Create new player
	m.numPlayer++
	playerID := m.numPlayer - 1

	// Keep generate player until it not collide with gameMap
	// TODO: check first and assign
	var tempPlayer playerpkg.Player
	for {
		tempPlayer = playerpkg.NewPlayer(playerID, clientID, name)
		if m.IsValidPosition(tempPlayer.GetCircle()) {
			break
		}
	}
	// Assign tempPlayer to player
	m.players[playerID] = tempPlayer
	player := m.players[playerID]
	player.SetEnable(true)
	log.Println("Register player", player.GetPlayerProto().GetId(), player.GetName())

	return player
}

// RangePlayers is to concurrently iterate through players list
//func (m *objManager) RangePlayers(f func(player playerpkg.Player)) {
//m.PlayersLock.Lock()
//defer m.PlayersLock.Unlock()

//for _, player := range m.players {
//f(player)
//}
//}

// RegisterShoot creates a new shoot based by time, direction and the firing time
func (m *objManager) RegisterShoot(player playerpkg.Player, x float32, y float32, dx float32, dy float32, startTime time.Time) shootpkg.Shoot {
	m.numShoot++
	shootID := m.numShoot - 1
	shoot := shootpkg.NewShoot(shootID, player.GetID(), x, y, dx, dy, startTime)
	m.shoots[player] = append(m.shoots[player], shoot)

	return shoot
}

// Remove Shoot from player
func (m *objManager) RemoveShoot(shoot shootpkg.Shoot) {
	player, ok := m.GetPlayerByID(shoot.GetPlayerID())
	if !ok {
		return
	}
	// Finding the shoots in all the shoot from the player and remove from m.shoots
	for i, s := range m.shoots[player] {
		if s == shoot {
			m.shoots[player] = append(m.shoots[player][:i], m.shoots[player][i+1:]...)
		}
	}
}

func (m *objManager) GetMap() mappkg.Map {
	return m.gameMap
}

func (m *objManager) GetPlayers() []playerpkg.Player {
	players := []playerpkg.Player{}
	for _, player := range m.players {
		players = append(players, player)
	}

	return players
}

func (m *objManager) GetShoots() []shootpkg.Shoot {
	shoots := []shootpkg.Shoot{}
	for _, player := range m.players {
		for _, shoot := range m.shoots[player] {
			shoots = append(shoots, shoot)
		}
	}

	return shoots
}

func (m *objManager) GetPlayerByID(id int32) (playerpkg.Player, bool) {
	//m.PlayersLock.Lock()
	//defer m.PlayersLock.Unlock()

	p, ok := m.players[id]
	if !ok {
		return nil, false
	}
	return p, true
}

func (m *objManager) Update() {
	shoots := m.GetShoots()
	now := time.Now()

	// Check Collision
	for _, shoot := range shoots {
		// Get current position of the shoot now
		curShootPosition := shoot.GetShootAtTime(now)

		// If shoot is out of boundary => remove
		if !m.isShootInWorld(curShootPosition) {
			m.RemoveShoot(shoot)
			continue
		}

		// Check if bullet hits players
		for _, player := range m.players {
			if shoot.GetPlayerID() != player.GetID() && m.isShootHitPlayer(curShootPosition, player) {
				log.Println("Shoot ", shoot.GetID(), " from player ", shoot.GetPlayerID(), " Hit Player", player.GetID())
				player.SetHealth(player.GetHealth() - 10)
				// Remove shoot
				m.RemoveShoot(shoot)
				// Remove player if health is under 0
				if player.GetHealth() <= 0 {
					log.Println("Push remove Player Event to event stream")
					// TODO: Removeplayer here, don't need send
					m.destroyPlayerStream <- common.DestroyPlayerEvent{
						PlayerID: player.GetID(),
						ClientID: -1,
					}
					// Add score for player who shoots
					if player, ok := m.GetPlayerByID(shoot.GetPlayerID()); ok {
						player.AddScore()
					}
				}
			}
		}

		// Check if bullet hits wall
		for _, block := range m.gameMap.GetRectBlocks() {
			// Consider checking by point instead of rectangle
			if m.isShootHitWall(curShootPosition, block) {
				// Remove shoot
				m.RemoveShoot(shoot)
				log.Println("Remove shoot", shoot)
				return
			}
		}
	}
}

// Consider generalize as shape package, just passing the shape
func (m *objManager) isShootInWorld(shoot shootpkg.Shoot) bool {
	return m.isInWorld(shoot.GetX(), shoot.GetY())
}

func (m *objManager) isInWorld(x, y float32) bool {
	return 0 <= x && x <= m.gameMap.GetWidth() && 0 <= y && y <= m.gameMap.GetHeight()
}

func (m *objManager) isShootHitPlayer(shoot shootpkg.Shoot, player playerpkg.Player) bool {
	return m.isCollidedPointCircle(shoot.GetPoint(), player.GetCircle())
}

func (m *objManager) isShootHitWall(shoot shootpkg.Shoot, wall shape.Rect) bool {
	return m.isCollidedPointRect(shoot.GetPoint(), wall)
}

func (m *objManager) isCollidedPointRect(point shape.Point, rect shape.Rect) bool {
	return rect.X1 <= point.X && point.X <= rect.X2 && rect.Y1 <= point.Y && point.Y <= rect.Y2
}

func (m *objManager) isCollidedPointCircle(point shape.Point, cir shape.Circle) bool {
	return m.dist(point.X, point.Y, cir.X, cir.Y) <= cir.Radius
}

func (m *objManager) isCollidedRectRect(rect1 shape.Rect, rect2 shape.Rect) bool {
	x1 := math.Max(float64(rect1.X1), float64(rect2.X1))
	x2 := math.Min(float64(rect1.X2), float64(rect2.X2))
	y1 := math.Max(float64(rect1.Y1), float64(rect2.Y1))
	y2 := math.Min(float64(rect1.Y2), float64(rect2.Y2))
	return (x1 < x2) && (y1 < y2)
}

// Checking rect circle collision givent the shape
func (m *objManager) isCollidedCircleRect(cir shape.Circle, rect shape.Rect) bool {
	if rect.X1 <= cir.X && cir.X <= rect.X2 && rect.Y1 <= cir.Y && cir.Y <= rect.Y2 {
		return true
	}
	var xnear = float32(math.Max(math.Min(float64(cir.X), float64(rect.X2)), float64(rect.X1)))
	var ynear = float32(math.Max(math.Min(float64(cir.Y), float64(rect.Y2)), float64(rect.Y1)))
	if m.dist(cir.X, cir.Y, xnear, ynear) <= cir.Radius {
		return true
	}
	return false
}

// dist calculate the distance between two points
func (m *objManager) dist(x1, y1, x2, y2 float32) float32 {
	return float32(math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2)))
}

func (m *objManager) isInside(box shape.Rect) bool {
	return box.X1 < 0 || box.X2 > gameconst.MapWidth || box.Y1 < 0 || box.Y2 > gameconst.MapHeight
}

func (m *objManager) IsValidPosition(circleObj shape.Circle) bool {
	for _, rectBlock := range m.gameMap.GetRectBlocks() {
		if m.isCollidedCircleRect(circleObj, rectBlock) == true {
			return false
		}
	}

	if circleObj.X-circleObj.Radius < 0 || circleObj.X+circleObj.Radius > m.gameMap.GetWidth() {
		return false
	}
	if circleObj.Y-circleObj.Radius < 0 || circleObj.Y+circleObj.Radius > m.gameMap.GetHeight() {
		return false
	}

	return true
}

func (m *objManager) MovePlayer(player playerpkg.Player, dx float32, dy float32, speed float32, timeElapsed float32) {
	pdx := dx * timeElapsed * speed
	pdy := dy * timeElapsed * speed
	if m.IsValidPosition(player.GetCircle().NextPosition(pdx, 0)) {
		player.Move(pdx, 0)
	}
	if m.IsValidPosition(player.GetCircle().NextPosition(0, pdy)) {
		player.Move(0, pdy)
	}
}

// RemovePlayerByClientID remove player with clientID. We need this because when a client disconnected, hub notified objManager with clientID
// We don't delete player here because this function is called from game master
//func (m *objManager) RemovePlayerByClientID(clientID int32) {
//// Send destroy event to game master
//m.eventStream <- common.DestroyPlayerEvent{
//ClientID: clientID,
//PlayerID: -1,
//}
//}

// RemovePlayer remove player according to player ID or ClientID, return playerID. This call need to be concurrently safe (by put behind event stream).
func (m *objManager) RemovePlayer(playerID int32, clientID int32) (removePlayerID int32) {
	// Remove by PlayerID
	if playerID != -1 {
		delete(m.players, playerID)
		return playerID
	}
	// Remove by ClientID
	for _, player := range m.players {
		if player.GetClientID() == clientID {
			delete(m.players, player.GetID())
			return player.GetID()
		}
	}
	return -1
}
