package controllers

import (
	"context"
	flinkv1 "github.com/123shang60/flink-session-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *FlinkSessionReconciler) updateNodePort(session *flinkv1.FlinkSession) error {
	klog.Info("开始更新 svc port 信息！")
	var svc corev1.ServiceList

	if err := r.List(context.Background(), &svc, client.MatchingLabels{
		"type": FlinkNativeType,
		"app":  session.GetName(),
	}, client.InNamespace(session.GetNamespace())); err != nil {
		klog.Error("svc 查询失败！", err)
		return err
	}

	klog.Info("获取svc信息成功，本次获取结果：", svc)
	for _, s := range svc.Items {
		if s.Spec.Type == corev1.ServiceTypeNodePort && len(s.Spec.Ports) != 0 {
			klog.Info("找到对外暴露的 svc ！ ", s)
			for _, port := range s.Spec.Ports {
				if port.Name == FlinkRestPortName {
					session.Status.Port = port.NodePort
					if err := r.Status().Update(context.Background(), session); err != nil {
						klog.Error("更新 status 失败！", err)
						return err
					}
				}
			}
		}
	}
	return nil
}
