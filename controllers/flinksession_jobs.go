package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	flinkv1 "github.com/123shang60/flink-session-operator/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

func (r *FlinkSessionReconciler) commitBootJob(session *flinkv1.FlinkSession) error {
	command, err := session.GenerateCommand()
	if err != nil {
		return err
	}

	jobName := fmt.Sprintf("boot-%s-%d", session.Name, time.Now().Unix())
	defaultMode := DefaultMode

	bootJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
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
							VolumeMounts: make([]apiv1.VolumeMount, 0),
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "config",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: jobName,
									},
									DefaultMode: &defaultMode,
								},
							},
						},
						{
							Name: "secret",
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName:  jobName,
									DefaultMode: &defaultMode,
								},
							},
						},
					},
				},
			},
		},
	}

	bootConfigMap := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: session.Namespace,
		},
		Data: make(map[string]string),
	}

	bootSecret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: session.Namespace,
		},
		Data: make(map[string][]byte),
		Type: apiv1.SecretTypeOpaque,
	}

	// kerberos
	if session.Spec.Security.Kerberos != nil {
		keytab, err := base64.StdEncoding.DecodeString(session.Spec.Security.Kerberos.Base64Keytab)
		if err != nil {
			klog.Error("keytab 解码失败，无效的base64 : ", err)
			klog.Error("keytab 原始码值: ", session.Spec.Security.Kerberos.Base64Keytab)
			return err
		}
		bootSecret.Data["flink.keytab"] = keytab
		bootConfigMap.Data["krb5.conf"] = session.Spec.Security.Kerberos.Krb5
		bootJob.Spec.Template.Spec.Containers[0].VolumeMounts = append(
			bootJob.Spec.Template.Spec.Containers[0].VolumeMounts,
			apiv1.VolumeMount{
				Name:      "config",
				MountPath: "/opt/flink/conf/krb5.conf",
				SubPath:   "krb5.conf",
			},
		)
		bootJob.Spec.Template.Spec.Containers[0].VolumeMounts = append(
			bootJob.Spec.Template.Spec.Containers[0].VolumeMounts,
			apiv1.VolumeMount{
				Name:      "secret",
				MountPath: "/opt/flink/conf/flink.keytab",
				SubPath:   "flink.keytab",
			},
		)
	}

	if session.Spec.ImageSecret != nil && len(*session.Spec.ImageSecret) != 0 {
		bootJob.Spec.Template.Spec.ImagePullSecrets = []apiv1.LocalObjectReference{
			{
				Name: *session.Spec.ImageSecret,
			},
		}
	}

	if len(session.Spec.Config.FlinkConf) != 0 {
		bootConfigMap.Data[`flink-conf.yaml`] = session.Spec.Config.FlinkConf
		bootJob.Spec.Template.Spec.Containers[0].VolumeMounts = append(
			bootJob.Spec.Template.Spec.Containers[0].VolumeMounts,
			apiv1.VolumeMount{
				Name:      "config",
				MountPath: "/opt/flink/conf/flink-conf.yaml",
				SubPath:   "flink-conf.yaml",
			},
		)
	}

	if len(session.Spec.Config.Log4j) != 0 {
		bootConfigMap.Data[`log4j-console.properties`] = session.Spec.Config.Log4j
		bootJob.Spec.Template.Spec.Containers[0].VolumeMounts = append(
			bootJob.Spec.Template.Spec.Containers[0].VolumeMounts,
			apiv1.VolumeMount{
				Name:      "config",
				MountPath: "/opt/flink/conf/log4j-console.properties",
				SubPath:   "log4j-console.properties",
			},
		)
	}

	if len(session.Spec.Config.LogBack) != 0 {
		bootConfigMap.Data[`logback-console.xml`] = session.Spec.Config.LogBack
		bootJob.Spec.Template.Spec.Containers[0].VolumeMounts = append(
			bootJob.Spec.Template.Spec.Containers[0].VolumeMounts,
			apiv1.VolumeMount{
				Name:      "config",
				MountPath: "/opt/flink/conf/logback-console.xml",
				SubPath:   "logback-console.xml",
			},
		)
	}

	if session.Spec.BalancedSchedule == flinkv1.PreferredDuringScheduling ||
		session.Spec.BalancedSchedule == flinkv1.RequiredDuringScheduling ||
		(session.Spec.Volumes != nil && len(session.Spec.Volumes) == 0) {
		bootConfigMap.Data[`pod-template.yaml`] = session.GeneratePodTemplate()
		bootJob.Spec.Template.Spec.Containers[0].VolumeMounts = append(
			bootJob.Spec.Template.Spec.Containers[0].VolumeMounts,
			apiv1.VolumeMount{
				Name:      "config",
				MountPath: "/opt/flink/template/pod-template.yaml",
				SubPath:   "pod-template.yaml",
			},
		)
	}

	err = controllerutil.SetControllerReference(session, bootJob, r.Scheme)
	if err != nil {
		klog.Error("job 设置 reference 失败!", err)
	}

	err = controllerutil.SetControllerReference(session, bootConfigMap, r.Scheme)
	if err != nil {
		klog.Error("configmap 设置 reference 失败!", err)
	}

	err = controllerutil.SetControllerReference(session, bootSecret, r.Scheme)
	if err != nil {
		klog.Error("secret 设置 reference 失败!", err)
	}

	err = r.Create(context.Background(), bootConfigMap)
	if err != nil {
		klog.Info("创建 configmap 失败!", err)
		return err
	}

	err = r.Create(context.Background(), bootSecret)
	if err != nil {
		klog.Info("创建 secret 失败!", err)
		return err
	}

	err = r.Create(context.Background(), bootJob)
	if err != nil {
		klog.Info("创建 job 失败!", err)
		return err
	}
	klog.Info("boot job 创建成功!")

	return nil
}

func (r *FlinkSessionReconciler) cleanBootJob(session *flinkv1.FlinkSession, success int32) error {
	jobList := batchv1.JobList{}
	if err := r.List(context.Background(), &jobList, client.MatchingLabels{
		"flink": "flink-session-operator",
	}, client.InNamespace(session.GetNamespace())); err != nil {
		klog.Error("获取列表失败!", err)
		return err
	} else {
		for _, job := range jobList.Items {
			if job.Status.Succeeded == success {
				for _, reference := range job.GetObjectMeta().GetOwnerReferences() {
					if reference.APIVersion == `flink.shang12360.cn/v1` &&
						reference.Kind == `FlinkSession` &&
						reference.Name == session.Name {
						jobName := job.Name
						propagationPolicy := metav1.DeletePropagationBackground
						if err := r.Delete(context.Background(), &job, &client.DeleteOptions{PropagationPolicy: &propagationPolicy}); err != nil {
							klog.Error("删除job失败!", err)
						} else {
							klog.Info("清理 job :", jobName)
						}
						if err := r.Delete(context.Background(), &apiv1.ConfigMap{
							ObjectMeta: metav1.ObjectMeta{
								Name:      jobName,
								Namespace: session.Namespace,
							},
						}); err != nil {
							klog.Error("删除job configmap失败!", err)
						} else {
							klog.Info("清理 job configmap:", jobName)
						}

						if err := r.Delete(context.Background(), &apiv1.Secret{
							ObjectMeta: metav1.ObjectMeta{
								Name:      jobName,
								Namespace: session.Namespace,
							},
						}); err != nil {
							klog.Error("删除job secret 失败!", err)
						} else {
							klog.Info("清理 job secret:", jobName)
						}

						break
					}
				}
			}
		}
	}
	return nil
}
