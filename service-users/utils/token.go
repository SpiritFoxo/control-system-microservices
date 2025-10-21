package utils

import (
	"errors"
	"fmt"
	"os"
	"service-users/internal/config"
	"service-users/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

func ValidateToken(c *gin.Context) error {
	token, err := GetToken(c)

	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}

	return errors.New("invalid token provided")
}

func GetToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("API_SECRET")), nil
	})
	return token, err
}

func getTokenFromRequest(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")

	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
