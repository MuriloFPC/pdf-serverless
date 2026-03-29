package router

import (
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"pdf_serverless/internal/core/domain/entities"
	"pdf_serverless/internal/core/domain/interfaces"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type AuthHandler struct {
	userRepo     entities.UserRepository
	emailService interfaces.EmailService
	jwtSecret    string
}

func NewAuthHandler(userRepo entities.UserRepository, emailService interfaces.EmailService, secret string) *AuthHandler {
	return &AuthHandler{
		userRepo:     userRepo,
		emailService: emailService,
		jwtSecret:    secret,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		log.Printf("Register: Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if req.Email == "" || req.Password == "" {
		log.Printf("Register: Missing email or password")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}

	// Check if user already exists
	existingUser, err := h.userRepo.GetByEmail(c.Context(), req.Email)
	if err == nil && existingUser != nil {
		log.Printf("Register: User already exists with email: %s", req.Email)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User already exists"})
	}

	// Hash password with Argon2
	salt := []byte(uuid.New().String())
	hash := argon2.IDKey([]byte(req.Password), salt, 1, 64*1024, 4, 32)
	encodedHash := fmt.Sprintf("%x.%x", salt, hash)

	user := &entities.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: encodedHash,
		CreatedAt:    time.Now(),
	}

	if err := h.userRepo.Create(c.Context(), user); err != nil {
		log.Printf("Register: Error creating user in repository: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user, err := h.userRepo.GetByEmail(c.Context(), req.Email)
	if err != nil {
		log.Printf("Login: Error getting user by email (%s): %v", req.Email, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Verify Argon2 hash
	parts := strings.Split(user.PasswordHash, ".")
	if len(parts) != 2 {
		log.Printf("Login: Invalid password hash format for email: %s", req.Email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	salt, err := hex.DecodeString(parts[0])
	if err != nil {
		log.Printf("Login: Error decoding salt for email %s: %v", req.Email, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	hash, err := hex.DecodeString(parts[1])
	if err != nil {
		log.Printf("Login: Error decoding hash for email %s: %v", req.Email, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	newHash := argon2.IDKey([]byte(req.Password), salt, 1, 64*1024, 4, 32)
	if subtle.ConstantTimeCompare(hash, newHash) != 1 {
		log.Printf("Login: Password mismatch for email: %s", req.Email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	t, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		log.Printf("Login: Error signing JWT: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func (h *AuthHandler) JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Printf("JWTMiddleware: Missing authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"})
		}

		const bearerPrefix = "Bearer "
		if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			log.Printf("JWTMiddleware: Invalid authorization header format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid authorization header format"})
		}

		tokenString := authHeader[len(bearerPrefix):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("JWTMiddleware: Unexpected signing method: %v", token.Header["alg"])
				return nil, errors.New("unexpected signing method")
			}
			return []byte(h.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			log.Printf("JWTMiddleware: Invalid token: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("JWTMiddleware: Invalid token claims")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		c.Locals("user_id", claims["user_id"])
		return c.Next()
	}
}
