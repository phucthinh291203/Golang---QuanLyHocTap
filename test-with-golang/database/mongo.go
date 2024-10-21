package database

import (
	"context"
	"log"
	"os"
	Models "test-with-golang/Models"
	auth "test-with-golang/auth"

	"github.com/robfig/cron/v3"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Biến toàn cục để chứa client và các collection
var (
	mongoClient *mongo.Client
	data        Models.MyData
)

func ConnectToDatabase() {
	godotenv.Load()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("DATABASE_URL")))
	if err != nil {
		log.Fatal("Kết nối đến database thất bại:", err)
		return
	}

	// Kiểm tra kết nối
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Kết nối đến database thất bại:", err)
		return
	}

	mongoClient = client

	data = Models.MyData{
		ClassCollection:             mongoClient.Database(GetDBName()).Collection("classCollection"),
		StudentCollection:           mongoClient.Database(GetDBName()).Collection("studentCollection"),
		UserCollection:              mongoClient.Database(GetDBName()).Collection("userCollection"),
		ScoreCollection:             mongoClient.Database(GetDBName()).Collection("scoreCollection"),
		SubjectCollection:           mongoClient.Database(GetDBName()).Collection("subjectCollection"),
		TeacherCollection:           mongoClient.Database(GetDBName()).Collection("teacherCollection"),
		BangDiemCollection:          mongoClient.Database(GetDBName()).Collection("bangDiemCollection"),
		OTPCollection:               mongoClient.Database(GetDBName()).Collection("OTPCollection"),
		SMSCollection:               mongoClient.Database(GetDBName()).Collection("SMSCollection"),
		FileCollection:              mongoClient.Database(GetDBName()).Collection("FileCollection"),
		HistoryUploadedCollection:   mongoClient.Database(GetDBName()).Collection("HistoryUploadedCollection"),
		HistoryDownloadedCollection: mongoClient.Database(GetDBName()).Collection("HistoryDownloadedCollection"),
	}

	if mongoClient == nil {
		log.Fatal("mongoClient chưa được khởi tạo")
	} else {
		log.Println("mongoClient đã khởi tạo thành công")
	}
	log.Print("Kết nối đến database thành công!")

}

// Hàm trả về collection của class
func GetClassCollection() *mongo.Collection {
	return data.ClassCollection
}

// Hàm trả về collection của student
func GetStudentCollection() *mongo.Collection {
	return data.StudentCollection
}

func GetUserCollection() *mongo.Collection {
	return data.UserCollection
}

func GetScoreCollection() *mongo.Collection {
	return data.ScoreCollection
}

func GetSubjectCollection() *mongo.Collection {
	return data.SubjectCollection
}

func GetTeacherCollection() *mongo.Collection {
	return data.TeacherCollection
}

func GetBangDiemCollection() *mongo.Collection {
	return data.BangDiemCollection
}

func GetOTPCollection() *mongo.Collection {
	return data.OTPCollection
}

func GetFileCollection() *mongo.Collection {
	return data.FileCollection
}

func GetHistoryUploadedCollection() *mongo.Collection {
	return data.HistoryUploadedCollection
}

func GetHistoryDownloadedCollection() *mongo.Collection {
	return data.HistoryDownloadedCollection
}

func GetDBName() string {
	return os.Getenv("dbName")
}

func GetClient() *mongo.Client {
	return mongoClient
}

func GetData() Models.MyData {
	return data
}

func SaveUser(user *Models.Admin) error {
	_, err := data.UserCollection.InsertOne(context.TODO(), user)
	return err
}

func FindUserByUsername(username string) (Models.Admin, error) {
	var user Models.Admin
	err := data.UserCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return Models.Admin{}, err
	}
	log.Print(user)
	return user, nil
}

func FindTeacherByUsername(username string) (Models.Teacher, error) {
	var teacher Models.Teacher
	err := data.TeacherCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&teacher)
	if err != nil {
		return Models.Teacher{}, err
	}
	log.Print(teacher)
	return teacher, nil
}

func FilterScoreFromCollection(semester string, startYear int, endYear int, studentId primitive.ObjectID) []Models.Score {
	filter := bson.M{
		"student_id":        studentId,
		"semester":          semester,
		"school_year_start": startYear,
		"school_year_end":   endYear,
	}
	cursor, err := data.ScoreCollection.Find(context.TODO(), filter)
	if err != nil {
		return []Models.Score{}
	}
	defer cursor.Close(context.TODO())

	var scores []Models.Score
	for cursor.Next(context.TODO()) {
		var score Models.Score
		if err := cursor.Decode(&score); err != nil {
			return []Models.Score{}
		}
		scores = append(scores, score)
	}

	return scores
}

func TinhDiemTrungBinh(scores []Models.Score) float32 {
	var totalWeightedScore float32 // Tổng điểm nhân hệ số
	var totalHeSo int              // Tổng hệ số

	// Duyệt qua từng điểm trong danh sách
	for _, score := range scores {
		heSo := Models.GetHeSoByExamType(score.Coefficient.ExamType) // Lấy hệ số dựa trên loại kỳ thi

		// môn học nào nếu hệ số thay đổi , sẽ thay đổi ở đây

		totalWeightedScore += score.Score * float32(heSo) // Cộng dồn điểm nhân hệ số
		totalHeSo += heSo                                 // Cộng dồn hệ số
	}

	// Tính điểm trung bình tổng
	if totalHeSo == 0 {
		return 0 // Tránh chia cho 0
	}
	return totalWeightedScore / float32(totalHeSo)
}

func StartOTPCleaner(data Models.MyData) {
	c := cron.New()
	c.AddFunc("*/2 * * * *", func() {
		auth.RemoveExpiredOTPs(data)
	})

	//Bắt đầu thực hiện
	c.Start()

	//Chạy vô tận
	select {}
}
