package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/unbot2313/go-streaming-service/config"
	"github.com/unbot2313/go-streaming-service/internal/models"
)

type AuthServiceImp struct{}

type AuthService interface {

	GenerateToken(User models.User) (string, error)

	ValidateToken(token string) (*models.User, error)

}

func NewAuthService() AuthService {
	return &AuthServiceImp{}
}

func (service *AuthServiceImp) GenerateToken(user models.User) (string, error) {

	SecretToken := []byte(config.GetConfig().JWTSecretKey)

	//crear token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,               // Identificador único del usuario
		"username": user.Username,         // Nombre de usuario para referencia
		"email":    user.Email,            
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	//firmar token
	tokenString, err := token.SignedString(SecretToken)
	if err != nil {
		return "", fmt.Errorf("error al firmar el token: %v", err)
	}

	return tokenString, nil
}

func (service *AuthServiceImp) ValidateToken(tokenString string) (*models.User, error) {
	
	SecretToken := []byte(config.GetConfig().JWTSecretKey)

	// Parsear y verificar el token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validar el método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return SecretToken, nil
	})

	if err != nil {
		// Error al parsear o verificar el token
		return nil, fmt.Errorf("error al parsear el token: %v", err)
	}

	// Extraer y validar los claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		// Validar y construir el objeto usuario
		idFloat, ok := claims["user_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("user_id no es válido")
		}

		id := int(idFloat) // Convertir a entero

		username, ok := claims["username"].(string)
		if !ok {
			return nil, fmt.Errorf("username no es válido")
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, fmt.Errorf("email no es válido")
		}

		user := &models.User{
			ID:       id,
			Username: username,
			Email:    email,
		}

		return user, nil
	}

	// Si el token no es válido o los claims no son correctos
	return nil, fmt.Errorf("token inválido o claims inválidos")
}


