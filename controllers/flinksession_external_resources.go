package controllers

import (
	"errors"
	flinkv1 "github.com/123shang60/flink-session-operator/api/v1"
	"github.com/123shang60/flink-session-operator/pkg"
	"k8s.io/klog/v2"
)

func (r *FlinkSessionReconciler) cleanExternalResources(f *flinkv1.FlinkSession) error {
	// 清理 zk
	if f.Spec.HA.Typ == flinkv1.ZKHA {
		klog.Info("zk ha 模式，开始清理！")
		zkCli, err := pkg.AutoConnZk(f.Spec.HA.Quorum)
		defer zkCli.Close()
		if err != nil {
			klog.Error("zk 不可用，请检查 zk 配置！", err)
			return errors.New("zk 不可用，请检查 zk 配置！")
		}
		if f.Spec.HA.Path != "" {
			err := zkCli.AutoDelete(f.Spec.HA.Path)
			if err != nil {
				klog.Error("zk 删除失败，忽略错误", err)
			} else {
				klog.Info("zk 指定路径 清理完成！")
			}
		} else {
			// TODO: 默认路径清理
		}
	}
	// 初始化 minio
	klog.Info("开始清理 minio")
	minioClient, err := pkg.AutoConnMinio(f.Spec.S3.EndPoint, f.Spec.S3.AccessKey, f.Spec.S3.SecretKey)
	if err != nil {
		klog.Error("minio 不可用，请检查 minio 配置！", err)
		return errors.New("minio 不可用，请检查 minio 配置！")
	}
	err = minioClient.TracBucket(f.Spec.S3.Bucket)
	if err != nil {
		klog.Error("minio 初始化 bucket 失败！，忽略", err)
	} else {
		klog.Info("minio 清理成功！")
	}
	return nil
}
