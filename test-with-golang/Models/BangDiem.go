package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type XepLoai string

const (
	Gioi      XepLoai = "Gioi"
	Kha       XepLoai = "Kha"
	Yeu       XepLoai = "Yeu"
	TrungBinh XepLoai = "TrungBinh"
	Kem       XepLoai = "Kem"
)

type BangDiem struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	StudentID       primitive.ObjectID `bson:"student_id"`
	SchoolYearStart int                `bson:"school_year_start"`
	SchoolYearEnd   int                `bson:"school_year_end"`
	Semester        string             `bson:"semester"`
	ScoreResponse   []ScoreResponse    `bson:"score_ids"`
	AverageScore    float32            `bson:"average_score"`
	Grade           XepLoai            `bson:"grade"`
}

type BangDiemOutPut struct {
	ID              primitive.ObjectID       `bson:"_id,omitempty" json:"_id"`
	StudentID       primitive.ObjectID       `bson:"student_id" json:"student_id"`
	SchoolYearStart int                      `bson:"school_year_start" json:"school_year_start"`
	SchoolYearEnd   int                      `bson:"school_year_end" json:"school_year_end"`
	Semester        string                   `bson:"semester" json:"semester"`
	ScoreResponse   []ProjectedScoreResponse `bson:"score_ids" json:"score_ids"`
	AverageScore    float32                  `bson:"average_score" json:"average_score"`
	Grade           XepLoai                  `bson:"grade" json:"grade"`
}
type ProjectedScoreResponse struct {
	SubjectName    string                `bson:"subject_name" json:"subject_name"`
	Score          []ScoreResponseOutput `bson:"score" json:"score"`
	AverageSubject float64               `bson:"average_subject" json:"average_subject"`
}

type TraCuu struct {
	StudentID       primitive.ObjectID `bson:"student_id"`
	SchoolYearStart int                `bson:"school_year_start"`
	SchoolYearEnd   int                `bson:"school_year_end"`
	Semester        string             `bson:"semester"`
}
