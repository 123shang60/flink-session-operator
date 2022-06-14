//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkConfig) DeepCopyInto(out *FlinkConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkConfig.
func (in *FlinkConfig) DeepCopy() *FlinkConfig {
	if in == nil {
		return nil
	}
	out := new(FlinkConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkHA) DeepCopyInto(out *FlinkHA) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkHA.
func (in *FlinkHA) DeepCopy() *FlinkHA {
	if in == nil {
		return nil
	}
	out := new(FlinkHA)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkResource) DeepCopyInto(out *FlinkResource) {
	*out = *in
	out.JobManager = in.JobManager
	out.TaskManager = in.TaskManager
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkResource.
func (in *FlinkResource) DeepCopy() *FlinkResource {
	if in == nil {
		return nil
	}
	out := new(FlinkResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkS3) DeepCopyInto(out *FlinkS3) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkS3.
func (in *FlinkS3) DeepCopy() *FlinkS3 {
	if in == nil {
		return nil
	}
	out := new(FlinkS3)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkSecurity) DeepCopyInto(out *FlinkSecurity) {
	*out = *in
	if in.Kerberos != nil {
		in, out := &in.Kerberos, &out.Kerberos
		*out = new(Kerberos)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkSecurity.
func (in *FlinkSecurity) DeepCopy() *FlinkSecurity {
	if in == nil {
		return nil
	}
	out := new(FlinkSecurity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkSession) DeepCopyInto(out *FlinkSession) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkSession.
func (in *FlinkSession) DeepCopy() *FlinkSession {
	if in == nil {
		return nil
	}
	out := new(FlinkSession)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FlinkSession) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkSessionList) DeepCopyInto(out *FlinkSessionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FlinkSession, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkSessionList.
func (in *FlinkSessionList) DeepCopy() *FlinkSessionList {
	if in == nil {
		return nil
	}
	out := new(FlinkSessionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FlinkSessionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkSessionSpec) DeepCopyInto(out *FlinkSessionSpec) {
	*out = *in
	if in.ImageSecret != nil {
		in, out := &in.ImageSecret, &out.ImageSecret
		*out = new(string)
		**out = **in
	}
	out.Resource = in.Resource
	out.S3 = in.S3
	out.HA = in.HA
	out.Config = in.Config
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]corev1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]corev1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Security.DeepCopyInto(&out.Security)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkSessionSpec.
func (in *FlinkSessionSpec) DeepCopy() *FlinkSessionSpec {
	if in == nil {
		return nil
	}
	out := new(FlinkSessionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlinkSessionStatus) DeepCopyInto(out *FlinkSessionStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlinkSessionStatus.
func (in *FlinkSessionStatus) DeepCopy() *FlinkSessionStatus {
	if in == nil {
		return nil
	}
	out := new(FlinkSessionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobManagerFlinkResource) DeepCopyInto(out *JobManagerFlinkResource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobManagerFlinkResource.
func (in *JobManagerFlinkResource) DeepCopy() *JobManagerFlinkResource {
	if in == nil {
		return nil
	}
	out := new(JobManagerFlinkResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Kerberos) DeepCopyInto(out *Kerberos) {
	*out = *in
	if in.UseTicketCache != nil {
		in, out := &in.UseTicketCache, &out.UseTicketCache
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Kerberos.
func (in *Kerberos) DeepCopy() *Kerberos {
	if in == nil {
		return nil
	}
	out := new(Kerberos)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskManagerFlinkResource) DeepCopyInto(out *TaskManagerFlinkResource) {
	*out = *in
	out.Framework = in.Framework
	out.Task = in.Task
	out.NetWork = in.NetWork
	out.Managed = in.Managed
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskManagerFlinkResource.
func (in *TaskManagerFlinkResource) DeepCopy() *TaskManagerFlinkResource {
	if in == nil {
		return nil
	}
	out := new(TaskManagerFlinkResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskManagerFrameworkFlinkResource) DeepCopyInto(out *TaskManagerFrameworkFlinkResource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskManagerFrameworkFlinkResource.
func (in *TaskManagerFrameworkFlinkResource) DeepCopy() *TaskManagerFrameworkFlinkResource {
	if in == nil {
		return nil
	}
	out := new(TaskManagerFrameworkFlinkResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskManagerManagedFlinkResource) DeepCopyInto(out *TaskManagerManagedFlinkResource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskManagerManagedFlinkResource.
func (in *TaskManagerManagedFlinkResource) DeepCopy() *TaskManagerManagedFlinkResource {
	if in == nil {
		return nil
	}
	out := new(TaskManagerManagedFlinkResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskManagerNetWorkFlinkResource) DeepCopyInto(out *TaskManagerNetWorkFlinkResource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskManagerNetWorkFlinkResource.
func (in *TaskManagerNetWorkFlinkResource) DeepCopy() *TaskManagerNetWorkFlinkResource {
	if in == nil {
		return nil
	}
	out := new(TaskManagerNetWorkFlinkResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskManagerTaskFlinkResource) DeepCopyInto(out *TaskManagerTaskFlinkResource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskManagerTaskFlinkResource.
func (in *TaskManagerTaskFlinkResource) DeepCopy() *TaskManagerTaskFlinkResource {
	if in == nil {
		return nil
	}
	out := new(TaskManagerTaskFlinkResource)
	in.DeepCopyInto(out)
	return out
}
