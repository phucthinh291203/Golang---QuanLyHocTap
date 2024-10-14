package services

import (
	"context"
	"log"
	"test-with-golang/Models"
	database "test-with-golang/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BangDiemDTO struct {
	collection *mongo.Collection
}

func NewBangDiemDTO(collection *mongo.Collection) *BangDiemDTO {
	return &BangDiemDTO{
		collection: collection,
	}
}

func (service *BangDiemDTO) CreateNewBangDiem(ctx *gin.Context, searchdata Models.TraCuu, data Models.MyData) Models.BangDiem {
	// existed := CheckingExist(searchdata.Semester, searchdata.SchoolYearStart, searchdata.SchoolYearEnd, searchdata.StudentID, data)
	// if existed == true {
	// 	log.Print("Đã tạo bảng điểm cho học sinh kỳ này năm này rồi")
	// 	return Models.BangDiem{}
	// }
	scores := database.FilterScoreFromCollection(searchdata.Semester, searchdata.SchoolYearStart, searchdata.SchoolYearEnd, searchdata.StudentID)

	var newBangDiem Models.BangDiem
	newBangDiem.StudentID = searchdata.StudentID
	newBangDiem.SchoolYearStart = (searchdata.SchoolYearStart)
	newBangDiem.SchoolYearEnd = (searchdata.SchoolYearEnd)
	newBangDiem.Semester = searchdata.Semester
	newBangDiem.AverageScore = database.TinhDiemTrungBinh(scores)

	newBangDiem.Grade = XepLoai(newBangDiem.AverageScore)

	log.Print(newBangDiem)

	var scoreResponses []Models.ScoreResponse
	for _, score := range scores {
		var subject Models.Subject
		err := data.SubjectCollection.FindOne(context.TODO(), bson.M{"_id": score.SubjectID}).Decode(&subject)
		if err != nil {
			log.Printf("Không tìm thấy môn học với ID: %v", score.SubjectID)
			continue
		}

		scoreResponses = append(scoreResponses, Models.ScoreResponse{
			SubjectName: subject.SubjectName, // Giả định rằng trường tên môn học trong bảng Subject là "Name"
			Score:       score.Score,
			ExamType:    score.Coefficient.ExamType,
			Multiply:    score.Coefficient.Multiply,
		})
	}

	newBangDiem.ScoreResponse = scoreResponses

	result, err := data.BangDiemCollection.InsertOne(context.TODO(), newBangDiem)
	if err == nil {
		newBangDiem.ID = result.InsertedID.(primitive.ObjectID)
		return newBangDiem
	}

	return Models.BangDiem{}
}

func XepLoai(average_score float32) Models.XepLoai {
	switch {
	case average_score >= 8.0:
		return "Gioi"
	case average_score >= 6.5:
		return "Kha"
	case average_score >= 5.0:
		return "TrungBinh"
	case average_score >= 3.5:
		return "Yeu"
	default:
		return "Kem"
	}
}

func (service *BangDiemDTO) GetBangDiem(idBangDiem primitive.ObjectID, data Models.MyData) Models.BangDiemOutPut {

	// Pipeline ngắn hơn chỉ gom nhóm và định dạng lại ScoreResponse
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: idBangDiem}}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$score_ids"},
			{Key: "includeArrayIndex", Value: "index"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$score_ids.subject_name"},
			{Key: "ScoreResponse", Value: bson.D{
				{Key: "$push", Value: bson.D{
					{Key: "exam_type", Value: "$score_ids.exam_type"},
					{Key: "score", Value: "$score_ids.score"},
				}},
			}},
			{Key: "SumOfScore", Value: bson.D{
				{Key: "$sum", Value: bson.D{
					{Key: "$multiply", Value: bson.A{
						"$score_ids.score",
						"$score_ids.multiply",
					}},
				}},
			}},
			{Key: "SumOfMutiply", Value: bson.D{
				{Key: "$sum", Value: "$score_ids.multiply"},
			}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0}, // Không trả về _id
			{Key: "subject_name", Value: "$_id"},
			{Key: "score", Value: bson.D{
				{Key: "$map", Value: bson.M{
					"input": "$ScoreResponse",
					"as":    "item",
					"in": bson.M{
						"exam_type": "$$item.exam_type", // Truy cập đúng trường LoaiKiemTra
						"score":     "$$item.score",    // Truy cập đúng trường Score
					},
				}},
			}},
			{Key: "average_subject", Value: bson.M{
				"$divide": bson.A{
					bson.M{"$toDouble": "$SumOfScore"},
					bson.M{"$toDouble": "$SumOfMutiply"},
				},
			}},
		}}},
	}

	var bangDiemOutPut Models.BangDiemOutPut
	var bangDiem Models.BangDiem
	data.BangDiemCollection.FindOne(context.TODO(), bson.M{"_id": idBangDiem}).Decode(&bangDiem)

	bangDiemOutPut.ID = bangDiem.ID
	bangDiemOutPut.StudentID = bangDiem.StudentID
	bangDiemOutPut.SchoolYearStart = bangDiem.SchoolYearStart
	bangDiemOutPut.SchoolYearEnd = bangDiem.SchoolYearEnd
	bangDiemOutPut.Semester = bangDiem.Semester

	var results []Models.ProjectedScoreResponse
	var AllAverageSubject []float64
	cursor, err := data.BangDiemCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var scoreResponse Models.ProjectedScoreResponse
		if err := cursor.Decode(&scoreResponse); err != nil {
			log.Fatal(err)
		}
		results = append(results, scoreResponse) // Append vào slice
		AllAverageSubject = append(AllAverageSubject, scoreResponse.AverageSubject)
	}

	bangDiemOutPut.ScoreResponse = results
	bangDiemOutPut.AverageScore = float32(TinhDiemTrungBinh(AllAverageSubject))
	bangDiemOutPut.Grade = XepLoai_test(AllAverageSubject, bangDiemOutPut.AverageScore)
	return bangDiemOutPut
}

func CheckingExist(semester string, startYear int, endYear int, studentId primitive.ObjectID, data Models.MyData) bool {
	filter := bson.M{
		"student_id":        studentId,
		"semester":          semester,
		"school_year_start": startYear,
		"school_year_end":   endYear,
	}
	log.Print(semester)
	log.Print(studentId)
	log.Print(startYear)
	log.Print(endYear)
	var existedData Models.BangDiem
	err := data.BangDiemCollection.FindOne(context.TODO(), filter).Decode(&existedData)
	log.Print(existedData)
	return err == nil
}

func TinhDiemTrungBinh(allSubject []float64) float64 {
	var sum float64
	for _, subject := range allSubject {
		sum += subject
	}

	return sum / float64(len(allSubject))
}

func XepLoai_test(allSubject []float64, averageScore float32) Models.XepLoai {
	var isGioi bool = true
	for _, subject := range allSubject {
		if subject < 6.5 {
			isGioi = false
		}
	}

	if isGioi {
		switch {
		case averageScore >= 8.0:
			return "Gioi"
		}
	} else {
		switch {
		case averageScore >= 6.5:
			return "Kha"
		case averageScore >= 5.0:
			return "TrungBinh"
		case averageScore >= 3.5:
			return "Yeu"
		default:
			return "Kem"
		}
	}
	return ""
}
