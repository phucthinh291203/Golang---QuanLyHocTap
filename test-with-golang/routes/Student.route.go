package routes

import (
	controllers "test-with-golang/Controllers"
	"test-with-golang/Models"
	database "test-with-golang/database"
	dto "test-with-golang/dto"
	"test-with-golang/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StudentRoutes(router *gin.Engine, data Models.MyData) {
	studentCollection := database.GetStudentCollection()
	StudentDTO := dto.NewStudent(studentCollection)
	studentController := controllers.NewStudent(StudentDTO)

	studentGroup := router.Group("/admin/students")
	studentGroup.Use(middleware.JWTAuthMiddleWare("Admin"))

	studentGroup.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"Message": "All data fetched successfully", "Data": studentController.FindAll(ctx)})
	})

	studentGroup.GET("/:id", func(ctx *gin.Context) {
		student, class := studentController.FindOne(ctx, data)
		if student.ID != primitive.NilObjectID {
			ctx.JSON(200, gin.H{"Message": "Successfully", "Data student": student, "Data class": class})
			return
		}

		ctx.JSON(500, gin.H{"Message": "Not found"})
	})

	studentGroup.POST("/", func(ctx *gin.Context) {
		result := studentController.Create(ctx, data)
		if result.ID != primitive.NilObjectID {
			ctx.JSON(200, gin.H{"Message": "Student created Successfully", "Student": result})
			return
		}

		ctx.JSON(500, gin.H{"Message": "Failed"})
	})

	studentGroup.PUT("/:id", func(ctx *gin.Context) {
		result := studentController.Update(ctx)
		if result.ID != primitive.NilObjectID {
			ctx.JSON(200, gin.H{"Message": "Student Updated Successfully", "Student_Updated": result})
			return
		}

		ctx.JSON(500, gin.H{"Message": "Failed"})
	})

	studentGroup.DELETE("/:id", func(ctx *gin.Context) {
		result := studentController.Delete(ctx)
		if result {
			ctx.JSON(200, gin.H{"Message": "Student Deleted Successfully"})
			return
		}

		ctx.JSON(500, gin.H{"Message": "Failed"})
	})

	studentGroup.GET("/filter", func(ctx *gin.Context) {
		result := studentController.FilterWithNationality(ctx)
		ctx.JSON(200, gin.H{"Message": "Filter Successfully", "Student_Updated": result})
	})
}
