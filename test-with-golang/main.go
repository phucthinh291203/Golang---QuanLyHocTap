package main

import (
	"log"
	"test-with-golang/database"
	"test-with-golang/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectToDatabase()
	data := database.GetData()
	database.CreateMinioClient()

	// Khởi động server
	server := gin.Default()

	//routes
	routes.AuthRoutes(server)
	routes.ClassRoutes(server, data)
	routes.StudentRoutes(server, data)
	routes.TeacherRoute(server, data)
	routes.ScoreRoute(server, data)
	routes.SubjectRoute(server, data)
	routes.BangDiemRoutes(server, data)
	routes.PasswordRecovery(server, data)
	routes.MinioRoute(server,data)
	go database.StartOTPCleaner(data)

	if err := server.Run(":8080"); err != nil {
		log.Print("Lỗi khi khởi động server:", err)
	}
}
