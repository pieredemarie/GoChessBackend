package game

import (
	chesslogic "backendChess/pkg/chessLogic"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MsgType string

const (
	MsgJoin     MsgType = "join"
	MsgMatch    MsgType = "match_found"
	MsgState    MsgType = "state"
	MsgMove     MsgType = "move"
	MsgUpdate   MsgType = "update"
	MsgError    MsgType = "error"
	MsgResign   MsgType = "resign"
	MsgGameOver MsgType = "game_over"
	MsgPing     MsgType = "ping"
	MsgPong     MsgType = "pong"
)

type Envelope struct {
	Type MsgType `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

type JoinData struct {
	Rating int `json:"rating"`
}

type MatchData struct { // it's better to move it to matchmaking.go
	RoomID string `json:"roomID"`
	Color chesslogic.Color `json:"color"`
}

type StateData struct { // state of the game
	BoardFEN string `json:"boardFEN,omitempty"` 
	Turn chesslogic.Color `json:"turn"`
	Status string `json:"status,omitempty"` // ok/check/mate/stalemate/
}

type MoveData struct { 
	From string `json:"from"`
	To string `json:"to"`
	Promotion string `json:"promotion,omitempty"` // queen,rook bishor or knight or empty if there's no promotion
}

type UpdateData struct {
	LastMove MoveData `json:"lastMove"`
	BoardFEN string `json:"boardFEN,omitempty"` 
	Turn chesslogic.Color `json:"turn"`
	Status string `json:"status,omitempty"` // // ok/check/mate/stalemate/
} 

type GameOverData struct {
	Reason string `json:"reason"`          // resign/checkmate/stalemate/disconnect
	Winner string `json:"winner,omitempty"` // "white"/"black"/"draw"
}

func mustRaw(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func readEnvelope(conn *websocket.Conn) (MsgType, json.RawMessage, error) {
	_, raw, err := conn.ReadMessage()
	if err != nil {
		return "", nil, err
	}
	var env Envelope
	if err := json.Unmarshal(raw,&env); err != nil {
		return "", nil, err
	}
	return env.Type, env.Data,nil
}

func errrorEnv(msg string) Envelope {
	return Envelope{Type: MsgError, Data: mustRaw(map[string]string{"message": msg})}
}

func abs(x int) int {
	if x < 0 {return -x}
	return x
}

func randRoomID() string {
	id := uuid.New()
	return id.String()
}

func winnerByColor(c chesslogic.Color) string {
	if c == chesslogic.White {return "white"}
	return "black"
}