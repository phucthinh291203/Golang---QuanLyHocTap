package Models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name           string             `bson:"name"`
	DateOfBirth    time.Time          `bson:"date_of_birth"`
	ClassID        primitive.ObjectID `bson:"class_id" json:"class_id"`
	Email          string             `bson:"email"`
	PhoneNumber    string             `bson:"phone_number"`
	Address        string             `bson:"address"`
	EnrollmentDate time.Time          `bson:"enrollment_date"`
	Gender         string             `bson:"gender"`
	Nationality    string             `bson:"nationality"`
	Avatar         string             `bson:"avatar"`
}
