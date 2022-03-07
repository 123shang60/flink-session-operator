package controllers

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type FlinkSessionServiceHandler struct {
}

// Create implements EventHandler.
func (e *FlinkSessionServiceHandler) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		klog.Error("CreateEvent received with no metadata", "event", evt)
		return
	}
	if svc, ok := evt.Object.(*corev1.Service); ok {
		if svc.GetLabels() != nil {
			if typ, ok := svc.GetLabels()["type"]; ok && typ == FlinkNativeType {
				if app, ok := svc.GetLabels()["app"]; ok {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app,
						Namespace: svc.GetNamespace(),
					}})
				}
			}
		}
	}
}

// Update implements EventHandler.
func (e *FlinkSessionServiceHandler) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	switch {
	case evt.ObjectNew != nil:
		if svc, ok := evt.ObjectNew.(*corev1.Service); ok {
			if svc.GetLabels() != nil {
				if typ, ok := svc.GetLabels()["type"]; ok && typ == FlinkNativeType {
					if app, ok := svc.GetLabels()["app"]; ok {
						q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
							Name:      app,
							Namespace: svc.GetNamespace(),
						}})
					}
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
func (e *FlinkSessionServiceHandler) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		klog.Error("DeleteEvent received with no metadata", "event", evt)
		return
	}
	//if svc, ok := evt.Object.(*corev1.Service); ok {
	//	klog.Info("del", svc)
	//}
}

// Generic implements EventHandler.
func (e *FlinkSessionServiceHandler) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		klog.Error("GenericEvent received with no metadata", "event", evt)
		return
	}
	if svc, ok := evt.Object.(*corev1.Service); ok {
		if svc.GetLabels() != nil {
			if typ, ok := svc.GetLabels()["type"]; ok && typ == FlinkNativeType {
				if app, ok := svc.GetLabels()["app"]; ok {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app,
						Namespace: svc.GetNamespace(),
					}})
				}
			}
		}
	}
}
