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
	// Search for the user in the database

	_, err := config.GetDB()

	if err != nil {
		return "", fmt.Errorf("error connecting to database: %v", err)
	}

	user, err := service.userService.GetUserByUserName(username)

	if err != nil {
		return "", fmt.Errorf("error when searching for user: %v", err)
	}

	// Verify password
	if !CheckPasswordHash(password, user.Password) {
		return "", fmt.Errorf("the password is invalid")
	}

	// Generate token
	token, err := service.GenerateToken(user)

	if err != nil {
		return "", fmt.Errorf("error generating token: %v", err)
	}

	return token, nil
}

func (service *AuthServiceImp) GenerateToken(user *models.User) (string, error) {

	SecretToken := []byte(config.GetConfig().JWTSecretKey)

	//create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.Id,               // User's unique identifier
		"username": user.Username,         // User name for reference
		"email":    user.Email,            
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
	})

	//sign token
	tokenString, err := token.SignedString(SecretToken)
	if err != nil {
		return "", fmt.Errorf("error when signing token: %v", err)
	}

	return tokenString, nil
}

func (service *AuthServiceImp) ValidateToken(tokenString string) (*models.User, error) {
	
	SecretToken := []byte(config.GetConfig().JWTSecretKey)

	// Parse and verify the token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signature method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signature method: %v", token.Header["alg"])
		}
		return SecretToken, nil
	})

	if err != nil {
		// Error parsing or verifying token
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// Extract and validate claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		// Validate and build the user object
		id, ok := claims["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("user_id is not valid")
		}

		username, ok := claims["username"].(string)
		if !ok {
			return nil, fmt.Errorf("username is not valid")
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, fmt.Errorf("email is not valid")
		}

		user := &models.User{
			Id:       id,
			Username: username,
			Email:    email,
		}

		return user, nil
	}

	// If the token is invalid or the claims are not correct
	return nil, fmt.Errorf("invalid token or invalid claims")
}

// Function to hash a password
func HashPassword(password string) (string, error) {
	// Generate the password hash with a default cost (bcrypt.DefaultCost)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error when generating the hash: %v", err)
	}
	return string(hashedPassword), nil
}

// Function to compare an unhashed password with its hash
func CheckPasswordHash(password, hashedPassword string) bool {
	// Compare the password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil // Returns true if there were no errors
}


