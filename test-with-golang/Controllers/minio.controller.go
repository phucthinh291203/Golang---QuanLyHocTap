package controllers

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"test-with-golang/Models"
	database "test-with-golang/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateBucket() {
	minioClient := database.GetMyClient()

	bucketName := "hinhanh"
	err := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		log.Panicln(err)
		return
	}
	found, _ := minioClient.BucketExists(context.Background(), bucketName)
	if found {
		log.Println("Bucket tồn tại")
	} else {
		log.Panicln("Không thấy bucket vừa tạo")
	}
}

func UploadImageToBucket() {
	minioClient := database.GetMyClient()
	bucketName := "hinhanh"
	objectName := "test-image.jpg"                                                             // Tên file khi lưu trên MinIO
	fileName := "C:\\Users\\GIGABYTE\\Desktop\\FIle_hoc_tap\\golang\\my-image\\test-image.jpg" // Đường dẫn của cái file mình up lên

	ui, err := minioClient.FPutObject(context.Background(), bucketName, objectName, fileName,
		minio.PutObjectOptions{ContentType: "image/jpg"})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Uploaded", ui.Key, "To", ui.Bucket, ui.ETag, ui.VersionID, ui.Size)
}

func UploadImageToBucket_2() string {
	minioClient := database.GetMyClient()
	bucketName := "photos"
	found, _ := minioClient.BucketExists(context.Background(), bucketName)
	if found {
		log.Println("Bucket tồn tại")
	} else {
		log.Panicln("Không thấy bucket vừa tạo")
	}
	// Tạo Presigned URL cho uploads
	url, err := minioClient.PresignedPutObject(context.Background(), bucketName, "file_khach_up.jpg", time.Hour*24)
	if err != nil {
		log.Println(err)
		return url.String()
	}
	log.Println("Presigned URL:", url)
	return url.String()
}

func CheckCredentials() {
	minioClient := database.GetMyClient()
	// Lấy danh sách bucket
	buckets, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		log.Println("Không thể xác thực với MinIO:", err)
		return
	}

	log.Println("Danh sách bucket:")
	for _, bucket := range buckets {
		log.Println(bucket.Name)
	}
}

func PutFileToCollection(file multipart.File, data Models.MyData) error {
	var students []Models.Student
	err := json.NewDecoder(file).Decode(&students)
	if err != nil {
		return err
	}

	_, err = data.StudentCollection.InsertMany(context.TODO(), toBsonArray(students))
	if err != nil {
		log.Println("Failed to insert records")
		return err
	}

	return nil
}

// hàm này sẽ loại bỏ ký tự BOM nếu có
func removeBOM(header []string) []string {
	for i, field := range header {
		if strings.HasPrefix(field, "\ufeff") {
			header[i] = strings.TrimPrefix(field, "\ufeff")
		}
	}
	return header
}

func HandleCSV(file io.Reader, data Models.MyData) error {
	reader := csv.NewReader(file)

	// Đọc tiêu đề (header) để biết thứ tự các trường
	header, err := reader.Read()
	if err != nil {
		return err
	}

	// Xóa BOM nếu có
	header = removeBOM(header)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		student, err := MapToStudent(header, record)
		if err != nil {
			return err
		}

		//Thêm dữ liệu trong file vào student collection
		data.StudentCollection.InsertOne(context.TODO(), student)

	}
	return nil
}

func HandleXLSX(file multipart.File, data Models.MyData) error {
	f, err := excelize.OpenReader(file) // Sử dụng OpenReader để mở file
	if err != nil {
		return err
	}

	// Lấy danh sách tên sheet
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		log.Println("no sheets found in the file")
		return err
	}

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		log.Println(err)
		return err
	}

	for i, row := range rows {
		if i == 0 { // Bỏ qua hàng đầu tiên (tiêu đề)
			continue
		}

		dateOfBirth, _ := time.Parse("1/2/2006", row[1])
		enrollmentDate, _ := time.Parse("1/2/2006", row[5])

		student := Models.Student{
			Name:           row[0],
			DateOfBirth:    dateOfBirth,
			Email:          row[2],
			PhoneNumber:    row[3],
			Address:        row[4],
			EnrollmentDate: enrollmentDate,
			Gender:         row[6],
			Nationality:    row[7],
			Avatar:         row[8],
		}
		log.Println(student)

		data.StudentCollection.InsertOne(context.TODO(), student)
	}
	return nil
}

func MapToStudent(header []string, record []string) (Models.Student, error) {
	if len(header) != len(record) {
		return Models.Student{}, errors.New("dữ liệu không khớp với tiêu đề")
	}

	var student Models.Student

	for i, field := range header {
		switch field {
		case "name":
			student.Name = record[i]
		case "date_of_birth":
			student.DateOfBirth, _ = time.Parse("2006-01-02", record[i])
		case "email":
			student.Email = record[i]
		case "phone_number":
			student.PhoneNumber = record[i]
		case "address":
			student.Address = record[i]
		case "enrollment_date":
			student.EnrollmentDate, _ = time.Parse("2006-01-02", record[i])
		case "gender":
			student.Gender = record[i]
		case "nationality":
			student.Nationality = record[i]
		case "avatar":
			student.Avatar = record[i]
		default:
			return Models.Student{}, errors.New("Trường dữ liệu không hợp lệ: " + field)
		}
	}
	return student, nil
}

func toBsonArray(students []Models.Student) []interface{} {
	var result []interface{}
	for _, student := range students {
		result = append(result, bson.M{
			"name":            student.Name,
			"email":           student.Email,
			"phone_number":    student.PhoneNumber,
			"address":         student.Address,
			"gender":          student.Gender,
			"nationality":     student.Nationality,
			"avatar":          student.Avatar,
			"date_of_birth":   student.DateOfBirth,
			"enrollment_date": student.EnrollmentDate,
			"class_id":        student.ClassID,
		})
	}
	return result
}

func UpFileToMinio(file multipart.File, fileObjectId primitive.ObjectID, fileType string) error {
	if file == nil {
		return errors.New("file rỗng")
	}
	minioClient := database.GetMyClient()
	bucketName := "filecontent"
	fileId := fileObjectId.Hex()
	objectName := fileId + fileType

	//Check nếu bucket tồn tại
	found, _ := minioClient.BucketExists(context.Background(), bucketName)
	if found {
		log.Println("Bucket tồn tại")
	} else {
		log.Panicln("Không thấy bucket vừa tạo")
	}

	// Tạo Presigned URL cho uploads
	url, err := minioClient.PresignedPutObject(context.Background(), bucketName, objectName, time.Minute*2)
	if err != nil {
		log.Println(err)
	}

	// Seek về đầu file để chắc chắn đọc từ đầu
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return errors.New("không thể reset con trỏ file")
	}

	log.Println(file)
	//viết hàm để đẩy file json vào url trên
	filecontent, err := io.ReadAll(file)
	if err != nil {
		return errors.New("có lỗi trong việc đọc nội dung file")
	}
	log.Printf("Nội dung file: %s", string(filecontent))

	req, err := http.NewRequest(http.MethodPut, url.String(), bytes.NewReader(filecontent))

	if err != nil {
		return errors.New("tạo method thất bại")
	}
	req.Header.Set("Content-Type", "application/json") // Thiết lập Content-Type cho request

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("thực hiện method thất bại")
	}
	defer resp.Body.Close()

	return nil

}

// func uploadUsingPresignedURL(url string, fileContent []byte) error {
// 	// Sử dụng http.NewRequest để upload file
// 	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(fileContent))
// 	if err != nil {
// 		return fmt.Errorf("failed to create request: %v", err)
// 	}
// 	req.Header.Set("Content-Type", "text/plain") // Thiết lập Content-Type cho request

// 	// Gửi yêu cầu
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("failed to upload file: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("failed to upload file, status code: %d", resp.StatusCode)
// 	}

// 	return nil
// }


func DownloadFile(ctx *gin.Context, fileId string, data Models.MyData) string {

	//Lấy token người dùng
	tokenId, _ := ctx.Get("teacher_id")
	tokenObjectId := tokenId.(primitive.ObjectID)

	minioClient := database.GetMyClient()
	bucketName := "filecontent"

	fileIdObject, _ := primitive.ObjectIDFromHex(fileId)
	var currentFile Models.File
	err := data.FileCollection.FindOne(context.TODO(), bson.M{"_id": fileIdObject}).Decode(&currentFile)
	if err != nil {
		return ""
	}
	objectName := fileId + currentFile.FileType

	if objectName == "" {
		log.Panicln("objectName is required")
		return ""
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\"tai_xuong\""+currentFile.FileType)

	switch currentFile.FileType {
	case ".csv":
		reqParams.Set("response-content-type", "application/csv")
	case ".xlsx":
		reqParams.Set("response-content-type", "application/xlsx")
	default:
		reqParams.Set("response-content-type", "application/txt")
	}

	// Tạo Presigned URL
	url, err := minioClient.PresignedGetObject(
		context.Background(),
		bucketName,
		objectName,
		time.Hour,
		reqParams,
	)
	if err != nil {
		log.Println("file not found")
		return ""
	}

	CreateHistoryDownloaded(tokenObjectId, fileIdObject, data)
	return url.String() //Trả về url dùng để get -> client lấy link về send & download
}
