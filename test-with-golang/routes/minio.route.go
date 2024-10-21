package routes

import (
	"net/http"
	"strings"
	controllers "test-with-golang/Controllers"
	"test-with-golang/Models"
	"test-with-golang/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MinioRoute(router *gin.Engine, data Models.MyData) {

	minioGroup := router.Group("/minio")
	minioGroup.POST("/newBucket", func(ctx *gin.Context) {
		controllers.CreateBucket()
	})

	minioGroup.POST("/uploadImage", func(ctx *gin.Context) {
		url := controllers.UploadImageToBucket_2()
		ctx.JSON(http.StatusOK, url)
	})

	minioGroup.GET("/checkBucket", func(ctx *gin.Context) {
		controllers.CheckCredentials()
	})

	minioGroup.GET("/admin/getAllHistoryUploaded", func(ctx *gin.Context) {
		allHistory := controllers.GetAllHistoryUploaded(data)
		ctx.JSON(http.StatusOK, Response{
			Message: "All data",
			Data:    allHistory,
		})
	})

	minioGroup.GET("/admin/getHistoryUploaded/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		idObject, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Nhập sai id",
			})
			return
		}

		history := controllers.GetHistoryUploadedById(idObject, data)
		if history.ID != primitive.NilObjectID {
			ctx.JSON(http.StatusOK, Response{
				Message: "Tìm thấy data",
				Data:    history,
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, Response{
			Message: "Không có kết quả",
			Data:    nil,
		})

	})

	minioGroup.GET("/admin/getAllFile", func(ctx *gin.Context) {
		allFile := controllers.GetAllFile(data)
		ctx.JSON(http.StatusOK, Response{
			Message: "All data",
			Data:    allFile,
		})
	})

	minioGroup.GET("/admin/getFile/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		idObject, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Nhập sai id",
			})
			return
		}

		file := controllers.GetFileById(idObject, data)
		if file.ID != primitive.NilObjectID {
			ctx.JSON(http.StatusOK, Response{
				Message: "Tìm thấy data",
				Data:    file,
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, Response{
			Message: "Không có kết quả",
			Data:    nil,
		})

	})

	minioGroup.Use(middleware.JWTAuthMiddleWare("Teacher"))
	minioGroup.GET("/download/:id", func(ctx *gin.Context) {
		fileId := ctx.Param("id")
		if fileId == "" {
			ctx.JSON(http.StatusBadRequest, Response{
				Message: "Chưa nhập Id file",
			})
		}
		url := controllers.DownloadFile(ctx, fileId, data)
		if url != "" {
			ctx.JSON(http.StatusOK, Response{
				Message: url,
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, Response{
			Message: "Đã có lỗi xảy ra",
		})
	})

	minioGroup.POST("/teacher/importListStudent", func(ctx *gin.Context) {
		//Lấy thông tin token
		teacherId, _ := ctx.Get("teacher_id")
		teacherObjectId, _ := teacherId.(primitive.ObjectID)
		if teacherObjectId == primitive.NilObjectID {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "không lấy được teacherId"})
			return
		}

		// Check xem có file không
		file, header, err := ctx.Request.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "File chưa được đính kèm"})
			return
		}
		defer file.Close()

		// Lấy tên file vừa upload
		fileName := header.Filename
		var fileType string

		//Thêm dữ liệu từ json vào collection học sinh
		if strings.HasSuffix(fileName, ".csv") {
			controllers.HandleCSV(file, data)
			fileType = ".csv"
		} else if strings.HasSuffix(fileName, ".xlsx") {
			controllers.HandleXLSX(file, data)
			fileType = ".xlsx"
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Định dạng không hợp lệ"})
			return
		}

		//Tạo thông tin file đó trong mongo
		fileID, err := controllers.CreateFile(fileName, data, fileType)
		if fileID != primitive.NilObjectID && err == nil {
			//Tạo history Uploaded trong mongo
			err = controllers.CreateHistoryUploaded(teacherObjectId, fileID, data)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Không khởi tạo lịch sử uploaded được"})
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Không khởi tạo file hoặc history uploaded được"})
			return
		}

		//Up nội dung json lên S3 theo id
		err = controllers.UpFileToMinio(file, fileID, fileType)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Không Up nội dung file lên S3 được"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Students imported successfully"})
	})

	minioGroup.GET("/admin/getHistoryUploaded", func(ctx *gin.Context) {
		allDownloaded := controllers.GetHistoryDownloaded(data)
		ctx.JSON(http.StatusOK, Response{
			Message: "All data",
			Data:    allDownloaded,
		})
	})

	minioForAdmin := router.Group("/minio")
	minioForAdmin.Use(middleware.JWTAuthMiddleWare("Admin"))

}
