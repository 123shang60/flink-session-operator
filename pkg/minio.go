package pkg

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"k8s.io/klog/v2"
)

type MinioClient struct {
	*minio.Client
}

func AutoConnMinio(endpoints, accessKeyID, secretAccessKey string) (*MinioClient, error) {
	klog.Info("准备连接minio！", endpoints)
	minioClient, err := minio.New(endpoints, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		klog.Error("minio 连接异常!，", err)
		return nil, err
	}

	return &MinioClient{
		minioClient,
	}, nil
}

// TracFile 清理 minio 中 ha 相关的文件 ，从而达到 flink 可用状态
// 逻辑如下：
// 1、如果没有 bucket ，要创建
// 2、如果有 bucket，要删除指定路径下的 object ，然后创建
// 2.1、清除逻辑参考 mc 源码 (https://github.com/minio/mc/blob/RELEASE.2022-02-23T03-15-59Z/cmd/rb-main.go#L153) ，批量删除后清除空 bucket
func (m *MinioClient) TracFile(bucket, prefix string) error {
	klog.Info("准备清理 m，bucket : ", bucket)
	ctx := context.Background()
	exists, err := m.BucketExists(ctx, bucket)
	if err != nil {
		klog.Error("m 清理失败！", err)
		return err
	}
	if exists {
		klog.Info("bucket 已经存在，开始清理")
		if err := m.removeFile(ctx, bucket, prefix); err != nil {
			klog.Error("minio bucket 清理失败！", err)
			return err
		}
	} else {
		if err := m.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			klog.Error("bucket 创建失败！,", err)
			return err
		}
	}

	klog.Info("minio 清理成功!")
	return nil
}

// removeBucket 清理 bucket，无论是否为空
func (m *MinioClient) removeFile(ctx context.Context, bucket, prefix string) error {
	objectCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectCh)
		for object := range m.ListObjects(ctx, bucket, minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		}) {
			objectCh <- object
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	for err := range m.RemoveObjects(ctx, bucket, objectCh, opts) {
		if err.Err != nil {
			klog.Error("object 删除失败！", err)
			return err.Err
		}
	}

	klog.Info("minio 清理成功!")
	return nil
}
