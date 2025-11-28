package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	controllerdto "vesuvio/internal/dto/controller"
	servicedto "vesuvio/internal/dto/service"
	"vesuvio/internal/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctl *AuthController) Register(c *gin.Context) {
	var req controllerdto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := ctl.authService.Register(c.Request.Context(), servicedto.RegisterUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case service.ErrEmailAlreadyExists, service.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		}
		return
	}

	c.JSON(http.StatusCreated, controllerdto.RegisterResponse{
		ID:      out.User.ID,
		Name:    out.User.Name,
		Email:   out.User.Email,
		IsAdmin: out.User.IsAdmin,
	})
}

func (ctl *AuthController) Login(c *gin.Context) {
	var req controllerdto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := ctl.authService.Login(c.Request.Context(), servicedto.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials, service.ErrInvalidInput:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		}
		return
	}

	c.JSON(http.StatusOK, controllerdto.LoginResponse{
		ID:      out.User.ID,
		Name:    out.User.Name,
		Email:   out.User.Email,
		IsAdmin: out.User.IsAdmin,
	})
}
