package testutil

import (
	"context"

	"user-management-api/models"
)

type MockLogService struct {
	Events []string
}

func NewMockLogService() *MockLogService {
	return &MockLogService{}
}

func (m *MockLogService) LogAsync(userID uint, event string, data map[string]interface{}) {
	m.Events = append(m.Events, event)
}

func (m *MockLogService) List(ctx context.Context, userID uint, page, perPage int) ([]models.UserLog, int64, error) {
	return nil, 0, nil
}
