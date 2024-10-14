package controllers

import (
	"context"
	"errors"
	"test-with-golang/Models"
	dto "test-with-golang/dto"

	"go.mongodb.org/mongo-driver/bson"
)

type OtpController struct {
	service dto.OtpDTO
}

func NewOtpController(service dto.OtpDTO) *OtpController {
	return &OtpController{
		service: service,
	}
}

func (ctrl *OtpController) SendOTP(email string, role string, data Models.MyData) error {
	switch role {
	case "Teacher":
		var Teacher Models.Teacher
		err := data.TeacherCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&Teacher)
		if err != nil {
			return errors.New("không tìm thấy email trong Teacher collection")
		}

	case "Admin":
		var Admin Models.Admin
		err := data.UserCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&Admin)
		if err != nil {
			return errors.New("không tìm thấy email trong Admin collection")
		}
	default:
		return errors.New("nhập sai chức vụ")
	}

	err := ctrl.service.SendOTP(email, role, data)
	return err
}

func (ctrl *OtpController) VerifyOTP(email string, otpInput string, data Models.MyData) error {
	err := ctrl.service.VerifyOTP(email, otpInput, data)
	return err
}

func (ctrl *OtpController) ResetPassword(email string, newPassword string, data Models.MyData) error {
	err := ctrl.service.ResetPassword(email, newPassword, data)
	return err
}

func (ctrl *OtpController) SendSMS(phoneNumber string, role string, data Models.MyData) error {
	err := ctrl.service.SendSMS(phoneNumber, role, data)
	return err
}

func (ctrl *OtpController) VerifySMS(phoneNumber string, role string, data Models.MyData) error {
	err := ctrl.service.VerifyOTP_SMS(phoneNumber, role, data)
	return err
}

func (ctrl *OtpController) ResetPassword_SMS(phoneNumber string, role string, data Models.MyData) error {
	err := ctrl.service.ResetPassword_SMS(phoneNumber, role, data)
	return err
}
