package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID            primitive.ObjectID `json:"user_id" bson:"user_id"`
	Filename          string             `json:"filename" bson:"filename"`
	Version           int                `json:"version" bson:"version"`
	PreviousVersionID primitive.ObjectID `json:"previous_version_id,omitempty" bson:"previous_version_id,omitempty"`
	UploadDate        time.Time          `json:"upload_date" bson:"upload_date"`
}
