package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pdf_serverless/internal/infra/database"
	"pdf_serverless/internal/infra/email"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegister_DuplicateEmail(t *testing.T) {
	app := fiber.New()
	userRepo := database.NewUserMemoryRepository()
	emailService := email.NewNoOpEmailService()
	authHandler := NewAuthHandler(userRepo, emailService, "secret")

	app.Post("/register", authHandler.Register)

	reqBody, _ := json.Marshal(RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	})

	// First registration
	req1 := httptest.NewRequest("POST", "/register", bytes.NewReader(reqBody))
	req1.Header.Set("Content-Type", "application/json")
	resp1, _ := app.Test(req1)

	if resp1.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", resp1.StatusCode)
	}

	// Second registration with same email
	req2 := httptest.NewRequest("POST", "/register", bytes.NewReader(reqBody))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2)

	if resp2.StatusCode != http.StatusConflict {
		t.Errorf("Expected status 409 Conflict, got %d", resp2.StatusCode)
	}

	var errorResp map[string]string
	json.NewDecoder(resp2.Body).Decode(&errorResp)
	if errorResp["error"] != "User already exists" {
		t.Errorf("Expected error message 'User already exists', got '%s'", errorResp["error"])
	}
}
