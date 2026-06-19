package auth

import (
	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(req RegisterRequest) (*models.User, error) {

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	err = database.DB.Create(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func LoginUser(req LoginRequest) (string, error) {

	var user models.User

	err := database.DB.
		Where("email = ?", req.Email).
		First(&user).Error

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return "", err
	}

	token, err := GenerateToken(
		user.ID.String(),
	)

	if err != nil {
		return "", err
	}

	return token, nil
}