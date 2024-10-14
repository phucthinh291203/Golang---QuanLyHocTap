package controllers

import (

	"test-with-golang/Models"
	dto "test-with-golang/dto"

	"github.com/gin-gonic/gin"
)

type IstudentController interface {
	Create(ctx *gin.Context,data Models.MyData) Models.Student
	FindAll(ctx *gin.Context) []Models.Student
	FindOne(ctx *gin.Context, data Models.MyData) (Models.Student, string)
	Update(ctx *gin.Context) Models.Student
	Delete(ctx *gin.Context) bool

	FilterWithNationality(ctx *gin.Context) []Models.Student
}

type studentController struct {
	service dto.StudentDTO
}

func NewStudent(service dto.StudentDTO) IstudentController {
	return &studentController{
		service: service,
	}
}

func (ctrl *studentController) Create(ctx *gin.Context,data Models.MyData) Models.Student {
	var newData Models.Student
	ctx.BindJSON(&newData)
	err := ctrl.service.Create(newData,data)
	if err != nil {
		return Models.Student{}

	}

	return newData
}
func (ctrl *studentController) FindAll(ctx *gin.Context) []Models.Student {
	return ctrl.service.FindAll()
}

func (ctrl *studentController) FindOne(ctx *gin.Context, data Models.MyData) (Models.Student, string) {
	id := ctx.Param("id")
	student, class := ctrl.service.FindById(id, data)
	return student, class
}
func (ctrl *studentController) Update(ctx *gin.Context) Models.Student {
	var data Models.Student
	id := ctx.Param("id")
	err := ctx.ShouldBindJSON(&data)
	if err == nil {
		result := ctrl.service.Update(id, data)
		return result
	}
	return Models.Student{}
}

func (ctrl *studentController) Delete(ctx *gin.Context) bool {
	id := ctx.Param("id")
	result := ctrl.service.Delete(id)
	return result
}

func (ctrl *studentController) FilterWithNationality(ctx *gin.Context) (result []Models.Student) {
	nationality := ctx.Query("nationality")
	result = ctrl.service.FilterWithNationality(nationality)
	return result
}

