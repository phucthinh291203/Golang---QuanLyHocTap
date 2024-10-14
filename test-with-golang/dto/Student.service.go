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

type studentDTO struct {
	collection *mongo.Collection
}

func (service *studentDTO) IndexNationality() {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{
			Key:   "nationality",
			Value: 1,
		}}, // Index theo trường nationality
		Options: options.Index().SetName("_indexNationality"),
	}

	service.collection.Indexes().CreateOne(context.TODO(), indexModel)
}

func NewStudent(collection *mongo.Collection) StudentDTO {

	dto := &studentDTO{
		collection: collection,
	}
	dto.IndexNationality()
	return dto
}

func (service *studentDTO) Create(newData Models.Student,data Models.MyData) error {
	//Tìm ra giáo viên hiện tại
	var class Models.Class
	err := service.collection.FindOne(context.TODO(), bson.M{"_id": newData.ClassID}).Decode(&class)
	if err != nil {
		log.Print("Khong tim thay lớp trong collection")
		return err
	}

	_, err = service.collection.InsertOne(context.TODO(), data)
	return err
}

func (service *studentDTO) FindAll() []Models.Student {
	var allData []Models.Student
	cursor, _ := service.collection.Find(context.TODO(), bson.D{})
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var data Models.Student
		_ = cursor.Decode(&data)
		allData = append(allData, data)
	}

	return allData

}

func (service *studentDTO) FindById(id string, data Models.MyData) (Models.Student, string) {
	var student Models.Student
	var class Models.Class
	classCollection := data.ClassCollection
	objectID, _ := primitive.ObjectIDFromHex(id)

	err := service.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&student)
	if err != nil {
		return Models.Student{}, ""
	}

	err = classCollection.FindOne(context.TODO(), bson.M{"_id": student.ClassID}).Decode(&class)
	if err != nil {
		return Models.Student{}, ""
	}

	// pipeline := mongo.Pipeline{
	// 	{{Key: "$match", Value: bson.M{"_id": objectID}}},
	// 	{{Key: "$lookup", Value: bson.M{
	// 		"from":         "classCollection",
	// 		"localField":   "class_id",
	// 		"foreignField": "_id",
	// 		"as":           "classInfo",
	// 	}}},
	// 	{{Key: "unwind", Value: bson.M{"path": "$classInfo", "preserveNullAndEmptyArrays": true}}},
	// }

	// cursor, _ := service.collection.Aggregate(context.TODO(), pipeline)
	// defer cursor.Close(context.TODO())

	// if cursor.Next(context.TODO()) {
	// 	if err := cursor.Decode(&result); err != nil {
	// 		return Models.Student{}
	// 	}
	// }

	return student, class.ClassName
}

func (service *studentDTO) Update(id string, change Models.Student) Models.Student {

	var result = change
	ObjectID, _ := primitive.ObjectIDFromHex(id)
	err := service.collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": ObjectID},
		bson.M{"$set": result},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&result)

	log.Print(ObjectID)
	log.Print(result)
	if err != nil {
		return Models.Student{} // Trả về lớp trống nếu không tìm thấy
	}
	return result // Trả về lớp đã cập nhật
}

func (service *studentDTO) Delete(id string) bool {
	objectID, _ := primitive.ObjectIDFromHex(id)

	result, err := service.collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})

	if err != nil || result.DeletedCount == 0 {
		return false
	}

	return true
}

func (service *studentDTO) FilterWithNationality(nation string) []Models.Student {

	var students []Models.Student
	// Tạo bộ lọc tìm kiếm cho nationality
	filter := bson.M{"nationality": nation}

	// Sử dụng Find với filter
	cursor, err := service.collection.Find(context.TODO(), filter, &options.FindOptions{
		Hint: "_indexNationality",
	})
	if err != nil {
		log.Println("Error finding students:", err)
		return students // Trả về danh sách rỗng nếu có lỗi
	}
	defer cursor.Close(context.TODO()) // Đảm bảo đóng con trỏ sau khi hoàn tất

	// Lặp qua kết quả và giải mã vào slice students
	for cursor.Next(context.TODO()) {
		var student Models.Student
		if err := cursor.Decode(&student); err != nil {
			log.Println("Error decoding student:", err)
			continue // Bỏ qua sinh viên này và tiếp tục
		}
		// Thêm sinh viên vào danh sách nếu quốc tịch khớp
		if student.Nationality == nation {
			students = append(students, student)
		}
		log.Print(student) // Ghi log thông tin sinh viên
	}

	// Kiểm tra lỗi trong quá trình lặp
	if err := cursor.Err(); err != nil {
		log.Println("Error during cursor iteration:", err)
	}

	return students // Trả về danh sách sinh viên
}
