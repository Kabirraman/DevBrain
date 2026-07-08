package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/models"
)

func RegisterHandler(c *gin.Context) {

	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	user, err := RegisterUser(req)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	token, err := GenerateToken(
		user.ID.String(),
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		Token: token,
	})
}

func LoginHandler(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	token, err := LoginUser(req)

	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})

		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
	})
}
func MeHandler(c *gin.Context) {

	userID := c.MustGet("userID")

	var user models.User

	err := database.DB.
		First(&user, "id = ?", userID).
		Error

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}