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

package v1

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/Masterminds/semver/v3"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var flinksessionlog = logf.Log.WithName("flinksession-resource")

func (r *FlinkSession) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-flink-shang12360-cn-v1-flinksession,mutating=true,failurePolicy=fail,sideEffects=None,groups=flink.shang12360.cn,resources=flinksessions,verbs=create;update,versions=v1,name=mflinksession.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &FlinkSession{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *FlinkSession) Default() {
	flinksessionlog.Info("default", "name", r.Name)

	b, _ := json.Marshal(r)

	klog.Info("default hook is :", string(b))

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-flink-shang12360-cn-v1-flinksession,mutating=false,failurePolicy=fail,sideEffects=None,groups=flink.shang12360.cn,resources=flinksessions,verbs=create;update,versions=v1,name=vflinksession.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &FlinkSession{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *FlinkSession) ValidateCreate() error {
	flinksessionlog.Info("validate create", "name", r.Name)

	b, _ := json.Marshal(r)

	klog.Info("validate create hook is :", string(b))
	// TODO(user): fill in your validation logic upon object creation.
	if r.Spec.Security.Kerberos != nil {
		if _, err := base64.StdEncoding.DecodeString(r.Spec.Security.Kerberos.Base64Keytab); err != nil {
			klog.Error("创建校验失败！不是正确的 base64！", err)
			return errors.New("error base64 format ! :" + err.Error())
		}
	}

	if r.Spec.FlinkVersion != nil {
		_, err := semver.NewVersion(*r.Spec.FlinkVersion)
		if err != nil {
			klog.Error("创建校验失败！不是符合语义标准的版本号！", err)
			return errors.New("error semver version format！：" + err.Error())
		}
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *FlinkSession) ValidateUpdate(old runtime.Object) error {
	flinksessionlog.Info("validate update", "name", r.Name)

	b, _ := json.Marshal(r)

	klog.Info("validate update hook is :", string(b))

	b, _ = json.Marshal(old)

	klog.Info("validate update old hook is :", string(b))
	// TODO(user): fill in your validation logic upon object update.
	if r.Spec.Security.Kerberos != nil {
		if _, err := base64.StdEncoding.DecodeString(r.Spec.Security.Kerberos.Base64Keytab); err != nil {
			klog.Error("创建校验失败！不是正确的 base64！", err)
			return errors.New("error base64 format ! :" + err.Error())
		}
	}

	if r.Spec.FlinkVersion != nil {
		_, err := semver.NewVersion(*r.Spec.FlinkVersion)
		if err != nil {
			klog.Error("创建校验失败！不是符合语义标准的版本号！", err)
			return errors.New("error semver version format！：" + err.Error())
		}
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *FlinkSession) ValidateDelete() error {
	flinksessionlog.Info("validate delete", "name", r.Name)

	b, _ := json.Marshal(r)

	klog.Info("validate delete hook is :", string(b))
	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
