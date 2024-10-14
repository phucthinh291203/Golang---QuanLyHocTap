package controllers

import (
	"net/http"
	"test-with-golang/Models"
	dto "test-with-golang/dto"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type classController struct {
	service    dto.ClassDTO
	collection *mongo.Collection
}

type IClassController interface {
	Save(ctx *gin.Context, data Models.MyData) Models.Class
	FindAll() []Models.Class
	FindById(ctx *gin.Context) Models.Class
	Update(ctx *gin.Context) Models.Class
	Delete(ctx *gin.Context)
}

func New(service dto.ClassDTO, collection *mongo.Collection) IClassController {
	return &classController{
		service:    service,
		collection: collection,
	}
}

// FindAll lấy tất cả lớp học
func (c *classController) FindAll() []Models.Class {
	return c.service.FindAll()
}

// Save lưu lớp học mới
func (c *classController) Save(ctx *gin.Context, data Models.MyData) Models.Class {
	var class Models.Class
	ctx.BindJSON(&class)
	result := c.service.Save(class, data)
	return result
}

func (c *classController) FindById(ctx *gin.Context) Models.Class {
	id := ctx.Param("id")
	class := c.service.FindById(id)
	return class
}

func (c *classController) Update(ctx *gin.Context) Models.Class {
	var data Models.Class
	id := ctx.Param("id")
	err := ctx.ShouldBindJSON(&data)
	if err == nil {
		result := c.service.Update(id, data)
		return result
	}
	return Models.Class{}
}

func (c *classController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	success := c.service.Delete(id)
	if !success {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "id not found or could not be deleted"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Class deleted successfully"})
}
