package game

import (
	chesslogic "backendChess/pkg/chessLogic"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {return true},
}

type Client struct {
	ID     string
	Conn   *websocket.Conn
	Send   chan []byte
	Rating int
	Color  chesslogic.Color

	room   *Room
	mu     sync.Mutex 
}

func (c *Client) writeJSON(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Println("writeJSON marshal: ", err)
	}
	select {
	case c.Send <- b:
	default:
		_ = c.Conn.Close()
	}
}

func (c *Client) writePump(ctx context.Context) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	defer func() {
		_ = c.Conn.Close()
	}()
	for {
			
	}
}

func WSHandler(w http.ResponseWriter,r *http.Request) {
	conn, err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		http.Error(w,"upgrade failed", http.StatusBadRequest)
		return
	}
	c := &Client{
		ID: randRoomID(),
		Conn: conn,
		Send: make(chan []byte,64),
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() { c.writePump(ctx); cancel() }()
	//c.readPump(ctx)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
