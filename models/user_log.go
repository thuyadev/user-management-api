package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	LogEventUserCreated    = "user.created"
	LogEventUserUpdated    = "user.updated"
	LogEventUserDeleted    = "user.deleted"
	LogEventUserLogin      = "user.login"
	LogEventUserRegister   = "user.register"
	LogEventCategoryCreate = "category.created"
	LogEventCategoryUpdate = "category.updated"
	LogEventCategoryDelete = "category.deleted"
	LogEventProductCreate  = "product.created"
	LogEventProductUpdate  = "product.updated"
	LogEventProductDelete  = "product.deleted"
)

type UserLog struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    uint               `json:"user_id" bson:"user_id"`
	Event     string             `json:"event" bson:"event"`
	Data      map[string]interface{} `json:"data" bson:"data"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
