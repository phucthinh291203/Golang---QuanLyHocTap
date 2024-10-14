package middleware

import (
	"net/http"
	"strings"
	auth "test-with-golang/auth"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleWare(requiredTitles string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Lấy token từ header Authorization
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			ctx.Abort()
			return
		}

		// Xóa "Bearer " khỏi token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Giải mã và kiểm tra token
		claims, err := auth.ParseJWT(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Kiểm tra vai trò của người dùng
		title := claims.Title // sử dụng claims từ custom claims

		if title != requiredTitles {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Sai chức vụ"})
			ctx.Abort()
		}

		ctx.Set("username", claims.Username)
		ctx.Set("name", claims.Name)

		if claims.Title == "Teacher" {
			ctx.Set("teacher_id", claims.UserID)
		}

		ctx.Next()
	}

	// Lưu username vào context để sử dụng trong các handler
}
