package main

import (
	"backendChess/internal/app/auth"
	"backendChess/internal/app/game"
	"backendChess/pkg/jwtutils"
	"backendChess/pkg/storage"
	"fmt"
	"os"
    _ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
)

func main() {
    connStr := "user=postgres password=gigixach1234 dbname=GoChess sslmode=disable"
    db, err := storage.NewPostgresStorage(connStr)
    if err != nil {
        fmt.Println("db connection error: " + err.Error())
        return
    }

    game.InitMatchmaker(db)

    authHandler := &auth.Handler{Storage: db}
    gameHandler := &game.Handler{Storage: db}
    secret := os.Getenv("JWT_SECRET")

    r := gin.Default()

    authRoutes := r.Group("/auth")
    {
        authRoutes.POST("/register", authHandler.RegisterHandler)
        authRoutes.POST("/login", authHandler.LoginHandler)
    }

    protected := r.Group("/")
    protected.Use(jwtutils.JWTMiddleware(secret))
    {
        protected.GET("/users/me", authHandler.GetInfoUser)

        protected.POST("/game/find", gameHandler.FindGameHandler)
        protected.GET("/game/:id", gameHandler.GetGameStateHandler)
        protected.GET("/game/ws", gameHandler.WebSocketHandler)
    }

    r.Run(":8080")
}
