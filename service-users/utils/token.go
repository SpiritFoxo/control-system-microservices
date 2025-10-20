package utils

import (
	"fmt"
	"service-users/internal/config"
	"service-users/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user models.User, cfg *config.Config) (string, error) {

	tokenLifespan, err := strconv.Atoi(strings.TrimSpace(cfg.TokenMinuteLifespan))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = user.ID
	claims["roles"] = user.Roles
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(cfg.TokenSecret))

}

func GenerateRefreshToken(user models.User, cfg *config.Config) (string, error) {
	tokenLifespan, err := strconv.Atoi(strings.TrimSpace(cfg.RefreshTokenHourLifespan))
	if err != nil {
		tokenLifespan = 72
		fmt.Println("Warning: REFRESH_TOKEN_HOUR_LIFESPAN not set or invalid. Using default 72 hours.")
	}

	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["roles"] = user.Roles
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.RefreshTokenSecret))
}
