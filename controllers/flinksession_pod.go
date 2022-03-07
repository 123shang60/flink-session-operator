package controllers

import (
	"context"
	"errors"
	flinkv1 "github.com/123shang60/flink-session-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func (r *FlinkSessionReconciler) waitDeletePod(session *flinkv1.FlinkSession) error {
	for i := 0; i < 30; i++ {
		var pods corev1.PodList

		if err := r.List(context.Background(), &pods, client.MatchingLabels{
			"app":  session.GetName(),
			"type": FlinkNativeType,
		}, client.InNamespace(session.GetNamespace())); err != nil || len(pods.Items) == 0 {
			klog.Error("获取 flink pods 相关列表失败或者pod不存在!", err)
			return err
		} else {
			klog.Info(session.GetName(), ",", session.GetNamespace(), " : 获取 pods 列表:", pods)
		}
		time.Sleep(time.Second * 10)
	}
	return errors.New("flink pods 等待删除完成超时！")
}
