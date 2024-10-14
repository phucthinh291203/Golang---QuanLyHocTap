package routes

import (
	"log"
	controllers "test-with-golang/Controllers"
	"test-with-golang/Models"
	"test-with-golang/middleware"

	dto "test-with-golang/dto"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SubjectRoute(router *gin.Engine, data Models.MyData) {
	subjectCollection := data.SubjectCollection
	subjectDTO := dto.NewSubjectDTO(subjectCollection)
	subjectController := controllers.NewSubject(*subjectDTO)

	subjectGroup := router.Group("admin/subjects")
	subjectGroup.Use(middleware.JWTAuthMiddleWare("Admin"))
	subjectGroup.POST("/", func(ctx *gin.Context) {
		createdData := subjectController.CreateNewSubject(ctx)
		log.Print(createdData)
		if createdData.ID != primitive.NilObjectID {
			ctx.JSON(200, gin.H{"Message": "Created successfully", "Data created": createdData})
			return
		}
		ctx.JSON(500, gin.H{"Message": "Created Failed"})
	})

	subjectGroup.PUT("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		idObject, _ := primitive.ObjectIDFromHex(id)
		updatedData := subjectController.UpdateSubject(idObject, ctx)
		if updatedData.ID != primitive.NilObjectID {
			ctx.JSON(200, gin.H{"Message": "Updated successfully", "Data updated": updatedData})
			return
		}
		ctx.JSON(500, gin.H{"Message": "Updated Failed"})
	})

	subjectGroup.DELETE("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		idObject, _ := primitive.ObjectIDFromHex(id)
		success := subjectController.DeleteSubject(idObject)
		if success == nil {
			ctx.JSON(200, gin.H{"Message": "Deleted Successfully"})
			return
		}
		ctx.JSON(500, gin.H{"Message": "Deleted Failed"})
	})

	subjectGroup.GET("/", func(ctx *gin.Context) {
		result := subjectController.GetAllSubject()
		if len(result) >= 1 {
			ctx.JSON(200, gin.H{"Message": "All data fetched Successfully", "Data": result})
			return
		}
		ctx.JSON(500, gin.H{"Message": "Not found"})
	})
}
