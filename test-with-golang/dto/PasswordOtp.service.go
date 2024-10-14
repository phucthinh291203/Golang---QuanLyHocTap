package services

import (
	"context"
	"errors"
	"log"
	"test-with-golang/Models"
	auth "test-with-golang/auth"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type OtpDTO struct {
	collection *mongo.Collection
}

func NewOtpDTO(collection *mongo.Collection) *OtpDTO {
	return &OtpDTO{
		collection: collection,
	}
}

func (service *OtpDTO) SendOTP(email string, role string, data Models.MyData) error {

	// Tìm kiếm OTP hiện tại dựa trên email
	var existingOTP Models.OTP
	err := service.collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&existingOTP)
	if err == nil {
		// Nếu tìm thấy OTP và nó vẫn còn hiệu lực và chưa xác thực
		if time.Now().Before(existingOTP.ExpiresAt) && !existingOTP.Verified {
			log.Println("OTP hiện tại vẫn còn hiệu lực, không tạo mã mới")

			//Cập nhật lại thời gian mới cho otp cũ
			newExpiresAt := time.Now().Add(10 * time.Minute)
			update := bson.M{
				"$set": bson.M{
					"expires_at": newExpiresAt,
				},
			}
			_, updateErr := service.collection.UpdateOne(context.TODO(), bson.M{"email": email}, update)

			if updateErr != nil {
				log.Print("Không thể cập nhật thời gian hết hạn OTP")
				return updateErr
			}

			go auth.SendMail(email, existingOTP.Code) // Gửi lại OTP hiện có
			return nil
		}
	}

	service.collection.DeleteMany(context.TODO(), bson.M{"email": email})
	otpCode := auth.GenerateOTP(6)
	go auth.SendMail(email, otpCode)

	expiresAt := time.Now().Add(10 * time.Minute)
	otp := Models.OTP{
		Email:     email,
		Role:      role,
		Code:      otpCode,
		Verified:  false,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}
	_, err = service.collection.InsertOne(context.TODO(), otp)
	if err != nil {
		log.Print("Không thể lưu OTP vào database")
		return err
	}

	return nil
}

func (service *OtpDTO) VerifyOTP(email string, otpInput string, data Models.MyData) error {

	var otpReceive Models.OTP
	filter := bson.M{
		"email":    email,
		"verified": true,
	}

	//Tìm email với otp
	err := service.collection.FindOne(context.TODO(), bson.M{"email": email, "code": otpInput}).Decode(&otpReceive)
	if err != nil {
		return errors.New("sai otp")
	}

	if otpReceive.Verified {
		return errors.New("đã xác thực rồi")
	}

	filter = bson.M{
		"email": email,
		"code":  otpInput,
	}
	result := service.collection.FindOneAndUpdate(
		context.TODO(), filter, bson.M{"$set": bson.M{"verified": true}})

	if result == nil {
		return errors.New("không thể xác thực OTP của bạn")
	}

	return nil
}

func (service *OtpDTO) ResetPassword(email string, newPassword string, data Models.MyData) error {

	//Check trong collection OTP có verified chưa
	var user Models.OTP
	err := data.OTPCollection.FindOne(context.TODO(), bson.M{"email": email, "verified": true}).Decode(&user)
	if err != nil {
		return errors.New("bạn chưa xác thực")
	}

	//Băm password mới
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	updateData := bson.M{
		"password": string(hashedPassword),
	}

	if user.Role == "Admin" {
		data.UserCollection.FindOneAndUpdate(
			context.TODO(),
			bson.M{"email": email},
			bson.M{"$set": updateData},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		)
	} else if user.Role == "Teacher" {
		data.TeacherCollection.FindOneAndUpdate(
			context.TODO(),
			bson.M{"email": email},
			bson.M{"$set": updateData},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		)
	} else {
		return errors.New("đã có lỗi xảy ra trong việc xác định role")
	}

	//Update xong thì xóa data trong OTP
	go data.OTPCollection.DeleteMany(context.TODO(), bson.M{"email": email})

	log.Println("Đã reset mật khẩu và xóa otp khỏi database")
	return nil
}

func (service *OtpDTO) SendSMS(phoneNumber string, role string, data Models.MyData) error {
	var existingSMS Models.SMS
	err := data.SMSCollection.FindOne(context.TODO(), bson.M{"phone_number": phoneNumber}).Decode(&existingSMS)
	if err == nil {
		// Nếu tìm thấy OTP và nó vẫn còn hiệu lực và chưa xác thực
		if time.Now().Before(existingSMS.ExpiresAt) && !existingSMS.Verified {
			log.Println("OTP hiện tại vẫn còn hiệu lực, không tạo mã mới")

			//Cập nhật lại thời gian mới cho otp cũ
			newExpiresAt := time.Now().Add(10 * time.Minute)
			update := bson.M{
				"$set": bson.M{
					"expires_at": newExpiresAt,
				},
			}
			_, updateErr := service.collection.UpdateOne(context.TODO(), bson.M{"phone_number": phoneNumber}, update)

			if updateErr != nil {
				log.Print("Không thể cập nhật thời gian hết hạn OTP")
				return updateErr
			}

			go auth.SendSMS(phoneNumber, existingSMS.Code) // Gửi lại OTP hiện có
			return nil
		}
	}

	data.SMSCollection.DeleteMany(context.TODO(), bson.M{"phone_number": phoneNumber})
	otpCode := auth.GenerateOTP(6)
	go auth.SendSMS(phoneNumber, otpCode) // Gửi lại OTP hiện có
	expiresAt := time.Now().Add(10 * time.Minute)
	otp := Models.SMS{
		PhoneNumber: phoneNumber,
		Role:        role,
		Code:        otpCode,
		Verified:    false,
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
	}
	_, err = data.SMSCollection.InsertOne(context.TODO(), otp)
	if err != nil {
		log.Print("Không thể lưu OTP vào database")
		return err
	}

	return nil
}

func (service *OtpDTO) VerifyOTP_SMS(phoneNumber string, otpInput string, data Models.MyData) error {

	var otpReceive Models.SMS
	filter := bson.M{
		"phone_number": phoneNumber,
		"verified":     true,
	}

	//Tìm email với otp
	err := data.SMSCollection.FindOne(context.TODO(), bson.M{"phone_number": phoneNumber, "code": otpInput}).Decode(&otpReceive)
	if err != nil {
		return errors.New("sai otp")
	}

	if otpReceive.Verified {
		return errors.New("đã xác thực rồi")
	}

	filter = bson.M{
		"phone_number": phoneNumber,
		"code":         otpInput,
	}
	result := data.SMSCollection.FindOneAndUpdate(
		context.TODO(), filter, bson.M{"$set": bson.M{"verified": true}})

	if result == nil {
		return errors.New("không thể xác thực OTP của bạn")
	}

	return nil
}

func (service *OtpDTO) ResetPassword_SMS(phoneNumber string, newPassword string, data Models.MyData) error {

	//Check trong collection SMS có verified chưa
	var user Models.SMS
	err := data.SMSCollection.FindOne(context.TODO(), bson.M{"phone_number": phoneNumber, "verified": true}).Decode(&user)
	if err != nil {
		return errors.New("bạn chưa xác thực")
	}

	//Băm password mới
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	updateData := bson.M{
		"password": string(hashedPassword),
	}

	if user.Role == "Admin" {
		data.UserCollection.FindOneAndUpdate(
			context.TODO(),
			bson.M{"phone_number": phoneNumber},
			bson.M{"$set": updateData},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		)
	} else if user.Role == "Teacher" {
		data.TeacherCollection.FindOneAndUpdate(
			context.TODO(),
			bson.M{"phone_number": phoneNumber},
			bson.M{"$set": updateData},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		)
	} else {
		return errors.New("đã có lỗi xảy ra trong việc xác định role")
	}

	//Update xong thì xóa data trong OTP
	go data.SMSCollection.DeleteMany(context.TODO(), bson.M{"phone_number": phoneNumber})

	log.Println("Đã reset mật khẩu và xóa otp khỏi database")
	return nil
}
