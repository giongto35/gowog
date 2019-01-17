package ws

type Client interface {
	WritePump()
	ReadPump()
	Send(message []byte)
	GetSend() chan []byte
	GetID() int32
	Close()
}

type Hub interface {
	Run()
	UnRegister(c Client)
	Register(c Client) <-chan int32
	ReceiveMessage(b []byte)
	Broadcast(b []byte)
	BroadcastExclude(b []byte, id int32)
	BindGameMaster(g IGame)
	Send(clientID int32, b []byte)
	//Close(clientID int32)
}

type IGame interface {
	ProcessInput(message []byte)
	RemovePlayerByClientID(clientID int32)
}
