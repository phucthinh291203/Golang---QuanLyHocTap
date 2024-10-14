package services

import (
	"context"
	"errors"
	"log"
	"test-with-golang/Models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type teacherDTO struct {
	collection *mongo.Collection
}

func NewTeacher(collection *mongo.Collection) TeacherDTO {

	dto := &teacherDTO{
		collection: collection,
	}
	return dto
}

func (service *teacherDTO) Create(newData Models.Teacher, data Models.MyData) error {
	// AllClassId := newData.ClassID
	// filter := bson.M{"_id": bson.M{"$in": AllClassId}}
	// cursor, _ := data.ClassCollection.Find(context.TODO(), filter)
	// // Kiểm tra nếu cursor không có dữ liệu
	// defer cursor.Close(context.TODO())

	// var allClass []Models.Class

	// for cursor.Next(context.TODO()) {
	// 	var eachClass Models.Class
	// 	_ = cursor.Decode(&eachClass)
	// 	allClass = append(allClass, eachClass)
	// }

	// if len(allClass) != len(AllClassId) {
	// 	log.Print("Có lớp giáo viên dạy không tồn tại trong collection")
	// 	return errors.New("có lớp giáo viên dạy không tồn tại trong collection")
	// }

	var existedTeacher Models.Teacher
	err := service.collection.FindOne(context.TODO(), bson.M{"username": newData.Username}).Decode(&existedTeacher)
	if err == nil && existedTeacher.ID != primitive.NilObjectID {
		log.Print("Đã có tài khoản giáo viên tồn tại")
		return errors.New("student not in teacher's class")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newData.Password), bcrypt.DefaultCost)
	teacher := Models.Teacher{
		Username:    newData.Username,
		Password:    string(hashedPassword), // Có thể mã hóa password nếu cần
		TeacherName: newData.TeacherName,
		Email:       newData.Email,
	}
	_, err = service.collection.InsertOne(context.TODO(), teacher)
	return err
}

func (service *teacherDTO) FindAll() []Models.Teacher {
	var allData []Models.Teacher
	cursor, _ := service.collection.Find(context.TODO(), bson.D{})
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var data Models.Teacher
		_ = cursor.Decode(&data)
		allData = append(allData, data)
	}

	return allData
}

func (service *teacherDTO) Update(id string, change Models.Teacher) Models.Teacher {
	updateData := bson.M{
		"email":       change.Email,
		"teacher_name": change.TeacherName,
		"class_ids":   change.ClassID,
		"subject_ids": change.SubjectID,
	}
	var result Models.Teacher
	ObjectID, _ := primitive.ObjectIDFromHex(id)
	err := service.collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": ObjectID},
		bson.M{"$set": updateData}, //Chỉ được cập nhật lớp, họ tên, password và username không được đổi
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&result)

	log.Print(ObjectID)
	log.Print(result)
	if err != nil {
		return Models.Teacher{} // Trả về lớp trống nếu không tìm thấy
	}
	return result // Trả về lớp đã cập nhật
}

func (service *teacherDTO) Delete(id string) bool {
	objectID, _ := primitive.ObjectIDFromHex(id)

	result, err := service.collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})

	if err != nil || result.DeletedCount == 0 {
		return false
	}

	return true
}

func (service *teacherDTO) GetStudentOfCurrentClass(myId primitive.ObjectID, data Models.MyData) ([]Models.Class, []Models.Student) {

	//Tìm ra giáo viên hiện tại
	var teacher Models.Teacher
	err := service.collection.FindOne(context.TODO(), bson.M{"_id": myId}).Decode(&teacher)
	if err != nil {
		return []Models.Class{}, nil
	}

	log.Print(teacher)

	if len(teacher.ClassID) == 0 {
		log.Print("Teacher has no classes assigned")
		return []Models.Class{}, nil
	}

	//Tìm ra các lớp mà họ dạy
	var classes []Models.Class
	classCollection := data.ClassCollection
	cursor, err := classCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": teacher.ClassID}})
	log.Print(teacher.ClassID)

	if err != nil {
		log.Print("Error finding classes for teacher: ", err)
		return []Models.Class{}, nil
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var eachClass Models.Class
		_ = cursor.Decode(&eachClass)
		classes = append(classes, eachClass)
	}

	log.Print("Thong tin cac lop hoc dang day: ", classes)

	var allStudents []Models.Student

	// Duyệt qua từng lớp để lấy học sinh
	for _, class := range classes {
		studentCollection := data.StudentCollection
		studentCursor, err := studentCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": class.StudentID}})
		if err != nil {
			log.Print("Error finding students for class: ", err)
			continue // Bỏ qua nếu gặp lỗi
		}
		defer studentCursor.Close(context.TODO())

		for studentCursor.Next(context.TODO()) {
			var eachStudentInClass Models.Student
			_ = studentCursor.Decode(&eachStudentInClass)
			allStudents = append(allStudents, eachStudentInClass)
		}
	}

	return classes, allStudents
}

// 	objectID, _ := primitive.ObjectIDFromHex(idScore)
// 	log.Print(teacherId)
// 	log.Print(objectID)
// 	filter := bson.M{"_id": objectID, "createdBy": teacherId}
// 	_, err := data.ScoreCollection.DeleteOne(context.TODO(), filter)
// 	return err
// }

func contains(allID []primitive.ObjectID, eachID primitive.ObjectID) bool {
	for _, id := range allID {
		if id == eachID {
			return true
		}
	}
	return false
}
