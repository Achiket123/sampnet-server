package auth

import (
	"errors"
	"net/http"
	domain "server/internal/domain/auth"
	usecase "server/internal/usecase/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc domain.UseCase
}

func NewHandler(uc domain.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) SignUp(c *gin.Context) {
	var req struct {
		FirstName   string `form:"first_name" binding:"required"`
		LastName    string `form:"last_name" binding:"required"`
		Email       string `form:"email" binding:"required,email"`
		Password    string `form:"password" binding:"required"`
		City        string `form:"city"`
		Country     string `form:"country"`
		PhoneNumber string `form:"phone_number" binding:"required"`
		ProfilePic  string `form:"profile_pic"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		City:        req.City,
		Country:     req.Country,
		PhoneNumber: req.PhoneNumber,
		ProfilePic:  req.ProfilePic,
	}

	pair, err := h.uc.SignUp(c.Request.Context(), user, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign up", "message": err.Error()})
		return
	}

	c.Header("Authorization", pair.AccessToken)
	c.JSON(http.StatusOK, gin.H{
		"message":       "User created successfully",
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
	})
}

func (h *Handler) SignIn(c *gin.Context) {
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")

	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	pair, err := h.uc.SignIn(c.Request.Context(), email, password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign in"})
		return
	}

	c.Header("Authorization", pair.AccessToken)
	c.JSON(http.StatusOK, gin.H{
		"message":       "Sign in successful",
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
	})
}

func (h *Handler) CompleteSignIn(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		PhoneNumber string `json:"phone_number" binding:"required"`
		Password    string `json:"password" binding:"required"`
		City        string `json:"city" binding:"required"`
		Country     string `json:"country" binding:"required"`
		ProfilePic  string `json:"profile_pic" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pair, err := h.uc.CompleteSignIn(c.Request.Context(), req.Email, req.PhoneNumber, req.Password, req.City, req.Country, req.ProfilePic)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete sign in"})
		return
	}

	c.Header("Authorization", pair.AccessToken)
	c.JSON(http.StatusOK, gin.H{
		"message":       "User updated successfully",
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
	})
}

func (h *Handler) ValidateEmployee(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID := userIDVal.(uint)

	token, err := h.uc.ValidateEmployee(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotAnEmployee) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate employee"})
		return
	}

	c.Header("Authorization", token)
	c.JSON(http.StatusOK, gin.H{"message": "Employee validated successfully"})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pair, err := h.uc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}

	c.Header("Authorization", pair.AccessToken)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
