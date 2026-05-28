package services

import (
	"testing"

	"user-management-api/utils"
)

func TestAIServiceFallback(t *testing.T) {
	cfg := &utils.Config{AIEnabled: false}
	svc := NewAIService(cfg)

	desc, err := svc.GenerateProductDescription("Headphones", "Electronics")
	if err != nil {
		t.Fatalf("GenerateProductDescription failed: %v", err)
	}
	if desc == "" {
		t.Error("expected non-empty description")
	}

	name, err := svc.SuggestCategoryName("wireless audio devices")
	if err != nil {
		t.Fatalf("SuggestCategoryName failed: %v", err)
	}
	if name == "" {
		t.Error("expected non-empty category name")
	}
}
