package repositories

import (
	"context"
	"time"

	"user-management-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogRepository interface {
	Create(ctx context.Context, log *models.UserLog) error
	List(ctx context.Context, userID uint, page, perPage int) ([]models.UserLog, int64, error)
}

type logRepository struct {
	collection *mongo.Collection
}

func NewLogRepository(collection *mongo.Collection) LogRepository {
	return &logRepository{collection: collection}
}

func (r *logRepository) Create(ctx context.Context, log *models.UserLog) error {
	now := time.Now()
	log.CreatedAt = now
	log.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, log)
	return err
}

func (r *logRepository) List(ctx context.Context, userID uint, page, perPage int) ([]models.UserLog, int64, error) {
	filter := bson.M{}
	if userID > 0 {
		filter["user_id"] = userID
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * perPage)).
		SetLimit(int64(perPage))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var logs []models.UserLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
