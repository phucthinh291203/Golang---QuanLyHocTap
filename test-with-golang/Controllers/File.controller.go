package controllers

import (
	"context"
	"errors"
	"test-with-golang/Models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateFile(fileName string, data Models.MyData, fileType string) (primitive.ObjectID, error) {
	var newFile Models.File
	newFile.FileName = fileName
	newFile.FileType = fileType
	insertResult, err := data.FileCollection.InsertOne(context.TODO(), newFile)
	if err == nil {
		newFile.ID = insertResult.InsertedID.(primitive.ObjectID)
		return newFile.ID, nil
	}
	return primitive.NilObjectID, errors.New("không thể tạo mới file")
}
