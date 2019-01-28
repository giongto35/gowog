package game

type Game interface {
	ProcessInput(message []byte)
	NewPlayerConnect(clientID int32)
	RemovePlayer(playerID int32, clientID int32)
	GetQuitChannel() chan bool
	Update()
}
