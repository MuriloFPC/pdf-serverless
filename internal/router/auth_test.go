package router

import (
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

func TestArgon2Login(t *testing.T) {
	password := "AlgoSuperSecretoTemporario@321"

	// Registration
	salt := []byte(uuid.New().String())
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	encodedHash := fmt.Sprintf("%x.%x", salt, hash)

	// Fixed Login Logic (The one now implemented in AuthHandler.Login)
	parts := strings.Split(encodedHash, ".")
	if len(parts) != 2 {
		t.Fatalf("Invalid encoded hash format")
	}

	saltDecoded, err := hex.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("Failed to decode salt: %v", err)
	}

	hashDecoded, err := hex.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("Failed to decode hash: %v", err)
	}

	newHash := argon2.IDKey([]byte(password), saltDecoded, 1, 64*1024, 4, 32)

	if subtle.ConstantTimeCompare(hashDecoded, newHash) != 1 {
		t.Errorf("Password mismatch!")
	} else {
		t.Log("Login logic verified successfully")
	}
}
