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
	YandexFolderID           string = "folderID"
	YandexServiceAccountJSON string = "serviceAccountJSON"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="VM Size",type=string,JSONPath=`.spec.properties.hardwareProfile.vmSize`
// +kubebuilder:printcolumn:name="Location",type=string,JSONPath=`.spec.location`,priority=1
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.\nPopulated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata"

// YandexMachineClass TODO
type YandexMachineClass struct {
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	metav1.TypeMeta `json:",inline"`

	// +optional
	Spec YandexMachineClassSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// YandexMachineClassList is a collection of PacketMachineClasses.
type YandexMachineClassList struct {
	// +optional
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// +optional
	Items []YandexMachineClass `json:"items"`
}

type YandexMachineClassSpec struct {
	Labels     map[string]string `json:"labels,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	RegionID   string            `json:"regionID,omitempty"`
	ZoneID     string            `json:"zoneID,omitempty"`
	PlatformID string            `json:"platformID,omitempty"`

	ResourcesSpec         YandexMachineClassSpecResourcesSpec           `json:"resourcesSpec,omitempty"`
	BootDiskSpec          YandexMachineClassSpecBootDiskSpec            `json:"bootDiskSpec,omitempty"`
	NetworkInterfaceSpecs []YandexMachineClassSpecNetworkInterfaceSpecs `json:"networkInterfaceSpecs,omitempty"`
	SchedulingPolicy      YandexMachineClassSpecSchedulingPolicy        `json:"schedulingPolicy,omitempty"`

	SecretRef            *corev1.SecretReference `json:"secretRef,omitempty"`
	CredentialsSecretRef *corev1.SecretReference `json:"credentialsSecretRef,omitempty"`
}

type YandexMachineClassSpecResourcesSpec struct {
	Cores        int64 `json:"cores,omitempty"`
	CoreFraction int64 `json:"coreFraction,omitempty"`
	Memory       int64 `json:"memory,omitempty"`
	GPUs         int64 `json:"gpus,omitempty"`
}

type YandexMachineClassSpecBootDiskSpec struct {
	AutoDelete bool   `json:"autoDelete,omitempty"`
	TypeID     string `json:"typeID,omitempty"`
	Size       int64  `json:"size,omitempty"`
	ImageID    string `json:"imageID,omitempty"`
}

type YandexMachineClassSpecNetworkInterfaceSpecs struct {
	SubnetID              string `json:"subnetID,omitempty"`
	AssignPublicIPAddress bool   `json:"assignPublicIPAddress,omitempty"`
}

type YandexMachineClassSpecSchedulingPolicy struct {
	Preemptible bool `json:"preemptible,omitempty"`
}
