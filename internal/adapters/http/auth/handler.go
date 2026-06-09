package auth

import (
	"encoding/json"
	"errors"
	"log"
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
	/*
			 "email": email,
		      "password": hashedPassword,
		      "first_name": firstName,
		      "last_name": lastName,
		      "phone_number": phoneNumber,
		      "profile_pic": profilePic.toString(),
		      "city": city,
		      "country": country,
		      "date_of_birth": dateOfBirth.toIso8601String(),
	*/
	var req struct {
		FirstName   string `json:"first_name" `
		LastName    string `json:"last_name" `
		Email       string `json:"email" `
		Password    string `json:"password" `
		City        string `json:"city"`
		Country     string `json:"country"`
		PhoneNumber string `json:"phone_number" `
		ProfilePic  string `json:"profile_pic"`
	}

	body := map[string]any{}
	body_buf, err := json.MarshalIndent(body, "	", "		")
	if err != nil {
		c.JSON(500, gin.H{"error": "Error"})
		return
	}

	log.Default().Println(string(body_buf))
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
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
	var req struct {
		Email    string `form:"email" `
		Password string `form:"password" `
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := req.Email
	password := req.Password
	pair, err := h.uc.SignIn(c.Request.Context(), email, password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvitePending) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign in", "message": err.Error()})
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

func (h *Handler) SendVerificationEmail(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDVal.(uint)

	if err := h.uc.SendVerificationEmail(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	if err := h.uc.VerifyEmail(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (h *Handler) GetMe(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDVal.(uint)

	pair, err := h.uc.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Authorization", pair.AccessToken)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
	})
}
