package controllers

import (
	"test-with-golang/Models"
	dto "test-with-golang/dto"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IteacherController interface {
	Create(newData Models.Teacher,data Models.MyData) error
	FindAll(ctx *gin.Context) []Models.Teacher
	Update(ctx *gin.Context) Models.Teacher
	Delete(ctx *gin.Context) bool

	GetStudentOfCurrentClass(myId primitive.ObjectID, data Models.MyData) ([]Models.Class, []Models.Student)
}

type teacherController struct {
	service dto.TeacherDTO
}

func NewTeacher(service dto.TeacherDTO) IteacherController {
	return &teacherController{
		service: service,
	}
}

func (ctrl *teacherController) Create(newData Models.Teacher,data Models.MyData) error {
	
	err := ctrl.service.Create(newData,data)
	return err
}
func (ctrl *teacherController) FindAll(ctx *gin.Context) []Models.Teacher {
	return ctrl.service.FindAll()
}

func (ctrl *teacherController) Update(ctx *gin.Context) Models.Teacher {
	var data Models.Teacher
	id := ctx.Param("id")
	err := ctx.ShouldBindJSON(&data)
	if err == nil {
		result := ctrl.service.Update(id, data)
		return result
	}
	return Models.Teacher{}
}

func (ctrl *teacherController) Delete(ctx *gin.Context) bool {
	id := ctx.Param("id")
	result := ctrl.service.Delete(id)
	return result
}

func (ctrl *teacherController) GetStudentOfCurrentClass(myId primitive.ObjectID, data Models.MyData) ([]Models.Class, []Models.Student) {
	class,students := ctrl.service.GetStudentOfCurrentClass(myId,data)
	return class,students
}

