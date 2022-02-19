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

	//+kubebuilder:validation:MinLength=0
	// CPU is
	CPU string `json:"cpu,omitempty"`
	//+kubebuilder:validation:MinLength=0
	// Memory is
	Memory string `json:"memory,omitempty"`
	//+kubebuilder:validation:MinLength=0
	// Image is
	Image string `json:"image,omitempty"`
	//+kubebuilder:validation:MinLength=0
	// BootCmd is
	BootCmd string `json:"bootCmd,omitempty"`
}

// FlinkSessionStatus defines the observed state of FlinkSession
type FlinkSessionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Ready bool `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FlinkSession is the Schema for the flinksessions API
// +kubebuilder:printcolumn:name="CPU",type="string",JSONPath=".spec.cpu"
// +kubebuilder:printcolumn:name="Memory",type="string",JSONPath=".spec.memory"
// +kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready"
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