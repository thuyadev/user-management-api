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
	List(ctx context.Context, userID uint, page, perPage int, search string) ([]models.UserLog, int64, error)
	EventStats(ctx context.Context, userID uint, days int) ([]models.LogEventStat, error)
	DailyStats(ctx context.Context, userID uint, days int) ([]models.LogDailyStat, error)
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

func (r *logRepository) List(ctx context.Context, userID uint, page, perPage int, search string) ([]models.UserLog, int64, error) {
	filter := bson.M{}
	if userID > 0 {
		filter["user_id"] = userID
	}
	if search != "" {
		filter["data.name"] = bson.M{"$regex": search, "$options": "i"}
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

func (r *logRepository) statsFilter(userID uint, days int) bson.M {
	since := time.Now().UTC().AddDate(0, 0, -(days - 1))
	since = time.Date(since.Year(), since.Month(), since.Day(), 0, 0, 0, 0, time.UTC)

	filter := bson.M{
		"created_at": bson.M{"$gte": since},
	}
	if userID > 0 {
		filter["user_id"] = userID
	}
	return filter
}

func (r *logRepository) EventStats(ctx context.Context, userID uint, days int) ([]models.LogEventStat, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: r.statsFilter(userID, days)}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$event",
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"count": -1}}},
		{{Key: "$project", Value: bson.M{
			"_id":   0,
			"event": "$_id",
			"count": 1,
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stats []models.LogEventStat
	if err := cursor.All(ctx, &stats); err != nil {
		return nil, err
	}
	if stats == nil {
		stats = []models.LogEventStat{}
	}
	return stats, nil
}

func (r *logRepository) DailyStats(ctx context.Context, userID uint, days int) ([]models.LogDailyStat, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: r.statsFilter(userID, days)}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"$dateToString": bson.M{
					"format": "%Y-%m-%d",
					"date":   "$created_at",
				},
			},
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"_id": 1}}},
		{{Key: "$project", Value: bson.M{
			"_id":   0,
			"date":  "$_id",
			"count": 1,
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stats []models.LogDailyStat
	if err := cursor.All(ctx, &stats); err != nil {
		return nil, err
	}
	return fillDailyStats(days, stats), nil
}

func fillDailyStats(days int, stats []models.LogDailyStat) []models.LogDailyStat {
	countByDate := make(map[string]int64, len(stats))
	for _, s := range stats {
		countByDate[s.Date] = s.Count
	}

	now := time.Now().UTC()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	start := end.AddDate(0, 0, -(days - 1))

	result := make([]models.LogDailyStat, 0, days)
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		result = append(result, models.LogDailyStat{
			Date:  dateStr,
			Count: countByDate[dateStr],
		})
	}
	return result
}
