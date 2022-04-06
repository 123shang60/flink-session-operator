package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	flinkv1 "github.com/123shang60/flink-session-operator/api/v1"
	yaml2 "github.com/ghodss/yaml"
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
		bootConfigMap.Data[`pod-template.yaml`] = generatePodTemplate(
			session.Spec.BalancedSchedule,
			session.Name,
			session.Namespace,
		)
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

	err = r.Create(context.Background(), bootConfigMap)
	if err != nil {
		klog.Info("创建 configmap 失败!", err)
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

func generatePodTemplate(strategy, appName, nameSpace string) string {
	//var podtemplate *apiv1.Pod
	podtemplate := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{},
		},
	}
	if strategy == flinkv1.PreferredDuringScheduling {
		podtemplate.Spec.Affinity = &apiv1.Affinity{
			PodAntiAffinity: &apiv1.PodAntiAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []apiv1.WeightedPodAffinityTerm{
					{
						Weight: 100,
						PodAffinityTerm: apiv1.PodAffinityTerm{
							LabelSelector: &metav1.LabelSelector{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									{
										Key:      "app",
										Operator: metav1.LabelSelectorOpIn,
										Values:   []string{appName},
									},
									{
										Key:      "type",
										Operator: metav1.LabelSelectorOpIn,
										Values:   []string{FlinkNativeType},
									},
								},
							},
							Namespaces:  []string{nameSpace},
							TopologyKey: DefaultTopologyKey,
						},
					},
				},
			},
		}
	}

	if strategy == flinkv1.RequiredDuringScheduling {
		podtemplate.Spec.Affinity = &apiv1.Affinity{
			PodAntiAffinity: &apiv1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []apiv1.PodAffinityTerm{
					{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      "app",
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{appName},
								},
								{
									Key:      "type",
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{FlinkNativeType},
								},
							},
						},
						Namespaces:  []string{nameSpace},
						TopologyKey: DefaultTopologyKey,
					},
				},
			},
		}
	}

	byte, _ := json.Marshal(podtemplate)
	yaml, _ := yaml2.JSONToYAML(byte)

	return string(yaml)
}
