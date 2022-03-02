package controllers

import (
	"context"
	"fmt"
	flinkv1 "github.com/123shang60/flink-session-operator/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

func (r *FlinkSessionReconciler) commitBootJob(session *flinkv1.FlinkSession) error {
	command, err := session.GenerateCommand()
	if err != nil {
		return err
	}

	bootJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("boot-%s-%d", session.Name, time.Now().Unix()),
			Namespace: session.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"flink": "flink-session-operator",
					},
				},
				Spec: apiv1.PodSpec{
					ServiceAccountName: session.Spec.Sa,
					RestartPolicy:      apiv1.RestartPolicyOnFailure,
					Containers: []apiv1.Container{
						{
							Name:  "start",
							Image: session.Spec.Image,
							Command: []string{
								"bash",
								"-c",
								command,
							},
						},
					},
				},
			},
		},
	}

	if session.Spec.ImageSecret != nil && len(*session.Spec.ImageSecret) != 0 {
		bootJob.Spec.Template.Spec.ImagePullSecrets = []apiv1.LocalObjectReference{
			{
				Name: *session.Spec.ImageSecret,
			},
		}
	}

	err = controllerutil.SetControllerReference(session, bootJob, r.Scheme)
	if err != nil {
		klog.Error("设置 reference 失败!")
	}

	err = r.Create(context.Background(), bootJob)
	if err != nil {
		klog.Info("创建 job 失败!", err)
		return err
	}
	klog.Info("boot job 创建成功!")

	return nil
}
