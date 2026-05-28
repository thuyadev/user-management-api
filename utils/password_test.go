package utils

import (
	"testing"
)

func TestHashPasswordAndCheckPassword(t *testing.T) {
	password := "secret123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !CheckPassword(password, hash) {
		t.Error("CheckPassword should return true for correct password")
	}

	if CheckPassword("wrong", hash) {
		t.Error("CheckPassword should return false for wrong password")
	}
}
