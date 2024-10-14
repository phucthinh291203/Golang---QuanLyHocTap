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

func TeacherRoute(router *gin.Engine, data Models.MyData) {
	teacherCollection := database.GetTeacherCollection()
	TeacherDTO := dto.NewTeacher(teacherCollection)
	teacherController := controllers.NewTeacher(TeacherDTO)

	teacherGroup := router.Group("/teachers")
	teacherGroup.Use(middleware.JWTAuthMiddleWare("Teacher"))
	teacherGroup.GET("/myClass", func(ctx *gin.Context) {
		teacherId, _ := ctx.Get("teacher_id")
		teacherObjectID, _ := teacherId.(primitive.ObjectID)
		log.Print(teacherObjectID)
		class, students := teacherController.GetStudentOfCurrentClass(teacherObjectID, data)
		ctx.JSON(200, gin.H{"Message": "All data fetched successfully", "Current Class": class, "All students in this class": students})
	})

	adminGroup := router.Group("/admin/teachers")
	adminGroup.Use(middleware.JWTAuthMiddleWare("Admin"))
	adminGroup.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"Message": "All data fetched successfully", "Data": teacherController.FindAll(ctx)})
	})

	adminGroup.POST("/", func(ctx *gin.Context) {

		var newData Models.Teacher
		ctx.BindJSON(&newData)

		err := teacherController.Create(newData, data)
		if err != nil {
			ctx.JSON(500, gin.H{"Message": "Teacher created Failed", "Error": err})
			return
		}
		ctx.JSON(200, gin.H{"Message": "Teacher created Successfully", "Teacher": newData})
	})

	adminGroup.PUT("/:id", func(ctx *gin.Context) {
		result := teacherController.Update(ctx)
		if result.ID != primitive.NilObjectID {
			ctx.JSON(200, gin.H{"Message": "Teacher Updated Successfully", "Teacher_Updated": result})
			return
		}

		ctx.JSON(500, gin.H{"Message": "Failed"})
	})

	adminGroup.DELETE("/:id", func(ctx *gin.Context) {
		result := teacherController.Delete(ctx)
		if result {
			ctx.JSON(200, gin.H{"Message": "Teacher Deleted Successfully"})
			return
		}

		ctx.JSON(500, gin.H{"Message": "Failed"})
	})
}
