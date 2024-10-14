package routes

import (
	"net/http"
	controllers "test-with-golang/Controllers"
	"test-with-golang/Models"
	database "test-with-golang/database"
	dto "test-with-golang/dto"
	"test-with-golang/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func BangDiemRoutes(router *gin.Engine, data Models.MyData) {
	BangDiemGroup := router.Group("/bangDiem")
	{
		BangDiemGroup.Use(middleware.JWTAuthMiddleWare("Admin"))
		BangDiemCollection := database.GetBangDiemCollection()
		BangDiemDTO := dto.NewBangDiemDTO(BangDiemCollection)
		BangDiemController := controllers.NewBangDiem(*BangDiemDTO)
		BangDiemGroup.POST("/create", func(ctx *gin.Context) {
			newBangDiem, _ := BangDiemController.CreateNewBangDiem(ctx, data)
			if newBangDiem.ID != primitive.NilObjectID {
				// ctx.JSON(200, gin.H{"Message": "In bảng điểm thành công", "Data": newBangDiem})

				ctx.JSON(http.StatusOK, Response{
					Message: "in diem thanh cong",
					Data:    newBangDiem,
				})
				return
			}

			ctx.JSON(500, gin.H{"Message": "In bảng điểm thất bại"})
		})

		BangDiemGroup.GET("/:id", func(ctx *gin.Context) {
			result := BangDiemController.GetBangDiem(ctx, data)
			// ctx.JSON(200, gin.H{"Message": "In bảng điểm thành công", "Data": result})
			ctx.JSON(http.StatusOK, Response{
				Message: "in diem thanh cong",
				Data:    result,
			})
		})
	}
}
