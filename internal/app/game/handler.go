package game

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


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


func WebSocketHandler(c *gin.Context) {
	WSHandler(c.Writer, c.Request)
}
