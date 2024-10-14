package services

import (
	"test-with-golang/Models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassDTO interface {
	Save(Models.Class, Models.MyData) Models.Class
	FindAll() []Models.Class
	FindById(id string) Models.Class
	Update(id string, change Models.Class) Models.Class
	Delete(id string) bool
}

type StudentDTO interface {
	Create(newData Models.Student, data Models.MyData) error
	FindAll() []Models.Student
	FindById(id string, data Models.MyData) (Models.Student, string)
	Update(id string, change Models.Student) Models.Student
	Delete(id string) bool
	FilterWithNationality(nation string) []Models.Student
}

type TeacherDTO interface {
	Create(Models.Teacher, Models.MyData) error
	FindAll() []Models.Teacher
	Update(id string, change Models.Teacher) Models.Teacher
	Delete(id string) bool
	GetStudentOfCurrentClass(myId primitive.ObjectID, data Models.MyData) ([]Models.Class, []Models.Student)
}
