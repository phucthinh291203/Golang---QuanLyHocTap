package controllers

import (
	"test-with-golang/Models"
	dto "test-with-golang/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IScoreController interface {
	CreateScoreForStudent(myId primitive.ObjectID, scoreData Models.Score, data Models.MyData) error
	UpdateScore(myId primitive.ObjectID, scoreID string, scoreData Models.Score, data Models.MyData) Models.Score
	DeleteScore(teacherId primitive.ObjectID, idScore string, data Models.MyData) error
}


type ScoreController struct{
	service dto.ScoreDTO
}


func NewScore(service dto.ScoreDTO) *ScoreController {
	return &ScoreController{
		service: service,
	}
}


func (ctrl *ScoreController) CreateScoreForStudent(myId primitive.ObjectID, scoreData Models.Score, data Models.MyData) error {
	err := ctrl.service.CreateScoreForStudent(myId, scoreData, data)
    if err != nil {
        return err
    }

    return nil
}

func (ctrl *ScoreController) UpdateScore(myId primitive.ObjectID,scoreID string, scoreData Models.Score,data Models.MyData) error{
	result := ctrl.service.UpdateScore(myId ,scoreID, scoreData, data)
    return result

}

func (ctrl *ScoreController) DeleteScore(teacherId primitive.ObjectID,idScore string,  data Models.MyData) error{
	err := ctrl.service.DeleteScore(teacherId,idScore,data)
	return err
}

func (ctrl *ScoreController) FindOne(idScore string,  data Models.MyData) Models.Score{
	result := ctrl.service.FindOne(idScore,data)
	return result
}

func (ctrl *ScoreController) FindAll(teacherId primitive.ObjectID, data Models.MyData) []Models.Score{
	result := ctrl.service.FindAll(teacherId,data)
	return result
}