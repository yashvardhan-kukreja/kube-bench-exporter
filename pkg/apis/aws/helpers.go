package aws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/yashvardhan-kukreja/kube-bench-exporter/pkg/global"
	"net/http"
	"os"
	"strings"
	"time"
)

func readCredentialsFromEnv() (string, string, string) {
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("REGION")
	if region == "" {
		region = "us-east-1"
	}
	return accessKey, secretAccessKey, region
}

func DeserializeInputJsonToS3Config(input map[string]interface{}) (global.Target, error) {
	jsonByte, err := json.Marshal(input)
	var config S3Config
	if err = json.Unmarshal(jsonByte, &config); err != nil {
		return S3Config{}, err
	}
	return config, nil
}

func connectAWS(accessKey, secretAccessKey, region string) (*session.Session, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(accessKey, secretAccessKey, ""),
		},
	)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func addFileToS3(s *session.Session, fileDir string, config S3Config) error {
	bucketName, prefix := config.Bucket, config.Prefix

	file, err := os.Open(fileDir)
	if err != nil {
		return fmt.Errorf("error occurred while opening the kube-bench logs file: %v", err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	objectName := fmt.Sprintf("%s-kube-bench-report.txt", time.Now().Format("02-01-2006"))
	if prefix != "" {
		prefix = strings.TrimSuffix(prefix, "/ ") // trimming slashes and spaces in the end
		objectName = prefix + "/" + objectName
	}

	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(bucketName),
		Key:                aws.String(objectName),
		Body:               bytes.NewReader(buffer),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
	})
	if err != nil {
		return fmt.Errorf("error occurred while uploading the report to the S3 bucket: %v", err)
	}
	return nil
}
