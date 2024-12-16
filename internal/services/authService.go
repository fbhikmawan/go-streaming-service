package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/unbot2313/go-streaming-service/config"
	"github.com/unbot2313/go-streaming-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImp struct{
	userService UserService
}

type AuthService interface {

	GenerateToken(User *models.User) (string, error)
	ValidateToken(token string) (*models.User, error)
	Login(username, password string) (string, error)

}

func NewAuthService() AuthService {
	return &AuthServiceImp{
		userService: NewUserService(),
	}
}

func (service *AuthServiceImp) Login(username, password string) (string, error) {
	// Buscar el usuario en la base de datos

	_, err := config.GetDB()

	if err != nil {
		return "", fmt.Errorf("error al conectar a la base de datos: %v", err)
	}

	user, err := service.userService.GetUserByUserName(username)

	if err != nil {
		return "", fmt.Errorf("error al buscar el usuario: %v", err)
	}

	// Verificar la contraseña
	if !CheckPasswordHash(password, user.Password) {
		return "", fmt.Errorf("la contraseña no es válida")
	}

	// Generar el token
	token, err := service.GenerateToken(user)

	if err != nil {
		return "", fmt.Errorf("error al generar el token: %v", err)
	}

	return token, nil
}

func (service *AuthServiceImp) GenerateToken(user *models.User) (string, error) {

	SecretToken := []byte(config.GetConfig().JWTSecretKey)

	//crear token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.Id,               // Identificador único del usuario
		"username": user.Username,         // Nombre de usuario para referencia
		"email":    user.Email,            
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Expira en 24 horas
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
		id, ok := claims["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("user_id no es válido")
		}

		username, ok := claims["username"].(string)
		if !ok {
			return nil, fmt.Errorf("username no es válido")
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, fmt.Errorf("email no es válido")
		}

		user := &models.User{
			Id:       id,
			Username: username,
			Email:    email,
		}

		return user, nil
	}

	// Si el token no es válido o los claims no son correctos
	return nil, fmt.Errorf("token inválido o claims inválidos")
}

// Función para hashear una contraseña
func HashPassword(password string) (string, error) {
	// Generar el hash de la contraseña con un costo predeterminado (bcrypt.DefaultCost)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error al generar el hash: %v", err)
	}
	return string(hashedPassword), nil
}

// Función para comparar una contraseña sin hashear con su hash
func CheckPasswordHash(password, hashedPassword string) bool {
	// Comparar la contraseña con el hash
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil // Devuelve true si no hubo errores
}


