package bqs

import (
	"log"
	"github.com/kataras/iris"
	"github.com/iris-contrib/websocket"
)

type BQS struct {
	World       *World
	PlayerMap   map[string]*Player
	ZoneMap     map[string]*Zone
	PlayerCount int
	Logger      *log.Logger
}

type Packet struct {
	Type int
	Data []byte
	Err  error
}

func recvPacket(c *websocket.Conn) *Packet {
	msg := new(Packet)
	msg.Type, msg.Data, msg.Err = c.ReadMessage()
	return msg
}

func NewBQS(cfg *Config) *BQS {
	bqs := &BQS{
		World: NewWrold(cfg),
		PlayerMap: make(map[string]*Player),
	}
	return bqs
}

func (b *BQS) HandleConnection(c *websocket.Conn) {
	defer c.Close()
	err := c.WriteMessage(websocket.TextMessage, []byte("go"))
	if err != nil {
		b.Logger.Println(err)
		return
	}
	player := NewPlayer(c, b.World)
	go player.HandleLoop()
	for {
		select {
		case <-player.exit:
			return
		case player.PacketChan <- recvPacket(c):
		}
	}
}

func (b *BQS) boardcast(e *Event) {
	for _, z := range b.ZoneMap {
		z.EventCh <- e
	}
}

func (b *BQS) ToIrisHandler() iris.HandlerFunc {
	return func(ctx *iris.Context) {
		var connService = websocket.New(b.HandleConnection) // use default options
		connService.Upgrade(ctx)
	}
}

func (b *BQS) Listen() {}
