package pkg

import (
	"context"
	"github.com/minio/minio-go/v7"
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
	conn, err := AutoConnMinio(EndPoint, AccessKeyID, SecretAccessKey)
	assert.Nil(t, err)

	err = conn.TracFile(TestBucketCreate, "")
	assert.Nil(t, err)
}

func TestCleanFile(t *testing.T) {
	conn, err := AutoConnMinio(EndPoint, AccessKeyID, SecretAccessKey)
	assert.Nil(t, err)

	err = conn.TracFile(TestBucketDelete, "")
	assert.Nil(t, err)

	// 上传一些文件
	reader := strings.NewReader("somethings test")

	_, err = conn.PutObject(
		context.Background(),
		TestBucketDelete,
		"deletes.txt",
		reader,
		reader.Size(),
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)

	assert.Nil(t, err)

	err = conn.TracFile(TestBucketDelete, "")
	assert.Nil(t, err)
}

func TestCleanFilePrefix(t *testing.T) {
	conn, err := AutoConnMinio(EndPoint, AccessKeyID, SecretAccessKey)
	assert.Nil(t, err)

	err = conn.TracFile(TestBucketDelete, "")
	assert.Nil(t, err)

	// 上传一些文件
	reader := strings.NewReader("somethings test")

	_, err = conn.PutObject(
		context.Background(),
		TestBucketDelete,
		"test/deletes.txt",
		reader,
		reader.Size(),
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	assert.Nil(t, err)

	err = conn.TracFile(TestBucketDelete, "test")
	assert.Nil(t, err)
}
