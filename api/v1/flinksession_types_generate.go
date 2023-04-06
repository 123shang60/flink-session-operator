package v1

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"k8s.io/klog/v2"
	"strconv"
	"strings"

	"github.com/123shang60/flink-session-operator/pkg"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func flinkHAConfig(f *FlinkSession, command *pkg.Command) {
	flinkConfigs := make([]func(f *FlinkSession, command *pkg.Command) bool, 0)
	flinkConfigs = append(flinkConfigs, flinkHAConfigMod114)
	flinkConfigs = append(flinkConfigs, flinkHAConfigMod116)
	flinkConfigs = append(flinkConfigs, flinkHAConfigMod117)
	flinkConfigs = append(flinkConfigs, flinkHAConfigModdefault)

	for _, fun := range flinkConfigs {
		if fun(f, command) {
			return
		}
	}
}

func flinkHAConfigModdefault(f *FlinkSession, command *pkg.Command) bool {
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

	return true
}

func flinkHAConfigMod114(f *FlinkSession, command *pkg.Command) bool {
	constraint, err := semver.NewConstraint("<= 1.15")
	if err != nil {
		klog.Error("it's bugs！", err)
		return false
	}

	if f.Spec.FlinkVersion != nil {
		v, err := semver.NewVersion(*f.Spec.FlinkVersion)
		if err != nil {
			klog.Error("it's bugs！", err)
			return false
		}

		if constraint.Check(v) {
			return flinkHAConfigModdefault(f, command)
		}
	}

	return false
}

func flinkHAConfigMod116(f *FlinkSession, command *pkg.Command) bool {
	constraint, err := semver.NewConstraint(">= 1.16 < 1.17")
	if err != nil {
		klog.Error("it's bugs！", err)
		return false
	}

	if f.Spec.FlinkVersion != nil {
		v, err := semver.NewVersion(*f.Spec.FlinkVersion)
		if err != nil {
			klog.Error("it's bugs！", err)
			return false
		}

		if constraint.Check(v) {
			switch f.Spec.HA.Typ {
			case ZKHA:
				command.FieldConfig("high-availability", "ZOOKEEPER")
				command.FieldConfig("high-availability.zookeeper.quorum", f.Spec.HA.Quorum)
				command.FieldConfig("high-availability.zookeeper.path.root", f.Spec.HA.Path)
			case CONFIGMAPHA:
				command.FieldConfig("high-availability", "KUBERNETES")
			default:
			}

			if f.Spec.HA.Typ == ZKHA || f.Spec.HA.Typ == CONFIGMAPHA {
				command.FieldConfig("high-availability.storageDir", fmt.Sprintf("s3://%s/%s/flink/ha/metadata", f.Spec.S3.Bucket, f.Name))
			}

			return true
		}
	}

	return false
}

func flinkHAConfigMod117(f *FlinkSession, command *pkg.Command) bool {
	constraint, err := semver.NewConstraint(">= 1.17")
	if err != nil {
		klog.Error("it's bugs！", err)
		return false
	}

	if f.Spec.FlinkVersion != nil {
		v, err := semver.NewVersion(*f.Spec.FlinkVersion)
		if err != nil {
			klog.Error("it's bugs！", err)
			return false
		}

		if constraint.Check(v) {
			switch f.Spec.HA.Typ {
			case ZKHA:
				command.FieldConfig("high-availability.type", "ZOOKEEPER")
				command.FieldConfig("high-availability.zookeeper.quorum", f.Spec.HA.Quorum)
				command.FieldConfig("high-availability.zookeeper.path.root", f.Spec.HA.Path)
			case CONFIGMAPHA:
				command.FieldConfig("high-availability.type", "KUBERNETES")
			default:
			}

			if f.Spec.HA.Typ == ZKHA || f.Spec.HA.Typ == CONFIGMAPHA {
				command.FieldConfig("high-availability.storageDir", fmt.Sprintf("s3://%s/%s/flink/ha/metadata", f.Spec.S3.Bucket, f.Name))
			}

			return true
		}
	}

	return false
}

func (f *FlinkSession) GenerateCommand() (string, error) {
	// 构建 name 以及 namespaces
	var command pkg.Command
	if f.Spec.ApplicationConfig == nil {
		command.WriteString(`$FLINK_HOME/bin/kubernetes-session.sh `)
	} else {
		command.WriteString(`$FLINK_HOME/bin/flink run-application --target kubernetes-application `)
	}

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
	flinkHAConfig(f, &command)

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

	// Security
	// Kerberos
	if f.Spec.Security.Kerberos != nil {
		command.FieldConfig("security.kerberos.login.keytab", "/opt/flink/conf/flink.keytab")
		command.FieldConfig("security.kerberos.login.principal", f.Spec.Security.Kerberos.Principal)
		command.FieldConfig("security.kerberos.login.contexts", f.Spec.Security.Kerberos.Contexts)
		command.FieldConfig("security.kerberos.krb5-conf.path", "/opt/flink/conf/krb5.conf")
		if *f.Spec.Security.Kerberos.UseTicketCache {
			command.FieldConfig("security.kerberos.login.use-ticket-cache", "true")
		} else {
			command.FieldConfig("security.kerberos.login.use-ticket-cache", "false")
		}
	}

	// 其他的必配项目
	// env.java.opts 仅在未特殊指定的情况下增加
	if !strings.Contains(f.Spec.Config.FlinkConf, "env.java.opts") {
		command.FieldConfig("env.java.opts", `"-XX:+UseG1GC"`)
	}
	command.FieldConfig("kubernetes.rest-service.exposed.type", "NodePort")

	if f.Spec.ApplicationConfig != nil {
		command.WriteString(fmt.Sprintf(`-p %d `, f.Spec.ApplicationConfig.Parallelism))
		command.WriteString(f.Spec.ApplicationConfig.JarPath)
		command.WriteString(` `)
		for _, arg := range f.Spec.ApplicationConfig.Args {
			command.WriteString(strings.Trim(arg, ` `))
			command.WriteString(` `)
		}
	}

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
			Containers: []apiv1.Container{
				{
					Name: "flink-main-container",
				},
			},
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

	if f.Spec.Volumes != nil && len(f.Spec.Volumes) != 0 {
		podtemplate.Spec.Volumes = f.Spec.Volumes
		if f.Spec.VolumeMounts != nil && len(f.Spec.VolumeMounts) != 0 {
			for k, _ := range podtemplate.Spec.Containers {
				podtemplate.Spec.Containers[k].VolumeMounts = f.Spec.VolumeMounts
			}
		}
	}

	byte, _ := yaml.Marshal(podtemplate)

	return string(byte)
}
