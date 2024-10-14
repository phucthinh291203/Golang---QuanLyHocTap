package Models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Teacher struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Email       string               `json:"email" bson:"email"`
	PhoneNumber string               `json:"phone_number" bson:"phone_number"`
	Username    string               `json:"username" bson:"username"`
	Password    string               `json:"password" bson:"password"`
	TeacherName string               `bson:"teacher_name"`
	ClassID     []primitive.ObjectID `bson:"class_ids"`
	SubjectID   []primitive.ObjectID `bson:"subject_ids"`
}
