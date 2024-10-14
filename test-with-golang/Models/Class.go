package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Class struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	ClassName       string               `bson:"class_name"`
	TeacherID       primitive.ObjectID   `bson:"teacher_id"`
	StudentID       []primitive.ObjectID `bson:"student_ids"`
	SchoolYearStart int                  `bson:"school_year_start"`
	SchoolYearEnd   int                  `bson:"school_year_end"`
}
