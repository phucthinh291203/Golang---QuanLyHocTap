package database

import (
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func CreateMinioClient() {
	endpoint := "localhost:9000"
	accessKeyID := os.Getenv("Access_KeyID_Minio")
	secretAccessKey := os.Getenv("Secret_Accesskey_Minio")
	useSSL := false //Vì local host nên không cần (giúp bảo mật)

	// Initialize minio client object.
	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Khởi tạo minio Client thành công") // minioClient is now setup
		MinioClient = mc
	}

	log.Println("Access Key ID:", accessKeyID)
	log.Println("Secret Access Key:", secretAccessKey)
}

func GetMyClient() *minio.Client {
	return MinioClient
}
