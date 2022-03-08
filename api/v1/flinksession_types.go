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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FlinkSessionSpec defines the desired state of FlinkSession
type FlinkSessionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// 填写 flink 运行镜像 ,这个镜像为 flink taskmanager 以及 jobmanager 共用，不支持双镜像，必填
	//+kubebuilder:validation:MinLength=1
	Image string `json:"image,omitempty"`

	// 填写 flink 镜像拉取 secret
	ImageSecret *string `json:"imageSecret,omitempty"`

	// SA，填写集群运行的 k8s service account 配置
	//+kubebuilder:validation:MinLength=1
	Sa string `json:"sa,omitempty"`

	// flink 运行资源配置
	Resource FlinkResource `json:"resource,omitempty"`

	// taskmanager 可用槽位，对应 taskmanager.numberOfTaskSlots
	//+kubebuilder:validation:Minimum=1
	NumberOfTaskSlots int32 `json:"numberOfTaskSlots,omitempty"`

	// S3 配置
	S3 FlinkS3 `json:"s3,omitempty"`

	// flink ha 配置
	HA FlinkHA `json:"ha,omitempty"`

	// 自定义配置项
	//+nullable
	Config FlinkConfig `json:"config,omitempty"`

	// NodeSelector
	//+nullable
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// 均衡调度策略 ，可选值： Required 必须每个节点调度一个 Preferred 尽可能每个节点调度一个 None 不设置均衡调度
	//+kubebuilder:validation:MinLength=1
	//+kubebuilder:validation:Enum={Required,Preferred,None}
	BalancedSchedule string `json:"balancedSchedule,omitempty"`
}

type FlinkResource struct {
	// jobmanager 资源配置
	JobManager JobManagerFlinkResource `json:"jobManager,omitempty"`

	// taskmanager 资源配置
	TaskManager TaskManagerFlinkResource `json:"taskManager,omitempty"`
}

// JobManagerFlinkResource job manager 资源配置
type JobManagerFlinkResource struct {
	// kubernetes.jobmanager.cpu
	CPU string `json:"cpu,omitempty"`
	// jobmanager.memory.process.size
	Memory string `json:"memory,omitempty"`
	// jobmanager.memory.jvm-metaspace.size
	JvmMetaspace string `json:"jvm-metaspace,omitempty"`
	// jobmanager.memory.off-heap.size
	OffHeap string `json:"off-heap,omitempty"`
}

// TaskManagerFlinkResource task manager 资源配置
type TaskManagerFlinkResource struct {
	// kubernetes.taskmanager.cpu
	CPU string `json:"cpu,omitempty"`
	// taskmanager.memory.process.size
	Memory string `json:"memory,omitempty"`
	// taskmanager.memory.jvm-metaspace.size
	JvmMetaspace string `json:"jvm-metaspace,omitempty"`

	Framework TaskManagerFrameworkFlinkResource `json:"framework,omitempty"`
	Task      TaskManagerTaskFlinkResource      `json:"task,omitempty"`
	NetWork   TaskManagerNetWorkFlinkResource   `json:"netWork,omitempty"`
	Managed   TaskManagerManagedFlinkResource   `json:"managed,omitempty"`
}

type TaskManagerFrameworkFlinkResource struct {
	// taskmanager.memory.framework.heap.size
	Heap string `json:"heap,omitempty"`
	// taskmanager.memory.framework.off-heap.size
	OffHeap string `json:"off-heap,omitempty"`
}

type TaskManagerTaskFlinkResource struct {
	// taskmanager.memory.task.off-heap.size
	OffHeap string `json:"off-heap,omitempty"`
}

type TaskManagerNetWorkFlinkResource struct {
	// taskmanager.memory.network.fraction ， 与 max min 仅其中一个配置生效，默认走max min
	Fraction string `json:"fraction,omitempty"`
	// taskmanager.memory.network.min， 与 max min 仅其中一个配置生效，默认走max min
	Min string `json:"min,omitempty"`
	// taskmanager.memory.network.max ， 与 max min 仅其中一个配置生效，默认走max min
	Max string `json:"max,omitempty"`
}

type TaskManagerManagedFlinkResource struct {
	// taskmanager.memory.managed.fraction ， 与 max min 仅其中一个配置生效，默认走max min
	Fraction string `json:"fraction,omitempty"`
	// taskmanager.memory.managed.min， 与 max min 仅其中一个配置生效，默认走max min
	Min string `json:"min,omitempty"`
	// taskmanager.memory.managed.max ， 与 max min 仅其中一个配置生效，默认走max min
	Max string `json:"max,omitempty"`
}

type FlinkS3 struct {
	// s3.endpoint
	//+kubebuilder:validation:MinLength=1
	EndPoint string `json:"endPoint,omitempty"`
	// s3.access-key
	//+kubebuilder:validation:MinLength=1
	AccessKey string `json:"accessKey,omitempty"`
	// s3.secret-key
	//+kubebuilder:validation:MinLength=1
	SecretKey string `json:"secretKey,omitempty"`
	// 部署 bucket
	//+kubebuilder:validation:MinLength=1
	Bucket string `json:"bucket"`
}

type FlinkHA struct {
	// ha 类型，允许值 zookeeper  kubernetes none
	//+kubebuilder:validation:Enum={zookeeper,kubernetes,none}
	Typ FlinkHAType `json:"type,omitempty"`

	// 仅 zookeeper ha 生效，配置 zk 地址
	//+nullable
	Quorum string `json:"quorum,omitempty"`
	// 仅 zookeeper ha 生效，配置 ha 路径前缀
	//+nullable
	Path string `json:"path,omitempty"`
}

type FlinkHAType string

const (
	ZKHA        FlinkHAType = "zookeeper"
	CONFIGMAPHA FlinkHAType = "kubernetes"
	NONE        FlinkHAType = "none"
)

const (
	RequiredDuringScheduling  string = `Required`
	PreferredDuringScheduling string = `Preferred`
	NoneScheduling            string = `None`
)

type FlinkConfig struct {
	// 对应 $FLINK_HOME/conf/flink-conf.yaml ，不写使用镜像内预制配置文件
	//+nullable
	FlinkConf string `json:"flink-conf.yaml,omitempty"`
	// 对应 $FLINK_HOME/conf/log4j.properties ，不写使用镜像内预制配置文件
	//+nullable
	Log4j string `json:"log4j.properties,omitempty"`
	// 对应 $FLINK_HOME/conf/logback.xml ，不写使用镜像内预制配置文件
	//+nullable
	LogBack string `json:"logback.xml,omitempty"`
}

// FlinkSessionStatus defines the observed state of FlinkSession
type FlinkSessionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Ready bool  `json:"ready,omitempty"`
	Port  int32 `json:"port,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FlinkSession is the Schema for the flinksessions API
// +kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready"
// +kubebuilder:printcolumn:name="Port",type="integer",JSONPath=".status.port"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type FlinkSession struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FlinkSessionSpec   `json:"spec,omitempty"`
	Status FlinkSessionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FlinkSessionList contains a list of FlinkSession
type FlinkSessionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FlinkSession `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FlinkSession{}, &FlinkSessionList{})
}

//+kubebuilder:docs-gen:collapse=Root Object Definitions
