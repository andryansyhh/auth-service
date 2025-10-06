package handler

import (
	"net/http"

	"github.com/andryansyhh/auth-service/internal/domain/dto"
	"github.com/andryansyhh/auth-service/internal/middleware"
	"github.com/andryansyhh/auth-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUc usecase.UserUsecase
	jwt    *middleware.JWTManager
}

func NewUserHandler(userUc usecase.UserUsecase, jwt *middleware.JWTManager) *UserHandler {
	return &UserHandler{userUc: userUc, jwt: jwt}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	auth := r.Group("/account", middleware.AuthMiddleware(h.jwt))
	auth.GET("/profile", h.Profile)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.userUc.Register(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	resp, err := h.userUc.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Profile(c *gin.Context) {
	username := c.GetString("username")

	resp, err := h.userUc.GetProfile(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
