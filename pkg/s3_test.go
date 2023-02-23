package pkg

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const (
	EndPoint         string = "127.0.0.1:9000"
	AccessKeyID      string = "minioadmin"
	SecretAccessKey  string = "minioadmin"
	TestBucketCreate string = "testcreate"
	TestBucketDelete string = "testdelete"
)

func TestCreateBucket(t *testing.T) {
	conn, err := AutoConnS3(EndPoint, AccessKeyID, SecretAccessKey)
	assert.Nil(t, err)

	err = conn.TracFile(TestBucketCreate, "")
	assert.Nil(t, err)
}

func TestCleanFile(t *testing.T) {
	conn, err := AutoConnS3(EndPoint, AccessKeyID, SecretAccessKey)
	assert.Nil(t, err)

	err = conn.TracFile(TestBucketDelete, "")
	assert.Nil(t, err)

	// 上传一些文件
	reader := strings.NewReader("somethings test")

	//_, err = conn.PutObject(
	//	context.Background(),
	//	TestBucketDelete,
	//	"deletes.txt",
	//	reader,
	//	reader.Size(),
	//	minio.PutObjectOptions{ContentType: "application/octet-stream"},
	//)
	_, err = conn.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(TestBucketDelete),
		Key:         aws.String("deletes.txt"),
		Body:        reader,
		ContentType: aws.String("application/octet-stream"),
	})

	assert.Nil(t, err)

	err = conn.TracFile(TestBucketDelete, "")
	assert.Nil(t, err)
}

func TestCleanFilePrefix(t *testing.T) {
	conn, err := AutoConnS3(EndPoint, AccessKeyID, SecretAccessKey)
	assert.Nil(t, err)

	err = conn.TracFile(TestBucketDelete, "")
	assert.Nil(t, err)

	// 上传一些文件
	reader := strings.NewReader("somethings test")

	//_, err = conn.PutObject(
	//	context.Background(),
	//	TestBucketDelete,
	//	"test/deletes.txt",
	//	reader,
	//	reader.Size(),
	//	minio.PutObjectOptions{ContentType: "application/octet-stream"},
	//)
	_, err = conn.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(TestBucketDelete),
		Key:         aws.String("test/deletes.txt"),
		Body:        reader,
		ContentType: aws.String("application/octet-stream"),
	})
	assert.Nil(t, err)

	err = conn.TracFile(TestBucketDelete, "test")
	assert.Nil(t, err)
}
