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

func ScoreRoute(router *gin.Engine, data Models.MyData) {
	scoreCollection := database.GetScoreCollection()
	scoreDTO := dto.NewScoreDTO(scoreCollection)
	scoreController := controllers.NewScore(*scoreDTO)

	scoreGroup := router.Group("teachers/scores")

	scoreGroup.Use(middleware.JWTAuthMiddleWare("Teacher"))
	scoreGroup.GET("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		result := scoreController.FindOne(id, data)
		if result.ID != primitive.NilObjectID {
			ctx.JSON(200, gin.H{"Message": "Tim kiem thanh cong", "ScoreData": result})
			return
		}

		ctx.JSON(500, gin.H{"Message": "Tim kiem that bai"})
	})

	scoreGroup.POST("/", func(ctx *gin.Context) {
		// Lấy teacher_id từ context
		teacherId, exists := ctx.Get("teacher_id")
		if !exists {
			ctx.JSON(400, gin.H{"Message": "Teacher ID is missing"})
			return
		}

		teacherObjectID, ok := teacherId.(primitive.ObjectID)
		if !ok {
			ctx.JSON(400, gin.H{"Message": "Invalid Teacher ID"})
			return
		}

		// Lấy dữ liệu điểm từ request body
		var scoreData Models.Score
		if err := ctx.ShouldBindJSON(&scoreData); err != nil {
			ctx.JSON(400, gin.H{"Message": "Invalid score data", "error": err.Error()})
			return
		}

		// Gọi controller để tạo điểm cho học sinh
		err := scoreController.CreateScoreForStudent(teacherObjectID, scoreData, data)
		if err != nil {
			ctx.JSON(500, gin.H{"Message": "Cham diem that bai", "error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"Message": "Cham diem thanh cong", "result": scoreData})
	})

	scoreGroup.PUT("/:id", func(ctx *gin.Context) {
		teacherId, exists := ctx.Get("teacher_id")
		if !exists {
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		// Chuyển đổi teacherId sang kiểu ObjectID
		teacherObjectID, ok := teacherId.(primitive.ObjectID)
		if !ok {
			ctx.JSON(400, gin.H{"error": "Invalid teacher ID"})
			return
		}

		// Lấy scoreID từ params
		scoreID := ctx.Param("id")
		if len(scoreID) == 0 {
			ctx.JSON(400, gin.H{"error": "Score ID is required"})
			return
		}

		// Bind dữ liệu JSON từ client vào struct Score
		var scoreData Models.Score
		if err := ctx.ShouldBindJSON(&scoreData); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}


		// Gọi hàm controller để cập nhật điểm cho học sinh
		err := scoreController.UpdateScore(teacherObjectID, scoreID, scoreData, data)
		if err != nil{
			ctx.JSON(500,gin.H{"Message":"Update failed","Error":err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "Score updated successfully", "new data": scoreData})
	})

	scoreGroup.DELETE("/:id", func(ctx *gin.Context) {
		teacherId, exists := ctx.Get("teacher_id")
		if !exists {
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		// Chuyển đổi teacherId sang kiểu ObjectID
		teacherObjectID, ok := teacherId.(primitive.ObjectID)
		if !ok {
			ctx.JSON(400, gin.H{"error": "Invalid teacher ID"})
			return
		}

		// Lấy scoreID từ params
		scoreID := ctx.Param("id")
		if len(scoreID) == 0 {
			ctx.JSON(400, gin.H{"error": "Score ID is required"})
			return
		}

		err := scoreController.DeleteScore(teacherObjectID, scoreID, data)
		if err == nil {
			ctx.JSON(200, gin.H{"message": "Score Deleted successfully"})
		}
	})

	scoreGroup.GET("/HaveCreated", func(ctx *gin.Context) {
		teacherId, exists := ctx.Get("teacher_id")
		if !exists {
			ctx.JSON(400, gin.H{"Message": "Teacher ID is missing"})
			return
		}

		teacherObjectId, _ := teacherId.(primitive.ObjectID)
		result := scoreController.FindAll(teacherObjectId, data)
		if len(result) >= 1 {
			ctx.JSON(200, gin.H{"Message": "All Score Fetch Successfully", "All Score": result})
			return
		}
		ctx.JSON(500, gin.H{"Message": "Not found"})
	})
}
