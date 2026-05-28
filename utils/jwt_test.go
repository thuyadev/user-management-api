package utils

import (
	"testing"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret"
	token, expiresAt, err := GenerateToken(1, "test@example.com", "admin", secret, 24)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if expiresAt.IsZero() {
		t.Error("expected expires_at")
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	if claims.UserID != 1 {
		t.Errorf("expected user_id 1, got %d", claims.UserID)
	}
	if claims.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", claims.Email)
	}
	if claims.Role != "admin" {
		t.Errorf("expected role admin, got %s", claims.Role)
	}
	if claims.ExpiresAt == nil {
		t.Error("expected JWT exp claim")
	}
}

func TestParseTokenInvalid(t *testing.T) {
	_, err := ParseToken("invalid.token.here", "secret")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}
