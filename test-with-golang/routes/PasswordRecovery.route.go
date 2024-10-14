package routes

import (
	"net/http"
	controllers "test-with-golang/Controllers"
	"test-with-golang/Models"
	dto "test-with-golang/dto"

	"github.com/gin-gonic/gin"
)

func PasswordRecovery(router *gin.Engine, data Models.MyData) {
	passwordGroup := router.Group("/password")
	OtpDTO := dto.NewOtpDTO(data.OTPCollection)
	OtpController := controllers.NewOtpController(*OtpDTO)
	passwordGroup.POST("/sendOTP", func(ctx *gin.Context) {
		var Request Models.RequestSendOTP
		err := ctx.ShouldBindJSON(&Request)
		if err != nil {
			ctx.JSON(http.StatusFailedDependency, Response{
				Message: "Chưa nhập đủ thông tin",
				Data:    nil,
			})
		}
		err = OtpController.SendOTP(Request.Email, Request.Role, data)
		if err == nil {
			ctx.JSON(http.StatusOK, Response{
				Message: "Đã gửi mã OTP đến email",
				Data:    Request.Email,
			})
		}
	})

	passwordGroup.POST("/verifyOTP", func(ctx *gin.Context) {
		var request Models.RequestVerifyOTP
		err := ctx.ShouldBindJSON(&request)
		if err != nil {
			ctx.JSON(http.StatusNotFound, Response{
				Message: "Chưa nhập đủ thông tin",
			})
			return
		}

		err = OtpController.VerifyOTP(request.Email, request.OTPInput, data)
		if err != nil {
			ctx.JSON(http.StatusConflict, Response{
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, Response{
			Message: "Đã xác thực thành công, vui lòng thực hiện bước tiếp theo",
		})
	})

	passwordGroup.POST("/resetPassword", func(ctx *gin.Context) {
		var request Models.RequestResetPassword
		err := ctx.ShouldBindJSON(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Chưa nhập đủ thông tin",
			})
			return
		}

		if request.NewPassword != request.ConfirmNewPassword {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Chưa nhập đủ thông tin",
			})
			return
		}

		OtpController.ResetPassword(request.Email, request.NewPassword, data)

		ctx.JSON(http.StatusOK, Response{
			Message: "Mật khẩu đã đổi thành công, thực hiện đăng nhập lại",
		})
	})

	passwordGroup.POST("/sendSMS", func(ctx *gin.Context) {
		var request Models.RequestSendSMS
		err := ctx.ShouldBindJSON(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Nhập thiếu nội dung",
			})
			return
		}

		err = OtpController.SendSMS(request.PhoneNumber, request.Role, data)
		if err == nil {
			ctx.JSON(http.StatusOK, Response{
				Message: "Đã gửi mã OTP đến email",
				Data:    request.PhoneNumber,
			})
		}
	})

	passwordGroup.POST("/verifySMS", func(ctx *gin.Context) {
		var request Models.RequestVerifySMS
		err := ctx.ShouldBindJSON(&request)
		if err != nil {
			ctx.JSON(http.StatusNotFound, Response{
				Message: "Chưa nhập đủ thông tin",
			})
			return
		}

		err = OtpController.VerifySMS(request.PhoneNumber, request.OTPInput, data)
		if err != nil {
			ctx.JSON(http.StatusConflict, Response{
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, Response{
			Message: "Đã xác thực thành công, vui lòng thực hiện bước tiếp theo",
		})
	})

	passwordGroup.POST("/resetPasswordSMS", func(ctx *gin.Context) {
		var request Models.RequestResetPasswordSMS
		err := ctx.ShouldBindJSON(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Chưa nhập đủ thông tin",
			})
			return
		}

		if request.NewPassword != request.ConfirmNewPassword {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Chưa nhập đủ thông tin",
			})
			return
		}

		OtpController.ResetPassword_SMS(request.PhoneNumber, request.NewPassword, data)

		ctx.JSON(http.StatusOK, Response{
			Message: "Mật khẩu đã đổi thành công, thực hiện đăng nhập lại",
		})
	})
}
