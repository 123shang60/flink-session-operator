package v1

import (
	"encoding/json"
	"fmt"
	"github.com/123shang60/flink-session-operator/pkg"
	yaml2 "github.com/ghodss/yaml"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
)

func (f *FlinkSession) GenerateCommand() (string, error) {
	// 构建 name 以及 namespaces
	var command pkg.Command
	command.WriteString(`$FLINK_HOME/bin/kubernetes-session.sh `)
	command.FieldConfig("kubernetes.cluster-id", f.Name).FieldConfig("kubernetes.namespace", f.Namespace)
	// service account
	command.FieldConfig("kubernetes.service-account", f.Spec.Sa)
	// images
	command.FieldConfig("kubernetes.container.image", f.Spec.Image)
	if f.Spec.ImageSecret != nil {
		command.FieldConfig("kubernetes.container.image.pull-secrets", *f.Spec.ImageSecret)
	}
	// resource
	// jobManager
	command.FieldConfig("kubernetes.jobmanager.cpu", f.Spec.Resource.JobManager.CPU)
	command.FieldConfig("jobmanager.memory.process.size", f.Spec.Resource.JobManager.Memory)
	command.FieldConfig("jobmanager.memory.jvm-metaspace.size", f.Spec.Resource.JobManager.JvmMetaspace)
	command.FieldConfig("jobmanager.memory.task.off-heap.size", f.Spec.Resource.JobManager.OffHeap)

	// taskManager
	command.FieldConfig("kubernetes.taskmanager.cpu", f.Spec.Resource.TaskManager.CPU)
	command.FieldConfig("taskmanager.memory.process.size", f.Spec.Resource.TaskManager.Memory)
	command.FieldConfig("taskmanager.memory.jvm-metaspace.size", f.Spec.Resource.TaskManager.JvmMetaspace)
	command.FieldConfig("taskmanager.memory.framework.heap.size", f.Spec.Resource.TaskManager.Framework.Heap)
	command.FieldConfig("taskmanager.memory.framework.off-heap.size", f.Spec.Resource.TaskManager.Framework.OffHeap)
	command.FieldConfig("taskmanager.memory.task.off-heap.size", f.Spec.Resource.TaskManager.Task.OffHeap)

	// taskManager NetWork
	if f.Spec.Resource.TaskManager.NetWork.Max != "" || f.Spec.Resource.TaskManager.NetWork.Min != "" {
		command.FieldConfig("taskmanager.memory.network.min", f.Spec.Resource.TaskManager.NetWork.Min)
		command.FieldConfig("taskmanager.memory.network.max", f.Spec.Resource.TaskManager.NetWork.Max)
	} else {
		command.FieldConfig("taskmanager.memory.network.fraction", f.Spec.Resource.TaskManager.NetWork.Fraction)
	}
	// taskManager Managed
	if f.Spec.Resource.TaskManager.Managed.Max != "" || f.Spec.Resource.TaskManager.Managed.Min != "" {
		command.FieldConfig("taskmanager.memory.managed.min", f.Spec.Resource.TaskManager.Managed.Min)
		command.FieldConfig("taskmanager.memory.managed.max", f.Spec.Resource.TaskManager.Managed.Max)
	} else {
		command.FieldConfig("taskmanager.memory.managed.fraction", f.Spec.Resource.TaskManager.Managed.Fraction)
	}

	// taskManager slot
	num := strconv.Itoa(int(f.Spec.NumberOfTaskSlots))
	command.FieldConfig("taskmanager.numberOfTaskSlots", num)

	// s3 状态后端
	command.FieldConfig("state.backend", "filesystem")
	command.FieldConfig("s3.endpoint", fmt.Sprintf("http://%s", f.Spec.S3.EndPoint))
	command.FieldConfig("s3.access-key", f.Spec.S3.AccessKey)
	command.FieldConfig("s3.secret-key", f.Spec.S3.SecretKey)
	command.FieldConfig("s3.path.style.access", "true")
	command.FieldConfig("state.checkpoints.dir", fmt.Sprintf("s3://%s/%s/flink/checkpoints", f.Spec.S3.Bucket, f.Name))
	command.FieldConfig("state.savepoints.dir", fmt.Sprintf("s3://%s/%s/flink/savepoints", f.Spec.S3.Bucket, f.Name))
	command.FieldConfig("historyserver.archive.fs.dir", fmt.Sprintf("s3://%s/%s/flink/completed-jobs", f.Spec.S3.Bucket, f.Name))
	command.FieldConfig("jobmanager.archive.fs.dir", fmt.Sprintf("s3://%s/%s/flink/archive", f.Spec.S3.Bucket, f.Name))
	command.FieldConfig("state.backend.incremental", "true")
	command.FieldConfig("fs.overwrite-files", "true")

	// ha
	switch f.Spec.HA.Typ {
	case ZKHA:
		command.FieldConfig("high-availability", "zookeeper")
		command.FieldConfig("high-availability.zookeeper.quorum", f.Spec.HA.Quorum)
		command.FieldConfig("high-availability.zookeeper.path.root", f.Spec.HA.Path)
	case CONFIGMAPHA:
		command.FieldConfig("high-availability", "org.apache.flink.kubernetes.highavailability.KubernetesHaServicesFactory")
	default:
	}

	if f.Spec.HA.Typ == ZKHA || f.Spec.HA.Typ == CONFIGMAPHA {
		command.FieldConfig("high-availability.storageDir", fmt.Sprintf("s3://%s/%s/flink/ha/metadata", f.Spec.S3.Bucket, f.Name))
	}

	// nodeSelector
	if f.Spec.NodeSelector != nil && len(f.Spec.NodeSelector) != 0 {
		selector := buildNodeSelector(f.Spec.NodeSelector)
		command.FieldConfig("kubernetes.taskmanager.node-selector", selector)
		command.FieldConfig("kubernetes.jobmanager.node-selector", selector)
	}
	if f.Spec.BalancedSchedule != NoneScheduling {
		command.FieldConfig("kubernetes.pod-template-file", "/opt/flink/template/pod-template.yaml")
		command.FieldConfig("kubernetes.pod-template-file.jobmanager", "/opt/flink/template/pod-template.yaml")
		command.FieldConfig("kubernetes.pod-template-file.taskmanager", "/opt/flink/template/pod-template.yaml")
	}

	// 其他的必配项目
	command.FieldConfig("env.java.opts", `"-XX:+UseG1GC"`)
	command.FieldConfig("kubernetes.rest-service.exposed.type", "NodePort")

	return command.Build(), nil
}

func buildNodeSelector(selector map[string]string) string {
	var builder strings.Builder

	for k, v := range selector {
		builder.WriteString(k)
		builder.WriteString(":")
		builder.WriteString(v)
		builder.WriteString(",")
	}
	str := builder.String()
	if len(str) > 0 {
		return str[:len(str)-1]
	} else {
		return str
	}
}

func (f *FlinkSession) GeneratePodTemplate() string {
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
	if f.Spec.BalancedSchedule == PreferredDuringScheduling {
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
										Values:   []string{f.Name},
									},
									{
										Key:      "type",
										Operator: metav1.LabelSelectorOpIn,
										Values:   []string{FlinkNativeType},
									},
								},
							},
							Namespaces:  []string{f.Namespace},
							TopologyKey: DefaultTopologyKey,
						},
					},
				},
			},
		}
	}

	if f.Spec.BalancedSchedule == RequiredDuringScheduling {
		podtemplate.Spec.Affinity = &apiv1.Affinity{
			PodAntiAffinity: &apiv1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []apiv1.PodAffinityTerm{
					{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      "app",
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{f.Name},
								},
								{
									Key:      "type",
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{FlinkNativeType},
								},
							},
						},
						Namespaces:  []string{f.Namespace},
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
