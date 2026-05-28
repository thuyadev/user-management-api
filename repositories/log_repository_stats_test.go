package repositories

import (
	"testing"
	"time"

	"user-management-api/models"
)

func TestFillDailyStatsIncludesQuietDays(t *testing.T) {
	now := time.Now().UTC()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	start := end.AddDate(0, 0, -2)

	stats := []models.LogDailyStat{
		{Date: start.Format("2006-01-02"), Count: 3},
		{Date: end.Format("2006-01-02"), Count: 5},
	}

	result := fillDailyStats(3, stats)
	if len(result) != 3 {
		t.Fatalf("expected 3 days, got %d", len(result))
	}
	if result[0].Count != 3 {
		t.Errorf("expected first day count 3, got %d", result[0].Count)
	}
	if result[1].Count != 0 {
		t.Errorf("expected quiet day count 0, got %d", result[1].Count)
	}
	if result[2].Count != 5 {
		t.Errorf("expected last day count 5, got %d", result[2].Count)
	}
}
