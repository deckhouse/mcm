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
	"time"

	"k8s.io/klog"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"

	"github.com/gardener/machine-controller-manager/pkg/apis/machine"
	"github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
)

// YandexMachineClassKind is used to identify the machineClassKind as Yandex
const YandexMachineClassKind = "YandexMachineClass"

func (c *controller) machineDeploymentToYandexMachineClassDelete(obj interface{}) {
	machineDeployment, ok := obj.(*v1alpha1.MachineDeployment)
	if machineDeployment == nil || !ok {
		return
	}
	if machineDeployment.Spec.Template.Spec.Class.Kind == YandexMachineClassKind {
		c.yandexMachineClassQueue.Add(machineDeployment.Spec.Template.Spec.Class.Name)
	}
}

func (c *controller) machineSetToYandexMachineClassDelete(obj interface{}) {
	machineSet, ok := obj.(*v1alpha1.MachineSet)
	if machineSet == nil || !ok {
		return
	}
	if machineSet.Spec.Template.Spec.Class.Kind == YandexMachineClassKind {
		c.yandexMachineClassQueue.Add(machineSet.Spec.Template.Spec.Class.Name)
	}
}

func (c *controller) machineToYandexMachineClassDelete(obj interface{}) {
	machine, ok := obj.(*v1alpha1.Machine)
	if machine == nil || !ok {
		return
	}
	if machine.Spec.Class.Kind == YandexMachineClassKind {
		c.yandexMachineClassQueue.Add(machine.Spec.Class.Name)
	}
}

func (c *controller) yandexMachineClassAdd(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		klog.Errorf("Couldn't get key for object %+v: %v", obj, err)
		return
	}
	c.yandexMachineClassQueue.Add(key)
}

func (c *controller) yandexMachineClassUpdate(oldObj, newObj interface{}) {
	old, ok := oldObj.(*v1alpha1.YandexMachineClass)
	if old == nil || !ok {
		return
	}
	new, ok := newObj.(*v1alpha1.YandexMachineClass)
	if new == nil || !ok {
		return
	}

	c.yandexMachineClassAdd(newObj)
}

// reconcileClusterYandexMachineClassKey reconciles an YandexMachineClass due to controller resync
// or an event on the yandexMachineClass.
func (c *controller) reconcileClusterYandexMachineClassKey(key string) error {
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	class, err := c.yandexMachineClassLister.YandexMachineClasses(c.namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.Infof("%s %q: Not doing work because it has been deleted", YandexMachineClassKind, key)
		return nil
	}
	if err != nil {
		klog.Infof("%s %q: Unable to retrieve object from store: %v", YandexMachineClassKind, key, err)
		return err
	}

	return c.reconcileClusterYandexMachineClass(class)
}

func (c *controller) reconcileClusterYandexMachineClass(class *v1alpha1.YandexMachineClass) error {
	klog.V(4).Info("Start Reconciling yandexmachineclass: ", class.Name)
	defer func() {
		c.enqueueYandexMachineClassAfter(class, 10*time.Minute)
		klog.V(4).Info("Stop Reconciling yandexmachineclass: ", class.Name)
	}()

	internalClass := &machine.YandexMachineClass{}
	err := c.internalExternalScheme.Convert(class, internalClass, nil)
	if err != nil {
		return err
	}

	// Manipulate finalizers
	if class.DeletionTimestamp == nil {
		err = c.addYandexMachineClassFinalizers(class)
		if err != nil {
			return err
		}
	}

	machines, err := c.findMachinesForClass(YandexMachineClassKind, class.Name)
	if err != nil {
		return err
	}

	if class.DeletionTimestamp != nil {
		if finalizers := sets.NewString(class.Finalizers...); !finalizers.Has(DeleteFinalizerName) {
			return nil
		}

		machineDeployments, err := c.findMachineDeploymentsForClass(YandexMachineClassKind, class.Name)
		if err != nil {
			return err
		}
		machineSets, err := c.findMachineSetsForClass(YandexMachineClassKind, class.Name)
		if err != nil {
			return err
		}
		if len(machineDeployments) == 0 && len(machineSets) == 0 && len(machines) == 0 {
			return c.deleteYandexMachineClassFinalizers(class)
		}

		klog.V(3).Infof("Cannot remove finalizer of %s because still Machine[s|Sets|Deployments] are referencing it", class.Name)
		return nil
	}

	for _, machine := range machines {
		c.addMachine(machine)
	}
	return nil
}

/*
	SECTION
	Manipulate Finalizers
*/

func (c *controller) addYandexMachineClassFinalizers(class *v1alpha1.YandexMachineClass) error {
	clone := class.DeepCopy()

	if finalizers := sets.NewString(clone.Finalizers...); !finalizers.Has(DeleteFinalizerName) {
		finalizers.Insert(DeleteFinalizerName)
		return c.updateYandexMachineClassFinalizers(clone, finalizers.List())
	}
	return nil
}

func (c *controller) deleteYandexMachineClassFinalizers(class *v1alpha1.YandexMachineClass) error {
	clone := class.DeepCopy()

	if finalizers := sets.NewString(clone.Finalizers...); finalizers.Has(DeleteFinalizerName) {
		finalizers.Delete(DeleteFinalizerName)
		return c.updateYandexMachineClassFinalizers(clone, finalizers.List())
	}
	return nil
}

func (c *controller) updateYandexMachineClassFinalizers(class *v1alpha1.YandexMachineClass, finalizers []string) error {
	// Get the latest version of the class so that we can avoid conflicts
	class, err := c.controlMachineClient.YandexMachineClasses(class.Namespace).Get(class.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	clone := class.DeepCopy()
	clone.Finalizers = finalizers
	_, err = c.controlMachineClient.YandexMachineClasses(class.Namespace).Update(clone)
	if err != nil {
		klog.Warning("Updating YandexMachineClass failed, retrying. ", class.Name, err)
		return err
	}
	klog.V(3).Infof("Successfully added/removed finalizer on the yandexmachineclass %q", class.Name)
	return err
}

func (c *controller) enqueueYandexMachineClassAfter(obj interface{}, after time.Duration) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	c.yandexMachineClassQueue.AddAfter(key, after)
}
