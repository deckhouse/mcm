/*
Copyright (c) 2017 SAP SE or an SAP affiliate company. All rights reserved.

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

// Package controller is used to provide the core functionalities of machine-controller-manager
package controller

import (
	"k8s.io/apimachinery/pkg/labels"

	"github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
)

// existsMachineClassForSecret checks for any machineClass
// referring to the passed secret object
func (c *controller) existsMachineClassForSecret(name string) (bool, error) {
	openStackMachineClasses, err := c.findOpenStackMachineClassForSecret(name)
	if err != nil {
		return false, err
	}

	gcpMachineClasses, err := c.findGCPMachineClassForSecret(name)
	if err != nil {
		return false, err
	}

	azureMachineClasses, err := c.findAzureMachineClassForSecret(name)
	if err != nil {
		return false, err
	}

	awsMachineClasses, err := c.findAWSMachineClassForSecret(name)
	if err != nil {
		return false, err
	}

	alicloudMachineClasses, err := c.findAlicloudMachineClassForSecret(name)
	if err != nil {
		return false, err
	}

	packetMachineClasses, err := c.findPacketMachineClassForSecret(name)
	if err != nil {
		return false, err
	}

	vsphereMachineClasses, err := c.findVsphereMachineClassForSecret(name)
	if err != nil {
		return false, err
	}

	yandexMachineClasses, err := c.findYandexMachineClassForSecret(name)
	if err != nil {
		return false, err
	}

	if len(openStackMachineClasses) == 0 &&
		len(gcpMachineClasses) == 0 &&
		len(azureMachineClasses) == 0 &&
		len(packetMachineClasses) == 0 &&
		len(alicloudMachineClasses) == 0 &&
		len(awsMachineClasses) == 0 &&
		len(vsphereMachineClasses) == 0 &&
		len(yandexMachineClasses) == 0 {
		return false, nil
	}

	return true, nil
}

// findOpenStackMachineClassForSecret returns the set of
// openStackMachineClasses referring to the passed secret
func (c *controller) findOpenStackMachineClassForSecret(name string) ([]*v1alpha1.OpenStackMachineClass, error) {
	machineClasses, err := c.openStackMachineClassLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	var filtered []*v1alpha1.OpenStackMachineClass
	for _, machineClass := range machineClasses {
		if (machineClass.Spec.SecretRef != nil && machineClass.Spec.SecretRef.Name == name) ||
			(machineClass.Spec.CredentialsSecretRef != nil && machineClass.Spec.CredentialsSecretRef.Name == name) {
			filtered = append(filtered, machineClass)
		}
	}
	return filtered, nil
}

// findGCPClassForSecret returns the set of
// GCPMachineClasses referring to the passed secret
func (c *controller) findGCPMachineClassForSecret(name string) ([]*v1alpha1.GCPMachineClass, error) {
	machineClasses, err := c.gcpMachineClassLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	var filtered []*v1alpha1.GCPMachineClass
	for _, machineClass := range machineClasses {
		if (machineClass.Spec.SecretRef != nil && machineClass.Spec.SecretRef.Name == name) ||
			(machineClass.Spec.CredentialsSecretRef != nil && machineClass.Spec.CredentialsSecretRef.Name == name) {
			filtered = append(filtered, machineClass)
		}
	}
	return filtered, nil
}

// findAzureClassForSecret returns the set of
// AzureMachineClasses referring to the passed secret
func (c *controller) findAzureMachineClassForSecret(name string) ([]*v1alpha1.AzureMachineClass, error) {
	machineClasses, err := c.azureMachineClassLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	var filtered []*v1alpha1.AzureMachineClass
	for _, machineClass := range machineClasses {
		if (machineClass.Spec.SecretRef != nil && machineClass.Spec.SecretRef.Name == name) ||
			(machineClass.Spec.CredentialsSecretRef != nil && machineClass.Spec.CredentialsSecretRef.Name == name) {
			filtered = append(filtered, machineClass)
		}
	}
	return filtered, nil
}

// findAlicloudClassForSecret returns the set of
// AlicloudMachineClasses referring to the passed secret
func (c *controller) findAlicloudMachineClassForSecret(name string) ([]*v1alpha1.AlicloudMachineClass, error) {
	machineClasses, err := c.alicloudMachineClassLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	var filtered []*v1alpha1.AlicloudMachineClass
	for _, machineClass := range machineClasses {
		if (machineClass.Spec.SecretRef != nil && machineClass.Spec.SecretRef.Name == name) ||
			(machineClass.Spec.CredentialsSecretRef != nil && machineClass.Spec.CredentialsSecretRef.Name == name) {
			filtered = append(filtered, machineClass)
		}
	}
	return filtered, nil
}

// findAWSClassForSecret returns the set of
// AWSMachineClasses referring to the passed secret
func (c *controller) findAWSMachineClassForSecret(name string) ([]*v1alpha1.AWSMachineClass, error) {
	machineClasses, err := c.awsMachineClassLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	var filtered []*v1alpha1.AWSMachineClass
	for _, machineClass := range machineClasses {
		if (machineClass.Spec.SecretRef != nil && machineClass.Spec.SecretRef.Name == name) ||
			(machineClass.Spec.CredentialsSecretRef != nil && machineClass.Spec.CredentialsSecretRef.Name == name) {
			filtered = append(filtered, machineClass)
		}
	}
	return filtered, nil
}

// findPacketClassForSecret returns the set of
// PacketMachineClasses referring to the passed secret
func (c *controller) findPacketMachineClassForSecret(name string) ([]*v1alpha1.PacketMachineClass, error) {
	machineClasses, err := c.packetMachineClassLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	var filtered []*v1alpha1.PacketMachineClass
	for _, machineClass := range machineClasses {
		if (machineClass.Spec.SecretRef != nil && machineClass.Spec.SecretRef.Name == name) ||
			(machineClass.Spec.CredentialsSecretRef != nil && machineClass.Spec.CredentialsSecretRef.Name == name) {
			filtered = append(filtered, machineClass)
		}
	}
	return filtered, nil
}

// findVsphereClassForSecret returns the set of
// VsphereMachineClasses referring to the passed secret
func (c *controller) findVsphereMachineClassForSecret(name string) ([]*v1alpha1.VsphereMachineClass, error) {
	machineClasses, err := c.vsphereMachineClassLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	var filtered []*v1alpha1.VsphereMachineClass
	for _, machineClass := range machineClasses {
		if machineClass.Spec.SecretRef.Name == name {
			filtered = append(filtered, machineClass)
		}
	}
	return filtered, nil
}

// findYandexClassForSecret returns the set of
// YandexMachineClasses referring to the passed secret
func (c *controller) findYandexMachineClassForSecret(name string) ([]*v1alpha1.YandexMachineClass, error) {
	machineClasses, err := c.yandexMachineClassLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	var filtered []*v1alpha1.YandexMachineClass
	for _, machineClass := range machineClasses {
		if machineClass.Spec.SecretRef.Name == name {
			filtered = append(filtered, machineClass)
		}
	}
	return filtered, nil
}
