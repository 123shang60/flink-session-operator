package pkg

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"k8s.io/klog/v2"
)

type S3Client struct {
	*s3.Client
}

func AutoConnS3(endpoints, accessKeyID, secretAccessKey string) (*S3Client, error) {
	klog.Info("准备连接s3！", endpoints)

	staticResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               fmt.Sprintf("http://%s", endpoints),
			SigningRegion:     region,
			HostnameImmutable: true,
		}, nil

	})

	s3Config, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithEndpointResolverWithOptions(staticResolver),
	)

	if err != nil {
		klog.Error("s3 连接异常!，", err)
		return nil, err
	}

	s3Client := s3.NewFromConfig(s3Config)

	return &S3Client{
		s3Client,
	}, nil
}

func (m *S3Client) BucketExists(bucket string) (bool, error) {
	buckets, err := m.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		klog.Error("获取 bucket 列表失败！")
		return false, err
	}
	for _, v := range buckets.Buckets {
		if fmt.Sprintf("%s", *v.Name) == fmt.Sprintf("%s", *aws.String(bucket)) {
			return true, nil
		}
	}
	return false, nil
}

// TracFile 清理 s3 中 ha 相关的文件 ，从而达到 flink 可用状态
// 逻辑如下：
// 1、如果没有 bucket ，要创建
// 2、如果有 bucket，要删除指定路径下的 object ，然后创建
// 2.1、清除逻辑参考 mc 源码 (https://github.com/minio/mc/blob/RELEASE.2022-02-23T03-15-59Z/cmd/rb-main.go#L153) ，批量删除后清除空 bucket
func (m *S3Client) TracFile(bucket, prefix string) error {
	klog.Info("准备清理 m，bucket : ", bucket)
	ctx := context.Background()
	exists, err := m.BucketExists(bucket)
	if err != nil {
		klog.Error("m 清理失败！", err)
		return err
	}
	if exists {
		klog.Info("bucket 已经存在，开始清理")
		if err := m.removeFile(ctx, bucket, prefix); err != nil {
			klog.Error("s3 bucket 清理失败！", err)
			return err
		}
	} else {
		if _, err := m.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		}); err != nil {
			klog.Error("bucket 创建失败！,", err)
			return err
		}
	}

	klog.Info("s3 清理成功!")
	return nil
}

// removeBucket 清理 bucket，无论是否为空
func (m *S3Client) removeFile(ctx context.Context, bucket, prefix string) error {
	listObjects, err := m.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		klog.Error("获取全量 s3 文件信息失败！", err)
		return err
	}

	if listObjects.Contents == nil || len(listObjects.Contents) == 0 {
		return nil
	}

	objects := make([]types.ObjectIdentifier, len(listObjects.Contents))

	for k, v := range listObjects.Contents {
		objects[k] = types.ObjectIdentifier{
			Key: v.Key,
		}
	}

	_, err = m.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &types.Delete{
			Objects: objects,
		},
	})
	if err != nil {
		klog.Error("删除失败！", err)
		return err
	}

	klog.Info("s3 清理成功!")
	return nil
}
