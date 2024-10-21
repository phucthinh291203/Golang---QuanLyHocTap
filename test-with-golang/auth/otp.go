package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"log"
	"os"
	"test-with-golang/Models"
	"time"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	LogoURL string
	Name    string
	OTP     string
}

func SendMail(receiverEmail string, otp string) {

	// Chuẩn bị dữ liệu cho email
	emailData := EmailData{
		LogoURL: "https://media2.giphy.com/media/v1.Y2lkPTc5MGI3NjExdDY3ZHh6Z3h6M3Fvd25iZGxmZTRmZmRsMWFoejVwM2xxYXU0NmUwayZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/K46xHSeu0XS8ADiP6u/giphy.webp", // URL logo
		OTP:     otp,                                                                                                                                                                              // Mã OTP
	}

	tmpl, err := template.ParseFiles("templates/email.template.html")
	if err != nil {
		log.Fatalf("Không thể đọc template: %v", err)
	}

	var emailHTML bytes.Buffer
	err = tmpl.Execute(&emailHTML, emailData)
	if err != nil {
		log.Fatalf("Không thể thực thi template: %v", err)
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", os.Getenv("password_recovery_username"))
	mail.SetHeader("To", receiverEmail)
	mail.SetHeader("Subject", "Gửi OTP")
	mail.SetBody("text/html", emailHTML.String())

	sendMail := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("password_recovery_username"), os.Getenv("password_recovery_password"))
	if err := sendMail.DialAndSend(mail); err == nil {
		log.Println("Đã gửi mail thành công")
	}
}

// Hàm gửi SMS qua Twilio
func SendSMS(phoneNumber string, otp string) {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("Account_SID"),
		Password: os.Getenv("Auth_Token_Twilio"),
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(os.Getenv("Sender_Number"))
	params.SetBody("Xin chào, mã otp xác nhận của bạn là: " + otp)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Println("Lỗi trong quá trình gửi SMS: " + err.Error())
	} else {
		log.Println("Đã gửi mã OTP")
	}
}

func GenerateOTP(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)[:length]
}

func RemoveExpiredOTPs(data Models.MyData) {
	filter := bson.M{
		"expires_at": bson.M{
			"$lt": time.Now(),
		},
		"verified": false,
	}
	_, err := data.OTPCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Println("Lỗi khi xóa OTP đã hết hạn:", err)
	} else {

	}
}
