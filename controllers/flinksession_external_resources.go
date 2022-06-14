package controllers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	flinkv1 "github.com/123shang60/flink-session-operator/api/v1"
	"github.com/123shang60/flink-session-operator/pkg"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *FlinkSessionReconciler) cleanExternalResources(f *flinkv1.FlinkSession) error {
	// 清理 ha
	if f.Spec.HA.Typ == flinkv1.ZKHA {
		klog.Info("zk ha 模式，开始清理！")
		zkCli, err := func() (*pkg.ZooKeeper, error) {
			if f.Spec.Security.Kerberos == nil || !pkg.SplitContainer(f.Spec.Security.Kerberos.Contexts, ",", "Client") {
				return pkg.AutoConnZk(f.Spec.HA.Quorum, nil)
			} else {
				keytab, err := base64.StdEncoding.DecodeString(f.Spec.Security.Kerberos.Base64Keytab)
				if err != nil {
					klog.Error("kerberos 解析失败！", err)
					return nil, err
				}
				return pkg.AutoConnZk(f.Spec.HA.Quorum, &pkg.KerberosConfig{
					Keytab:       keytab,
					Krb5:         f.Spec.Security.Kerberos.Krb5,
					PrincipalStr: f.Spec.Security.Kerberos.Principal,
				})
			}
		}()
		if err != nil {
			klog.Error("zk 不可用，请检查 zk 配置！", err)
			return errors.New("zk 不可用，请检查 zk 配置！")
		}
		defer zkCli.Close()
		if f.Spec.HA.Path != "" {
			err := zkCli.AutoDelete(f.Spec.HA.Path)
			if err != nil {
				klog.Error("zk 删除失败，忽略错误", err)
			} else {
				klog.Info("zk 指定路径 清理完成！")
			}
		} else {
			// 默认路径为 /flink/${high-availability.cluster-id}
			klog.Info("清理默认路径！：/flink/" + f.Name)
			err := zkCli.AutoDelete(fmt.Sprintf("/flink/%s", f.Name))
			if err != nil {
				klog.Error("zk 删除失败，忽略错误", err)
			} else {
				klog.Info("zk 默认路径 清理完成！")
			}
		}
	} else if f.Spec.HA.Typ == flinkv1.CONFIGMAPHA {
		klog.Info("k8s native ha 模式，开始清理！")
		var haConfigMaps corev1.ConfigMapList
		if err := r.List(context.Background(), &haConfigMaps, client.MatchingLabels{
			"app":            f.GetName(),
			"configmap-type": FlinkHAConfigType,
			"type":           FlinkNativeType,
		},
			client.InNamespace(f.GetNamespace()),
		); err != nil {
			klog.Error("获取 ha 相关 configmap 列表失败！", err)
			return err
		}

		for _, configMap := range haConfigMaps.Items {
			klog.Info("开始清理 ha configmap ，uuid：", configMap.UID)
			if err := r.Delete(context.Background(), &configMap); err != nil {
				klog.Error("删除 configmap 失败！", err)
			}
		}
		klog.Info("k8s native ha 模式，清理完毕！")
	} else {
		klog.Info("未开启 ha ，无需清理！")
	}
	// 初始化 minio
	klog.Info("开始清理 minio")
	minioClient, err := pkg.AutoConnMinio(f.Spec.S3.EndPoint, f.Spec.S3.AccessKey, f.Spec.S3.SecretKey)
	if err != nil {
		klog.Error("minio 不可用，请检查 minio 配置！", err)
		return errors.New("minio 不可用，请检查 minio 配置！")
	}
	err = minioClient.TracFile(f.Spec.S3.Bucket, fmt.Sprintf("%s/flink/ha/metadata", f.Name))
	if err != nil {
		klog.Error("minio 初始化 bucket 失败！，忽略", err)
	} else {
		klog.Info("minio 清理成功！")
	}
	return nil
}
