package game

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {return true},
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
		select {
		case <- ctx.Done():
			return
		case msg, ok := <-c.Send:
			if !ok {
				return
			}
			c.mu.Lock()
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			c.mu.Unlock()
			if err != nil {
				return
			}
			case <- ticker.C: {
				c.mu.Lock()
				c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				err := c.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"ping"}`))
				c.mu.Unlock()
				if err != nil {
					return
				}
			}
		}	
	}
}

func (c *Client) readPump(ctx context.Context) {
	defer func() {
		if c.room != nil {
			c.room.onDisconnect(c)
		} else {
			matchmaker.RemoveFromQueue(c)
		}
		_ = c.Conn.Close()
	}()

	c.Conn.SetReadLimit(1 << 16) // 64KB
	_ = c.Conn.SetReadDeadline(time.Now().Add(60 *time.Second))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	typ, data, err := readEnvelope(c.Conn)
	if err != nil {
		c.writeJSON(errrorEnv("failed to read first message: " + err.Error()))
		return
	}
	if typ != MsgJoin {
		c.writeJSON(errrorEnv("first message must be 'join'"))
		return
	}
	var j JoinData
	_ = json.Unmarshal(data, &j)
	if j.Rating == 0 {
		j.Rating = 900 // 900 is default rating for all 
	}
	c.Rating = j.Rating
	matchmaker.Enqueue(c)

	for {
		evType, evData, err := readEnvelope(c.Conn)
		if err != nil {
			return
		}
		switch evType {
		case MsgPing:
			c.writeJSON(Envelope{Type: MsgPong})
		case MsgMove:
			if c.room == nil {
				c.writeJSON(errrorEnv("not in a room"))
			}
			var m MoveData
			if err := json.Unmarshal(evData,&m); err != nil {
				c.writeJSON(errrorEnv("bad move payload"))
				continue
			}
			c.room.onMove(c,m)
		case MsgResign:
			if c.room != nil {
				c.room.onResign(c)
			}
		default:
			c.writeJSON(errrorEnv("unknow message: " + string(evType)))
		}
	}
}

var jwtKey = os.Getenv("JWT_SECRET")

func getUserIDFromToken(r *http.Request) (int, error) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return 0, fmt.Errorf("no authorization header")
    }

    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return 0, fmt.Errorf("invalid auth header format")
    }
    tokenStr := parts[1]

    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        return 0, fmt.Errorf("invalid token")
    }

    return claims.UserID, nil
}

func WSHandler(w http.ResponseWriter,r *http.Request) {
	userID, err := getUserIDFromToken(r)
    if err != nil {
        http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
        return
    }

	conn, err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		http.Error(w,"upgrade failed", http.StatusBadRequest)
		return
	}
	c := &Client{
		ID: randRoomID(),
		UserID: userID,
		Conn: conn,
		Send: make(chan []byte,64),
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() { c.writePump(ctx); cancel() }()
	c.readPump(ctx)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
