package services

import (
	"context"
	"test-with-golang/Models"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SubjectDTO struct {
	collection *mongo.Collection
}

func NewSubjectDTO(collection *mongo.Collection) *SubjectDTO {
	return &SubjectDTO{
		collection: collection,
	}
}

func (service *SubjectDTO) CreateNewSubject(ctx *gin.Context) Models.Subject {
	var newData Models.Subject
	ctx.ShouldBindJSON(&newData)
	insertResult, err := service.collection.InsertOne(context.TODO(), newData)
	if err == nil {
		newData.ID = insertResult.InsertedID.(primitive.ObjectID)
		return newData
	}
	return Models.Subject{}
}

func (service *SubjectDTO) GetAllSubject() []Models.Subject {
	cursor, _ := service.collection.Find(context.TODO(), bson.D{})

	defer cursor.Close(context.TODO())

	var AllSubject []Models.Subject
	for cursor.Next(context.TODO()) {
		var eachSubject Models.Subject
		_ = cursor.Decode(&eachSubject)
		AllSubject = append(AllSubject, eachSubject)
	}

	return AllSubject
}

func (service *SubjectDTO) UpdateSubject(id primitive.ObjectID, ctx *gin.Context) Models.Subject {
	var updatedData Models.Subject
	ctx.ShouldBindJSON(&updatedData)

	err := service.collection.FindOneAndUpdate(context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": updatedData},
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedData)

	if err == nil {
		return updatedData
	}
	return Models.Subject{}
}

func (service *SubjectDTO) DeleteSubject(id primitive.ObjectID) error {
	_, err := service.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}
