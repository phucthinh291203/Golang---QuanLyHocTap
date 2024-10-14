package controllers

import (
	"context"
	"log"
	"test-with-golang/Models"
	dto "test-with-golang/dto"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BangDiemController struct {
	service dto.BangDiemDTO
}

func NewBangDiem(service dto.BangDiemDTO) *BangDiemController {
	return &BangDiemController{
		service: service,
	}
}

func (ctrl *BangDiemController) CreateNewBangDiem(ctx *gin.Context, data Models.MyData) (Models.BangDiem, error) {
	var searchdata Models.TraCuu
	err := ctx.ShouldBindJSON(&searchdata)
	if err != nil {
		return Models.BangDiem{}, err
	}

	var student Models.Student
	err = data.StudentCollection.FindOne(context.TODO(), bson.M{"_id": searchdata.StudentID}).Decode(&student)
	if err != nil {
		log.Print("Không tìm thấy student trong collection")
		return Models.BangDiem{},err
	}
	result := ctrl.service.CreateNewBangDiem(ctx, searchdata, data)
	return result,nil
}


func (ctrl *BangDiemController) GetBangDiem(ctx *gin.Context, data Models.MyData) Models.BangDiemOutPut {
	id := ctx.Param("id")
	idObjectID,_ := primitive.ObjectIDFromHex(id)
	result := ctrl.service.GetBangDiem(idObjectID,data)
	return result
}
