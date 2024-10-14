package Models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subject struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	SubjectName string             `bson:"subject_name"`
	Credit      int                `bson:"credit"`
}


