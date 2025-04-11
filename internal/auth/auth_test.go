package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWTAndValidateJWT(t *testing.T) {
	userID := uuid.New()	
	tokenSecret := "mysecretekey"
	expiresIn := 2 * time.Minute

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("failed to created JWT: %v", err)
	}

	parsedID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}
	if parsedID != userID {
		t.Errorf("expected userID %v, got %v", userID, parsedID)
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecretkey"
	expiresIn := -1 * time.Minute

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestInvalidSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecretkey"
	wrongSecret := "wrongsecret"
	expiresIn := 2 * time.Minute

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("expected error for invalid secret, got nil")
	}
}
