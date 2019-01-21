package mappkg

import (
	"strings"

	"github.com/giongto35/gowog/server/Message_proto"
	"github.com/giongto35/gowog/server/game/config"
	"github.com/giongto35/gowog/server/game/shape"
)

// Map doesn't contain proto because proto cannot store 2d
// TODO: remove proto from other type
type mapImpl struct {
	blocks      [][]int
	rectBlocks  []shape.Rect
	numRows     int
	numCols     int
	blockWidth  float32
	blockHeight float32
	Width       float32
	Height      float32
}

// NewMap creates a new map
func NewMap(blockWidth float32, blockHeight float32) Map {
	// consider getting blockWidth blockHeight from file map instead of config
	gameMap := &mapImpl{}
	gameMap.blockWidth = blockWidth
	gameMap.blockHeight = blockHeight

	// Load map from a name, the map is a grid of wall
	gameMap.loadMap()
	return gameMap
}

func (m *mapImpl) loadMap() {
	strmap := config.RawMap
	sm := strings.Split(strmap, "\n")
	m.numRows = len(sm) - 1 // Note -1 because omit the last one
	m.numCols = len(sm[0])
	m.blocks = make([][]int, m.numRows)
	m.rectBlocks = make([]shape.Rect, 0)

	// Create map
	for i := 0; i < m.numRows; i++ {
		m.blocks[i] = make([]int, m.numCols)
		for j := 0; j < m.numCols; j++ {
			m.blocks[i][j] = int(sm[i][j] - '0')
			if m.blocks[i][j] != 0 {
				m.rectBlocks = append(m.rectBlocks, shape.Rect{
					Y1: float32(i) * m.blockHeight,
					Y2: float32(i+1) * m.blockHeight,
					X1: float32(j) * m.blockWidth,
					X2: float32(j+1) * m.blockWidth,
				})
			}
		}
	}

}

func (m *mapImpl) ToProto() *Message_proto.Map {
	proto := &Message_proto.Map{
		NumCols:     int32(m.numCols),
		NumRows:     int32(m.numRows),
		BlockWidth:  m.blockWidth,
		BlockHeight: m.blockHeight,
	}
	for i := 0; i < m.numRows; i++ {
		for j := 0; j < m.numCols; j++ {
			proto.Block = append(proto.Block, int32(m.blocks[i][j]))
		}
	}
	return proto
}

func (m *mapImpl) IsCollide(x float32, y float32) bool {
	bx := int(x / m.blockWidth)
	by := int(y / m.blockHeight)
	return (m.blocks[by][bx] != 0)
}

func (m *mapImpl) GetNumCols() int {
	return m.numCols
}

func (m *mapImpl) GetNumRows() int {
	return m.numRows
}

func (m *mapImpl) GetWidth() float32 {
	return float32(m.numCols) * m.blockWidth
}

func (m *mapImpl) GetHeight() float32 {
	return float32(m.numRows) * m.blockHeight
}

func (m *mapImpl) GetRectBlocks() []shape.Rect {
	return m.rectBlocks
}
