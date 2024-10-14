package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoaiKiemTra string

const (
	Mieng   LoaiKiemTra = "Mieng"
	PHUT_15 LoaiKiemTra = "15Phut"
	PHUT_45 LoaiKiemTra = "45Phut"
	GK      LoaiKiemTra = "GiuaKy"
	CK      LoaiKiemTra = "CuoiKy"
)

type CachTinhDiem struct {
	ExamType LoaiKiemTra `bson:"exam_type" json:"exam_type"`
	Multiply int         `bson:"mutiply" json:"mutiply"`
}

type Score struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	StudentID       primitive.ObjectID `bson:"student_id"`
	ClassID         primitive.ObjectID `bson:"class_id"`
	SubjectID       primitive.ObjectID `bson:"subject_id"`
	Score           float32            `bson:"score"`
	Semester        string             `bson:"semester"`
	SchoolYearStart int                `bson:"school_year_start"`
	SchoolYearEnd   int                `bson:"school_year_end"`
	Coefficient     CachTinhDiem       `bson:"coefficient"`
	CreatedBy       primitive.ObjectID `bson:"created_by" json:"created_by"`
}

type ScoreResponse struct {
	SubjectName string      `bson:"subject_name" json:"subject_name"`
	ExamType    LoaiKiemTra `bson:"exam_type" json:"exam_type"`
	Score       float32     `bson:"score" json:"score"`
	Multiply    int         `bson:"multiply" json:"multiply"`
}

type ScoreResponseOutput struct {
	ExamType LoaiKiemTra `bson:"exam_type" json:"exam_type"`
	Score    float32     `bson:"score" json:"score"`
}

func GetHeSoByExamType(examType LoaiKiemTra) int {
	switch examType {
	case PHUT_45:
		return 2
	case GK:
		return 2
	case CK:
		return 4
	default:
		return 1 // Giá trị mặc định cho các kỳ thi khác (ví dụ: Miệng, 15 phút)
	}
}
