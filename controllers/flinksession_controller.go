/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"
	"github.com/cnf/structhash"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	flinkv1 "github.com/123shang60/flink-session-operator/api/v1"
)

// FlinkSessionReconciler reconciles a FlinkSession object
type FlinkSessionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type changeReconciler struct {
	Spec              flinkv1.FlinkSessionSpec
	DeletionTimestamp *v1.Time
}

//+kubebuilder:rbac:groups=flink.shang12360.cn,resources=flinksessions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flink.shang12360.cn,resources=flinksessions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=flink.shang12360.cn,resources=flinksessions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FlinkSession object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *FlinkSessionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	var session flinkv1.FlinkSession

	if err := r.Get(ctx, req.NamespacedName, &session); err != nil {
		klog.Error("unable to fetch flinksession:", err)
		return ctrl.Result{}, nil
	}

	changeHash := changeReconciler{
		Spec:              session.Spec,
		DeletionTimestamp: session.DeletionTimestamp,
	}

	if val, ok := session.ObjectMeta.Annotations[HashAnnotations]; !ok {
		session.ObjectMeta.Annotations[HashAnnotations], _ = structhash.Hash(changeHash, 1)
		if err := r.Update(ctx, &session); err != nil {
			klog.Error("Add empty HashAnnotations error:: ", err)
			return ctrl.Result{}, err
		}
	} else {
		thisHash, _ := structhash.Hash(changeHash, 1)
		if val != thisHash {
			klog.Info("spec change!")
		} else {
			klog.Info("spec nochange!")
			return ctrl.Result{}, nil
		}
	}

	b, _ := json.Marshal(session)

	klog.Info("session is :", string(b))

	// examine DeletionTimestamp to determine if object is under deletion
	if session.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(&session, FlinkSessionFinalizerName) {
			controllerutil.AddFinalizer(&session, FlinkSessionFinalizerName)
			if err := r.Update(ctx, &session); err != nil {
				klog.Error("Update AddFinalizer error: ", err)
				return ctrl.Result{}, err
			}
		}

		if err := r.updateExternalResources(&session); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&session, FlinkSessionFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.deleteExternalResources(&session); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				klog.Error("Delete ExternalResources error: ", err)
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&session, FlinkSessionFinalizerName)
			if err := r.Update(ctx, &session); err != nil {
				klog.Error("Delete AddFinalizer error:: ", err)
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FlinkSessionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&flinkv1.FlinkSession{}).
		Complete(r)
}

// updateExternalResources create or delete all other resources
func (r *FlinkSessionReconciler) updateExternalResources(session *flinkv1.FlinkSession) error {
	//
	// delete any external resources associated with the cronJob
	//
	// Ensure that delete implementation is idempotent and safe to invoke
	// multiple times for same object.
	klog.Info("do somethings update!:", session.Spec.Image)
	return nil
}

// deleteExternalResources remove all other resources
func (r *FlinkSessionReconciler) deleteExternalResources(session *flinkv1.FlinkSession) error {
	//
	// delete any external resources associated with the cronJob
	//
	// Ensure that delete implementation is idempotent and safe to invoke
	// multiple times for same object.
	klog.Info("do somethings delete!:", session.Spec.Image)
	return nil
}