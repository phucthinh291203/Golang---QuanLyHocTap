package Models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HistoryUploaded struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	UploadedDate time.Time          `bson:"uploaded_date"`
	FileID       primitive.ObjectID `bson:"file_id" json:"file_id"`
}

type HistoryDownloaded struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	DownloadDate time.Time          `bson:"downloaded_date"`
	FileID       primitive.ObjectID `bson:"file_id" json:"file_id"`
}

type File struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	FileName string             `bson:"file_name"`
	FileType string             `bson:"file_type"`
}
