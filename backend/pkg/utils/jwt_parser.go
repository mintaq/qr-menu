package utils

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// TokenMetadata struct to describe metadata in JWT.
type TokenMetadata struct {
	UserID      int
	Credentials map[string]bool
	Expires     int64
}

type RenewToken struct {
	RefreshToken string `json:"refresh_token"`
}

// ExtractTokenMetadata func to extract metadata from JWT.
func ExtractTokenMetadata(c *fiber.Ctx) (*TokenMetadata, error) {
	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}

	// Create a new renew refresh token struct.
	renew := &RenewToken{}

	// Checking received data from JSON body.
	if err := c.BodyParser(renew); err != nil {
		// Return, if JSON data is not correct.
		return nil, err
	}

	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// User ID.
		userID, err := strconv.Atoi(claims["id"].(string))
		if err != nil {
			return nil, err
		}

		now := time.Now().Unix()

		if renew.RefreshToken != "" {
			expiresRefreshToken, err := ParseRefreshToken(renew.RefreshToken)
			if err != nil {
				return nil, err
			}
			if now >= expiresRefreshToken {
				return nil, errors.New("refresh token is expired")
			}
		} else if claims.VerifyExpiresAt(now, true) {
			return nil, errors.New("token is expired")
		}

		// Expires time.
		expires := int64(claims["expires"].(float64))
		// User credentials.
		credentials := map[string]bool{
			"book:create": claims["book:create"].(bool),
			"book:update": claims["book:update"].(bool),
			"book:delete": claims["book:delete"].(bool),
		}

		return &TokenMetadata{
			UserID:      userID,
			Credentials: credentials,
			Expires:     expires,
		}, nil
	}

	return nil, err
}

func extractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")

	// Normally Authorization HTTP header.
	onlyToken := strings.Split(bearToken, " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := extractToken(c)

	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
