package services

import (
	"context"
	"log"
	"time"

	"user-management-api/models"
	"user-management-api/repositories"
)

type LogService interface {
	LogAsync(userID uint, event string, data map[string]interface{})
	List(ctx context.Context, userID uint, page, perPage int) ([]models.UserLog, int64, error)
}

type logService struct {
	repo repositories.LogRepository
}

func NewLogService(repo repositories.LogRepository) LogService {
	return &logService{repo: repo}
}

func (s *logService) LogAsync(userID uint, event string, data map[string]interface{}) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userLog := &models.UserLog{
			UserID: userID,
			Event:  event,
			Data:   data,
		}

		if err := s.repo.Create(ctx, userLog); err != nil {
			log.Printf("async log failed: event=%s userID=%d err=%v", event, userID, err)
		}
	}()
}

func (s *logService) List(ctx context.Context, userID uint, page, perPage int) ([]models.UserLog, int64, error) {
	return s.repo.List(ctx, userID, page, perPage)
}
