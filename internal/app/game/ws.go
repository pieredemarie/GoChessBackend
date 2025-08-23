package game

import (
	chesslogic "backendChess/pkg/chessLogic"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {return true},
}

type Client struct {
	Conn *websocket.Conn
	Game *chesslogic.Game
	Color chesslogic.Color	
}

