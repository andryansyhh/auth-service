package handler

import (
	"net/http"
	"strconv"

	"github.com/andryansyhh/auth-service/pkg/domain/dto"
	"github.com/andryansyhh/auth-service/pkg/middleware"
	"github.com/andryansyhh/auth-service/pkg/usecase"
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

	userRoutes := r.Group("/users", middleware.AuthMiddleware(h.jwt))
	{
		userRoutes.GET("", h.ListUsers)
		userRoutes.PUT("/:id", h.UpdateUser)
		userRoutes.DELETE("/:id", h.DeleteUser)
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.LoginRegisterRequest
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
	var req dto.LoginRegisterRequest
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

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.userUc.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.userUc.UpdateUser(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	if err := h.userUc.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
