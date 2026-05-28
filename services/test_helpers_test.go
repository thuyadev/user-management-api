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

func (m *mockLogService) List(ctx context.Context, userID uint, page, perPage int) ([]models.UserLog, int64, error) {
	return nil, 0, nil
}
