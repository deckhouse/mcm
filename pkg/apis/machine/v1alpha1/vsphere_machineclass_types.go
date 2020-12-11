// Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// WARNING!
// IF YOU MODIFY ANY OF THE TYPES HERE COPY THEM TO ../types.go
// AND RUN `make generate`

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// VsphereHost is a constant for a key name that is part of the Vsphere cloud credentials.
	VsphereHost string = "host"
	// VsphereUsername is a constant for a key name that is part of the Vsphere cloud credentials.
	VsphereUsername string = "username"
	// VspherePassword is a constant for a key name that is part of the Vsphere cloud credentials.
	VspherePassword string = "password"
	// VsphereInsecure is a constant for a key name that is part of the Vsphere cloud credentials.
	VsphereInsecure string = "insecure"
	// VsphereRegionTagCategory is a constant for a key name that is part of the Vsphere cloud credentials.
	VsphereRegionTagCategory string = "regionTagCategory"
	// VsphereZoneTagCategory is a constant for a key name that is part of the Vsphere cloud credentials.
	VsphereZoneTagCategory string = "zoneTagCategory"
	// VsphereClusterNameTagCategory is a constant for a key name that is part of the Vsphere cloud credentials.
	VsphereClusterNameTagCategory string = "clusterNameTagCategory"
	// VsphereNodeRoleTagCategory is a constant for a key name that is part of the Vsphere cloud credentials.
	VsphereNodeRoleTagCategory string = "nodeRoleTagCategory"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="VM Size",type=string,JSONPath=`.spec.properties.hardwareProfile.vmSize`
// +kubebuilder:printcolumn:name="Location",type=string,JSONPath=`.spec.location`,priority=1
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.\nPopulated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata"

// VsphereMachineClass TODO
type VsphereMachineClass struct {
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	metav1.TypeMeta `json:",inline"`

	// +optional
	Spec VsphereMachineClassSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// VsphereMachineClassList is a collection of VsphereMachineClass.
type VsphereMachineClassList struct {
	// +optional
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// +optional
	Items []VsphereMachineClass `json:"items"`
}

// VsphereMachineClassSpec is the specification of a cluster.
type VsphereMachineClassSpec struct {
	NumCPUs              uint64            `json:"numCPUs,omitempty"`
	Memory               uint64            `json:"memory,omitempty"`
	Region               string            `json:"region,omitempty"`
	Zone                 string            `json:"zone,omitempty"`
	ClusterNameTag       string            `json:"clusterNameTag,omitempty"`
	NodeRoleTag          string            `json:"nodeRoleTag,omitempty"`
	Template             string            `json:"template,omitempty"`
	ResourcePool         string            `json:"resourcePool,omitempty"`
	VirtualMachineFolder string            `json:"virtualMachineFolder,omitempty"`
	MainNetwork          string            `json:"mainNetwork,omitempty"`
	AdditionalNetworks   []string          `json:"additionalNetworks,omitempty"`
	Datastore            string            `json:"datastore,omitempty"`
	RootDiskSize         uint64            `json:"rootDiskSize,omitempty"`
	DisableTimesync      bool              `json:"disableTimesync,omitempty"`
	SshKeys              []string          `json:"sshKeys,omitempty"`
	ExtraConfig          map[string]string `json:"extraConfig,omitempty"`

	RuntimeOptions VsphereMachineClassSpecRuntimeOptions `json:"runtimeOptions,omitempty"`

	UserData             string                  `json:"userData,omitempty"`
	SecretRef            *corev1.SecretReference `json:"secretRef,omitempty"`
	CredentialsSecretRef *corev1.SecretReference `json:"credentialsSecretRef,omitempty"`
}
type VsphereMachineClassSpecRuntimeOptions struct {
	NestedHardwareVirtualization bool `json:"nestedHardwareVirtualization,omitempty"`

	ResourceAllocationInfo VsphereMachineClassSpecResourceAllocationInfo `json:"resourceAllocationInfo,omitempty"`
}

type VsphereMachineClassSpecResourceAllocationInfo struct {
	CpuShares         *int32 `json:"cpuShares,omitempty"`
	CpuLimit          *int64 `json:"cpuLimit,omitempty"`
	CpuReservation    *int64 `json:"cpuReservation,omitempty"`
	MemoryShares      *int32 `json:"memoryShares,omitempty"`
	MemoryLimit       *int64 `json:"memoryLimit,omitempty"`
	MemoryReservation *int64 `json:"memoryReservation,omitempty"`
}
