package Models

import "go.mongodb.org/mongo-driver/mongo"

type MyData struct {
	ClassCollection   *mongo.Collection
	StudentCollection *mongo.Collection
	UserCollection    *mongo.Collection
	ScoreCollection   *mongo.Collection
	SubjectCollection *mongo.Collection
	TeacherCollection *mongo.Collection
	BangDiemCollection *mongo.Collection
	OTPCollection *mongo.Collection
	SMSCollection	*mongo.Collection
}
