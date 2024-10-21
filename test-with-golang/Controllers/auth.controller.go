package controllers

import (
	"net/http"
	"test-with-golang/Models"
	auth "test-with-golang/auth"
	database "test-with-golang/database"
	dto "test-with-golang/dto"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Không dùng nữa
func Register(c *gin.Context) {
	var user Models.Admin
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Kiểm tra xem username đã tồn tại trong bảng Teacher hay User
	existingUser, _ := database.FindUserByUsername(user.Username)
	existingTeacher, _ := database.FindTeacherByUsername(user.Username)

	if existingUser.Username != "" || existingTeacher.Username != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	//Ko nhập role
	if user.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chưa nhập role"})
		return
	}

	//Mã hóa mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if user.Title == "Teacher" {
		teacherCollection := database.GetTeacherCollection()
		TeacherDTO := dto.NewTeacher(teacherCollection)
		teacherController := NewTeacher(TeacherDTO)

		teacher := Models.Teacher{
			Username:    user.Username,
			Password:    string(hashedPassword), // Có thể mã hóa password nếu cần
			TeacherName: user.Name,
			Email:       user.Email,
		}

		teacherErr := teacherController.Create(teacher, database.GetData()) // Gọi phương thức create ở đây
		if teacherErr != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Tạo thất bại"})
			return
		}
	}

	// if user.Title == "Admin" {
	// 	//Lưu mật khẩu mã hóa vào collection User
	// }

	user.Password = string(hashedPassword)

	if err := database.SaveUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var input Models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if input.Type == "Admin" {
		Login_Admin(c, input)
	} else if input.Type == "Teacher" {
		Login_Teacher(c, input)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Chưa nhập type"})
		return
	}
}

func Login_Admin(c *gin.Context, input Models.LoginInput) {
	//Tìm xem username có tồn tại ko
	user, err := database.FindUserByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	//Giải mã password về lại ban đầu
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	claims := auth.BaseClaims{
		UserID:   user.ID,
		Username: user.Username,
		Name:     user.Name,
		Title:    "Admin",
	}

	//Tạo JWT cho người dùng
	token, err := auth.GenerateJWT(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
	c.JSON(200, gin.H{"Hello: ": user.Name})
}

func Login_Teacher(c *gin.Context, input Models.LoginInput) {
	//Tìm xem username có tồn tại ko
	teacher, err := database.FindTeacherByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	//Giải mã password về lại ban đầu
	if err := bcrypt.CompareHashAndPassword([]byte(teacher.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	claims := auth.BaseClaims{
		Title:    "Teacher",
		Username: teacher.Username,
		Name:     teacher.TeacherName,
		UserID:   teacher.ID,
	}

	//Tạo JWT cho người dùng
	token, err := auth.GenerateJWT(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
	c.JSON(200, gin.H{"Hello: ": teacher.TeacherName})
}
