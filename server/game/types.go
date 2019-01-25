package game

type Game interface {
	ProcessInput(message []byte)
	NewPlayerConnect(clientID int32)
	GetQuitChannel() chan bool
	Update()
}
