package controllers

import (
	"test-with-golang/Models"
	dto "test-with-golang/dto"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubjectController struct{
	service dto.SubjectDTO
}


func NewSubject (service dto.SubjectDTO) *SubjectController{
	return &SubjectController{
		service: service,
	}
}

func (ctrl *SubjectController) CreateNewSubject(ctx *gin.Context) Models.Subject{
	result := ctrl.service.CreateNewSubject(ctx)
	return result
}

func (ctrl *SubjectController) GetAllSubject() [] Models.Subject{
	result := ctrl.service.GetAllSubject()
	return result
}

func (ctrl *SubjectController) UpdateSubject(id primitive.ObjectID ,ctx *gin.Context) Models.Subject{
	
	result := ctrl.service.UpdateSubject(id,ctx)
	return result
}

func (ctrl *SubjectController) DeleteSubject(id primitive.ObjectID) error{
	success := ctrl.service.DeleteSubject(id)
	return success
}