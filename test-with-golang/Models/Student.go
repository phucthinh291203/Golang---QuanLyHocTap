package Models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID                        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name                      string             `bson:"name"`
	DateOfBirth               time.Time          `bson:"date_of_birth"`
	ClassID                   primitive.ObjectID `bson:"class_id" json:"class_id"`
	Email                     string             `bson:"email"`
	PhoneNumber               string             `bson:"phone_number"`
	Address                   string             `bson:"address"`
	EnrollmentDate            time.Time          `bson:"enrollment_date"`
	Gender                    string             `bson:"gender"`
	Nationality               string             `bson:"nationality"`
	Grade                     float64            `bson:"grade"`
	GuardianName              string             `bson:"guardian_name"`
	Scholarship               bool               `bson:"scholarship"`
	ExtracurricularActivities []string           `bson:"extracurricular_activities"`
	CourseLoad                []string           `bson:"course_load"`           // Danh sách các khóa học đang theo học
	AttendanceRate            float64            `bson:"attendance_rate"`       // Tỷ lệ đi học
	DisciplinaryActions       []string           `bson:"disciplinary_actions"`  // Các hành động kỷ luật đã bị áp dụng
	ExtracurricularClubs      []string           `bson:"extracurricular_clubs"` // Các câu lạc bộ ngoại khóa tham gia
	Height                    float64            `bson:"height"`                // Chiều cao (cm)
	Weight                    float64            `bson:"weight"`                // Cân nặng (kg)
	MedicalConditions         []string           `bson:"medical_conditions"`    // Các tình trạng y tế
	FavoriteSubjects          []string           `bson:"favorite_subjects"`     // Các môn học yêu thích
}
