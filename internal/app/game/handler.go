package game

import (
	"backendChess/pkg/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Storage *storage.PostgresStorage
}

func FindGameHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "connect via WebSocket at /game/ws",
	})
}

func GetGameStateHandler(c *gin.Context) {
	id := c.Param("id")
	val, ok := rooms.Load(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}
	r := val.(*Room)

	c.JSON(http.StatusOK, gin.H{
		"id":     r.ID,
		"white":  r.White.ID,
		"black":  r.Black.ID,
		"turn":   r.Game.Board.Turn,
		"status": "active", 
	})
}

func (h *Handler) SaveMove(c *gin.Context) {
	var req storage.DBMove
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": "bad request"})
		return
	}

	err := h.Storage.SaveMove(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK,gin.H{"responce": "move saved"})
}

func (h *Handler) SaveGame(c *gin.Context) {
	var req storage.DBGame
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": "bad request"})
		return
	}

	err := h.Storage.SaveGame(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK,gin.H{"responce": "game saved"})
}

func (h *Handler) UpdateGameResult(c *gin.Context) {
	var req UpdateGameRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": "bad request"})
		return
	}

	err := h.Storage.UpdateGameResult(req.GameID,req.Winner,req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK,gin.H{"responce": "game updated"})
}

func WebSocketHandler(c *gin.Context) {
	WSHandler(c.Writer, c.Request)
}
