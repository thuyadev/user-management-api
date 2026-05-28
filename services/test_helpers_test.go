package services

import (
	"context"

	"user-management-api/models"
)

type mockLogService struct {
	events []string
}

func newMockLogService() *mockLogService {
	return &mockLogService{}
}

func (m *mockLogService) LogAsync(userID uint, event string, data map[string]interface{}) {
	m.events = append(m.events, event)
}

func (m *mockLogService) List(ctx context.Context, userID uint, page, perPage int, search string) ([]models.UserLog, int64, error) {
	return nil, 0, nil
}

func (m *mockLogService) EventStats(ctx context.Context, userID uint, days int) ([]models.LogEventStat, error) {
	return []models.LogEventStat{}, nil
}

func (m *mockLogService) DailyStats(ctx context.Context, userID uint, days int) ([]models.LogDailyStat, error) {
	return []models.LogDailyStat{}, nil
}
