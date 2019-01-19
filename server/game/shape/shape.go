// package shape is to define shape interface of game object to handle object collision proper
package shape

// NextPosition returns a circle which is the predicted circle if player move in dx, dy direction
func (c Circle) NextPosition(dx float32, dy float32) Circle {
	return Circle{X: c.X + dx, Y: c.Y + dy, Radius: c.Radius}
}
