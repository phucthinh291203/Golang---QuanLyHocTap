package services

import (
	"context"
	"log"
	"test-with-golang/Models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type classDto struct {
	collection *mongo.Collection
}

func NewClassDto(collection *mongo.Collection) ClassDTO {
	return &classDto{
		collection: collection,
	}
}

func (service *classDto) Save(class Models.Class, data Models.MyData) Models.Class {
	teacherFound := data.TeacherCollection.FindOne(context.TODO(), bson.M{"_id": class.TeacherID})
	if teacherFound.Err() != nil {
		log.Print("Teacher not found")
		return Models.Class{}
	}

	allStudentID := class.StudentID

	var students []Models.Student
	cursor, _ := data.StudentCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": allStudentID}})

	// Kiểm tra nếu cursor không có dữ liệu
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var eachStudent Models.Student
		_ = cursor.Decode(&eachStudent)
		students = append(students, eachStudent)
	}
	log.Print(len(students))
	log.Print(len(allStudentID))
	if len(students) != len(allStudentID) {
		log.Print("Hoc sinh khong có trong collection")
		return Models.Class{}
	}

	insertResult, err := service.collection.InsertOne(context.TODO(), class)
	if err == nil {
		class.ID = insertResult.InsertedID.(primitive.ObjectID)
		log.Print(class)
		return class
	}
	return Models.Class{}
}

func (service *classDto) FindAll() []Models.Class {
	var classes []Models.Class
	cursor, _ := service.collection.Find(context.TODO(), bson.D{})
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var class Models.Class
		_ = cursor.Decode(&class) // Bỏ qua lỗi
		classes = append(classes, class)
	}
	return classes
}

func (service *classDto) FindById(id string) Models.Class {
	var class Models.Class
	objectID, _ := primitive.ObjectIDFromHex(id) // Chuyển đổi id thành object id
	err := service.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&class)
	if err != nil {
		return Models.Class{} // Trả về lớp trống nếu không tìm thấy
	}
	return class
}

func (service *classDto) Update(id string, change Models.Class) Models.Class {
	var result = change
	ObjectID, _ := primitive.ObjectIDFromHex(id)
	err := service.collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": ObjectID},
		bson.M{"$set": result},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&result)

	if err != nil {
		return Models.Class{} // Trả về lớp trống nếu không tìm thấy
	}
	return result // Trả về lớp đã cập nhật
}

func (service *classDto) Delete(id string) bool {
	objectID, _ := primitive.ObjectIDFromHex(id)

	result, err := service.collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil || result.DeletedCount == 0 {
		return false // Trả về false nếu không xóa được hoặc không có tài liệu nào bị xóa
	}
	return true
}
