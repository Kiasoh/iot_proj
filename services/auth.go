package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"iot_proj/models"
)

var jwtKey []byte

func init() {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		key = "super-secret-key"
	}
	jwtKey = []byte(key)
}

func (s *Service) Register(creds *models.Credentials) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:    creds.Email,
		Password: string(hashedPassword),
	}

	return s.Repo.CreateUser(user)
}

func (s *Service) Login(creds *models.Credentials) (string, error) {
	user, err := s.Repo.GetUserByEmail(creds.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func (s *Service) GetUserProfile(id int) (*models.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *Service) ValidateJWT(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid signature")
		}
		return nil, errors.New("bad token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
