package miscallenous

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// generateJWTToken creates a new JWT token for the given user.
func GenerateJWTToken(object any,tokenName string,ID uint) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		tokenName: object,
			"sub": ID,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(), // Expiration time (30 days from now)
	}).SignedString([]byte(os.Getenv("SECRET")))
}

func DecodeJWTToken(tokenString string) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("SECRET")), nil
	}

	return jwt.Parse(tokenString, keyFunc)
}

// hashPassword hashes the given password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}


func VerifyPassword(hashedPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}