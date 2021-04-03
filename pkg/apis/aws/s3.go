package aws

import (
	"fmt"
)

type S3Config struct {
	Bucket string `json:"bucket"`
	Region string `json:"region"`
	Prefix string `json:"prefix,omitempty"`
}

func (config S3Config) Export() error {
	ak, sak, region := readCredentialsFromEnv()
	sess, err := connectAWS(ak, sak, region)
	if err != nil {
		return fmt.Errorf("error occurred while establishing a session with AWS: %v", err)
	}
	filePath := "/export/kube-bench/report.txt"
	if err := addFileToS3(sess, filePath, config); err != nil {
		return fmt.Errorf("error occurred while exporting the report to S3: %v", err)
	}
	return nil
}
