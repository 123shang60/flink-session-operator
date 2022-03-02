package controllers

import (
	"context"
	flinkv1 "github.com/123shang60/flink-session-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *FlinkSessionReconciler) deleteDeploy(session *flinkv1.FlinkSession) error {
	klog.Info("开始查询 flink session 集群 deployment")

	var deploy appsv1.Deployment
	if err := r.Get(context.Background(), client.ObjectKeyFromObject(session), &deploy); err != nil {
		klog.Error("Get flink session deployment error!", err)
		return err
	}

	klog.Infof("找到session 集群，uuid : %s, 开始删除！", deploy.UID)
	if err := r.Delete(context.Background(), &deploy); err != nil {
		klog.Error("deployment 删除失败！", err)
		return err
	}

	klog.Info("删除完毕！")

	klog.Info("deployment is ", deploy)
	return nil
}
