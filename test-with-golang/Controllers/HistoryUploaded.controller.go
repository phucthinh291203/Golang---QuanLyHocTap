package controllers

import (
	"context"
	"errors"
	"test-with-golang/Models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateHistoryUploaded(userCreated primitive.ObjectID, IdFile primitive.ObjectID, data Models.MyData) error {
	var newHitory Models.HistoryUploaded
	newHitory.UserID = userCreated
	newHitory.FileID = IdFile
	newHitory.UploadedDate = time.Now().Add(7 * time.Hour)
	_, err := data.HistoryUploadedCollection.InsertOne(context.TODO(), newHitory)
	if err != nil {
		return errors.New("thêm lịch sử uploaded thất bại")
	}
	return nil
}

func GetAllHistoryUploaded(data Models.MyData) []Models.HistoryUploaded {
	var allHistory []Models.HistoryUploaded
	cursor, _ := data.HistoryUploadedCollection.Find(context.TODO(), bson.D{})
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var data Models.HistoryUploaded
		_ = cursor.Decode(&data)
		allHistory = append(allHistory, data)
	}

	return allHistory
}

func GetHistoryUploadedById(historyId primitive.ObjectID, data Models.MyData) Models.HistoryUploaded {
	var history Models.HistoryUploaded
	data.HistoryUploadedCollection.FindOne(context.TODO(), bson.M{"_id": historyId}).Decode(&history)
	return history
}

func GetAllFile(data Models.MyData) []Models.File {
	var allFile []Models.File
	cursor, _ := data.FileCollection.Find(context.TODO(), bson.D{})
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var data Models.File
		_ = cursor.Decode(&data)
		allFile = append(allFile, data)
	}

	return allFile
}

func GetFileById(fileId primitive.ObjectID, data Models.MyData) Models.File {
	var file Models.File
	data.FileCollection.FindOne(context.TODO(), bson.M{"_id": fileId}).Decode(&file)
	return file
}

func CreateHistoryDownloaded(user_id primitive.ObjectID, file_id primitive.ObjectID, data Models.MyData) error {
	var historyDownloaded Models.HistoryDownloaded
	historyDownloaded.DownloadDate = time.Now().Add(7 * time.Hour)
	historyDownloaded.FileID = file_id
	historyDownloaded.UserID = user_id

	_, err := data.HistoryDownloadedCollection.InsertOne(context.TODO(), historyDownloaded)
	if err != nil {
		return errors.New("không thể tạo mới lịch sử tải về")
	}

	return nil
}

func GetHistoryDownloaded(data Models.MyData) []Models.HistoryDownloaded {
	var allDownloaded []Models.HistoryDownloaded
	cursor, _ := data.HistoryDownloadedCollection.Find(context.TODO(), bson.D{})
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var data Models.HistoryDownloaded
		_ = cursor.Decode(&data)
		allDownloaded = append(allDownloaded, data)
	}

	return allDownloaded
}
