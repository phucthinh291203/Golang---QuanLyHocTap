package services

import (
	"context"
	"errors"
	"log"
	Models "test-with-golang/Models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ScoreDTO struct {
	collection *mongo.Collection
}

func NewScoreDTO(collection *mongo.Collection) *ScoreDTO {
	return &ScoreDTO{
		collection: collection,
	}
}

func (service *ScoreDTO) CreateScoreForStudent(teacherId primitive.ObjectID, scoreData Models.Score, data Models.MyData) error {
	// Tìm giáo viên hiện tại
	var teacher Models.Teacher
	err := data.TeacherCollection.FindOne(context.TODO(), bson.M{"_id": teacherId}).Decode(&teacher)
	if err != nil {
		log.Print("Teacher not found")
		return errors.New("teacher not found")
	}

	log.Print(scoreData.SubjectID)

	//Tìm xem môn học có trong collection không
	var subject Models.Subject
	err = data.SubjectCollection.FindOne(context.TODO(), bson.M{"_id": scoreData.SubjectID}).Decode(&subject)
	log.Print(subject)
	if err != nil {
		log.Print("Subject not found")
		return errors.New("subject not found")
	}

	// Tìm toàn bộ lớp của giáo viên đang dạy
	var classes []Models.Class
	cursor, _ := data.ClassCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": teacher.ClassID}})

	// Kiểm tra nếu cursor không có dữ liệu
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var eachClass Models.Class
		_ = cursor.Decode(&eachClass)
		classes = append(classes, eachClass)
	}

	// Tìm lớp đã nhập vào xem có match với toàn bộ lớp giáo viên ko
	var targetClass Models.Class
	for _, class := range classes {
		if class.ID == scoreData.ClassID {
			targetClass = class
			break
		}
	}

	//Nếu lớp nhập vào không match với gv thì trả lỗi
	if targetClass.ID == primitive.NilObjectID {
		log.Print("Lớp id nhập vào không thuộc lớp của bạn đang dạy")
		return errors.New("lớp id nhập vào không thuộc lớp của bạn đang dạy")
	}

	//Gán năm học của lớp hiện tại vô score
	scoreData.SchoolYearStart = targetClass.SchoolYearStart
	scoreData.SchoolYearEnd = targetClass.SchoolYearEnd

	// Tìm toàn bộ môn của giáo viên dạy
	var subjects []Models.Subject
	cursor, _ = data.SubjectCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": teacher.SubjectID}})

	// Kiểm tra nếu cursor không có dữ liệu
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var eachSubject Models.Subject
		_ = cursor.Decode(&eachSubject)
		subjects = append(subjects, eachSubject)
	}

	// Tìm môn đã nhập vào xem có match với toàn bộ môn giáo viên phụ trách ko
	var targetSubject Models.Subject
	for _, subject := range subjects {
		if subject.ID == scoreData.SubjectID {
			targetSubject = subject
			break
		}
	}

	//Nếu môn nhập vào không match với gv thì trả lỗi
	if targetSubject.ID == primitive.NilObjectID {
		log.Print("Môn id nhập vào không phải của bạn phụ trách")
		return errors.New("môn id nhập vào không phải của bạn phụ trách")
	}

	// Kiểm tra xem học sinh có thuộc lớp này không
	if !contains(targetClass.StudentID, scoreData.StudentID) {
		log.Print("Student not in teacher's class")
		return errors.New("student not in teacher's class")
	}

	//Tạo hệ số
	scoreData.Coefficient.Multiply = Models.GetHeSoByExamType(scoreData.Coefficient.ExamType)
	//Gán id người tạo bài kiểm tra
	scoreData.CreatedBy = teacherId

	// Lưu điểm vào database
	scoreCollection := data.ScoreCollection
	_, err = scoreCollection.InsertOne(context.TODO(), scoreData)
	if err != nil {
		log.Print("Error inserting score data")
		return err
	}

	return nil
}

func (service *ScoreDTO) UpdateScore(teacherId primitive.ObjectID, scoreID string, scoreData Models.Score, data Models.MyData) error {
	// Tìm giáo viên hiện tại
	var teacher Models.Teacher
	err := data.TeacherCollection.FindOne(context.TODO(), bson.M{"_id": teacherId}).Decode(&teacher)
	if err != nil {
		log.Print("Teacher not found")
		return errors.New("teacher not found")
	}

	//Tìm xem môn học có trong collection không
	subjectFound := data.SubjectCollection.FindOne(context.TODO(), bson.M{"_id": scoreData.SubjectID})
	if subjectFound.Err() != nil {
		log.Print("Subject not found")
		return errors.New("subject not found")
	}

	// Tìm toàn bộ lớp của giáo viên đang dạy
	var classes []Models.Class
	cursor, _ := data.ClassCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": teacher.ClassID}})

	// Kiểm tra nếu cursor không có dữ liệu
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var eachClass Models.Class
		_ = cursor.Decode(&eachClass)
		classes = append(classes, eachClass)
	}

	// Tìm lớp đã nhập vào xem có match với toàn bộ lớp giáo viên ko
	var targetClass Models.Class
	for _, class := range classes {
		if class.ID == scoreData.ClassID {
			targetClass = class
			break
		}
	}

	//Nếu lớp nhập vào không match với gv thì trả lỗi
	if targetClass.ID == primitive.NilObjectID {
		log.Print("Lớp id nhập vào không thuộc lớp của bạn đang dạy")
		return errors.New("lớp id nhập vào không thuộc lớp của bạn đang dạy")
	}

	//Gán năm học của lớp hiện tại vô score
	scoreData.SchoolYearStart = targetClass.SchoolYearStart
	scoreData.SchoolYearEnd = targetClass.SchoolYearEnd

	// Tìm toàn bộ môn của giáo viên dạy
	var subjects []Models.Subject
	cursor, _ = data.SubjectCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": teacher.SubjectID}})

	// Kiểm tra nếu cursor không có dữ liệu
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var eachSubject Models.Subject
		_ = cursor.Decode(&eachSubject)
		subjects = append(subjects, eachSubject)
	}

	// Tìm môn đã nhập vào xem có match với toàn bộ môn giáo viên phụ trách ko
	var targetSubject Models.Subject
	for _, subject := range subjects {
		if subject.ID == scoreData.SubjectID {
			targetSubject = subject
			break
		}
	}

	//Nếu môn nhập vào không match với gv thì trả lỗi
	if targetSubject.ID == primitive.NilObjectID {
		log.Print("Môn id nhập vào không phải của bạn phụ trách")
		return errors.New("môn id nhập vào không phải của bạn phụ trách")
	}

	// Kiểm tra xem học sinh có thuộc lớp này không
	if !contains(targetClass.StudentID, scoreData.StudentID) {
		log.Print("Student not in teacher's class")
		return errors.New("student not in teacher's class")
	}

	// Tìm kiếm điểm người đăng nhập đã chấm
	ObjectID, _ := primitive.ObjectIDFromHex(scoreID)
	filter := bson.M{"_id": ObjectID, "createdBy": teacherId}

	if scoreData.CreatedBy != teacherId {
		log.Print("không được phép đổi id người chấm")
		return errors.New("không được phép đổi id người chấm")
	}

	// Tùy chọn: Trả về tài liệu sau khi cập nhật
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedScore Models.Score

	// Gọi findOneAndUpdate để cập nhật dữ liệu trong MongoDB
	err = data.ScoreCollection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": scoreData}, options).Decode(&updatedScore)
	if err != nil {
		log.Println("Error updating score:", err)
		return errors.New("lỗi trong quá trình update")
	}

	log.Print(updatedScore)
	return nil
}

func (service *ScoreDTO) DeleteScore(teacherId primitive.ObjectID, idScore string, data Models.MyData) error {
	objectID, _ := primitive.ObjectIDFromHex(idScore)
	log.Print(teacherId)
	log.Print(objectID)
	filter := bson.M{"_id": objectID, "createdBy": teacherId}
	_, err := data.ScoreCollection.DeleteOne(context.TODO(), filter)
	return err
}

func (service *ScoreDTO) FindOne(idScore string, data Models.MyData) Models.Score {
	var ScoreData Models.Score
	idScoreObjectId, _ := primitive.ObjectIDFromHex(idScore)
	service.collection.FindOne(context.TODO(), bson.M{"_id": idScoreObjectId}).Decode(&ScoreData)
	return ScoreData
}

func (service *ScoreDTO) FindAll(teacherId primitive.ObjectID, data Models.MyData) []Models.Score {
	var allScore []Models.Score
	filter := bson.M{"createdBy": teacherId}

	cursor, err := service.collection.Find(context.TODO(), filter)
	if err != nil {
		return []Models.Score{} // Trả về lỗi nếu có
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var scoreData Models.Score
		_ = cursor.Decode(&scoreData)
		allScore = append(allScore, scoreData)
	}

	return allScore
}


