package auth

import (
	"backendChess/pkg/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterHandler(c *gin.Context) {
	var newUser RegisterRequest 

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest,ErrorResponce{
			Message: err.Error(),
		})
		return
	}

	user := storage.User{
		Username: newUser.Username,
		Email: newUser.Email,
		Password: newUser.Password,
	}
	if err := h.Storage.Register(user); err != nil {
		c.JSON(http.StatusInternalServerError,ErrorResponce{
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var req LoginRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponce{
			Message: err.Error(),
		})
	}
	
	token, err := h.Storage.Login(req.Email,req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponce{
			Message: err.Error(),
		})
	}
	c.JSON(http.StatusOK,LoginResponce{
		Token: token,
	})
}

func (h *Handler) GetInfoUser(c *gin.Context) {
	userId, _ := c.Get("userID")

	user, err := h.Storage.GetUserInfo(userId.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError,ErrorResponce{
			Message: err.Error(),
		})
	}

	c.JSON(http.StatusOK,gin.H{
		"username": user.Username,
		"rating": user.Elo,
		"gamesCount": user.GamesCount,
	})
}
