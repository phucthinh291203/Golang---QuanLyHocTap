package routes

import (
	"log"
	controllers "test-with-golang/Controllers"
	"test-with-golang/Models"
	database "test-with-golang/database"
	dto "test-with-golang/dto"
	"test-with-golang/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ClassRoutes(router *gin.Engine, data Models.MyData) {
	classCollection := database.GetClassCollection()
	ClassDTO := dto.NewClassDto(classCollection)
	classController := controllers.New(ClassDTO, classCollection)

	classGroup := router.Group("/admin/classes")
	{
		classGroup.Use(middleware.JWTAuthMiddleWare("Admin"))
		classGroup.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, classController.FindAll())
		})
		classGroup.GET("/:id", func(ctx *gin.Context) {
			ctx.JSON(200, classController.FindById(ctx))
		})
		classGroup.POST("/", func(ctx *gin.Context) {

			result := classController.Save(ctx, data)
			log.Print(result)
			if result.ID != primitive.NilObjectID {
				ctx.JSON(200, gin.H{"Message": "Tạo lớp thành công", "Data": result})
				return
			}
			ctx.JSON(500, gin.H{"Message": "tạo lớp thất bại"})
		})
		classGroup.PUT("/:id", func(ctx *gin.Context) {
			updatedClass := classController.Update(ctx)

			if updatedClass.ID != primitive.NilObjectID {
				ctx.JSON(200, gin.H{"Message": "Student Updated Successfully", "Student_Updated": updatedClass})
				return
			}

			ctx.JSON(500, gin.H{"Message": "Failed"})

		})
		classGroup.DELETE("/:id", func(ctx *gin.Context) {
			classController.Delete(ctx)
		})

	}
}
