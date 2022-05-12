package controllers

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type FlinkSessionJobHandler struct {
}

// Create implements EventHandler.
func (e *FlinkSessionJobHandler) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		klog.Error("CreateEvent received with no metadata", "event", evt)
		return
	}
	if job, ok := evt.Object.(*batchv1.Job); ok {
		for _, reference := range job.GetObjectMeta().GetOwnerReferences() {
			if reference.APIVersion == `flink.shang12360.cn/v1` &&
				reference.Kind == `FlinkSession` {
				q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
					Name:      reference.Name,
					Namespace: job.GetNamespace(),
				}})
			}
		}
	}
}

// Update implements EventHandler.
func (e *FlinkSessionJobHandler) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	switch {
	case evt.ObjectNew != nil:
		if job, ok := evt.ObjectNew.(*batchv1.Job); ok {
			for _, reference := range job.GetObjectMeta().GetOwnerReferences() {
				if reference.APIVersion == `flink.shang12360.cn/v1` &&
					reference.Kind == `FlinkSession` {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      reference.Name,
						Namespace: job.GetNamespace(),
					}})
				}
			}
		}
	case evt.ObjectOld != nil:
		if svc, ok := evt.ObjectOld.(*corev1.Service); ok {
			klog.Info("old", svc)
		}
	default:
		klog.Error("UpdateEvent received with no metadata", "event", evt)
	}
}

// Delete implements EventHandler.
func (e *FlinkSessionJobHandler) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		klog.Error("DeleteEvent received with no metadata", "event", evt)
		return
	}
	//if svc, ok := evt.Object.(*corev1.Service); ok {
	//	klog.Info("del", svc)
	//}
}

// Generic implements EventHandler.
func (e *FlinkSessionJobHandler) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		klog.Error("GenericEvent received with no metadata", "event", evt)
		return
	}
	if job, ok := evt.Object.(*batchv1.Job); ok {
		for _, reference := range job.GetObjectMeta().GetOwnerReferences() {
			if reference.APIVersion == `flink.shang12360.cn/v1` &&
				reference.Kind == `FlinkSession` {
				q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
					Name:      reference.Name,
					Namespace: job.GetNamespace(),
				}})
			}
		}
	}
}
