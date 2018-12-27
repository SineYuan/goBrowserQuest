package bqs

import (
	"log"

	"github.com/labstack/echo"
	"github.com/gorilla/websocket"
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
		World:     NewWrold(cfg),
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

var (
	upgrader = websocket.Upgrader{}
)

func (b *BQS) ToEchoHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()
		b.HandleConnection(ws)

		return nil
	}
}

func (b *BQS) Listen() {}
