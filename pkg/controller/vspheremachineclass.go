/*
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

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	"github.com/gardener/machine-controller-manager/pkg/apis/machine"
	"github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	"github.com/gardener/machine-controller-manager/pkg/apis/machine/validation"
)

// VsphereMachineClassKind is used to identify the machineClassKind as Vsphere
const VsphereMachineClassKind = "VsphereMachineClass"

func (c *controller) machineDeploymentToVsphereMachineClassDelete(obj interface{}) {
	machineDeployment, ok := obj.(*v1alpha1.MachineDeployment)
	if machineDeployment == nil || !ok {
		return
	}
	if machineDeployment.Spec.Template.Spec.Class.Kind == VsphereMachineClassKind {
		c.vsphereMachineClassQueue.Add(machineDeployment.Spec.Template.Spec.Class.Name)
	}
}

func (c *controller) machineSetToVsphereMachineClassDelete(obj interface{}) {
	machineSet, ok := obj.(*v1alpha1.MachineSet)
	if machineSet == nil || !ok {
		return
	}
	if machineSet.Spec.Template.Spec.Class.Kind == VsphereMachineClassKind {
		c.vsphereMachineClassQueue.Add(machineSet.Spec.Template.Spec.Class.Name)
	}
}

func (c *controller) machineToVsphereMachineClassDelete(obj interface{}) {
	machine, ok := obj.(*v1alpha1.Machine)
	if machine == nil || !ok {
		return
	}
	if machine.Spec.Class.Kind == VsphereMachineClassKind {
		c.vsphereMachineClassQueue.Add(machine.Spec.Class.Name)
	}
}

func (c *controller) vsphereMachineClassAdd(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		klog.Errorf("Couldn't get key for object %+v: %v", obj, err)
		return
	}
	c.vsphereMachineClassQueue.Add(key)
}

func (c *controller) vsphereMachineClassUpdate(oldObj, newObj interface{}) {
	old, ok := oldObj.(*v1alpha1.VsphereMachineClass)
	if old == nil || !ok {
		return
	}
	new, ok := newObj.(*v1alpha1.VsphereMachineClass)
	if new == nil || !ok {
		return
	}

	c.vsphereMachineClassAdd(newObj)
}

// reconcileClusterVsphereMachineClassKey reconciles a VsphereMachineClass due to controller resync
// or an event on the vsphereMachineClass.
func (c *controller) reconcileClusterVsphereMachineClassKey(key string) error {
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	class, err := c.vsphereMachineClassLister.VsphereMachineClasses(c.namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.Infof("%s %q: Not doing work because it has been deleted", VsphereMachineClassKind, key)
		return nil
	}
	if err != nil {
		klog.Infof("%s %q: Unable to retrieve object from store: %v", VsphereMachineClassKind, key, err)
		return err
	}

	return c.reconcileClusterVsphereMachineClass(class)
}

func (c *controller) reconcileClusterVsphereMachineClass(class *v1alpha1.VsphereMachineClass) error {
	klog.V(4).Info("Start Reconciling vspheremachineclass: ", class.Name)
	defer func() {
		c.enqueueVsphereMachineClassAfter(class, 10*time.Minute)
		klog.V(4).Info("Stop Reconciling vspheremachineclass: ", class.Name)
	}()

	internalClass := &machine.VsphereMachineClass{}
	err := c.internalExternalScheme.Convert(class, internalClass, nil)
	if err != nil {
		return err
	}
	// TODO this should be put in own API server
	validationerr := validation.ValidateVsphereMachineClass(internalClass)
	if validationerr.ToAggregate() != nil && len(validationerr.ToAggregate().Errors()) > 0 {
		klog.V(2).Infof("Validation of %s failed %s", VsphereMachineClassKind, validationerr.ToAggregate().Error())
		return nil
	}

	// Manipulate finalizers
	if class.DeletionTimestamp == nil {
		err = c.addVsphereMachineClassFinalizers(class)
		if err != nil {
			return err
		}
	}

	machines, err := c.findMachinesForClass(VsphereMachineClassKind, class.Name)
	if err != nil {
		return err
	}

	if class.DeletionTimestamp != nil {
		if finalizers := sets.NewString(class.Finalizers...); !finalizers.Has(DeleteFinalizerName) {
			return nil
		}

		machineDeployments, err := c.findMachineDeploymentsForClass(VsphereMachineClassKind, class.Name)
		if err != nil {
			return err
		}
		machineSets, err := c.findMachineSetsForClass(VsphereMachineClassKind, class.Name)
		if err != nil {
			return err
		}
		if len(machineDeployments) == 0 && len(machineSets) == 0 && len(machines) == 0 {
			return c.deleteVsphereMachineClassFinalizers(class)
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

func (c *controller) addVsphereMachineClassFinalizers(class *v1alpha1.VsphereMachineClass) error {
	clone := class.DeepCopy()

	if finalizers := sets.NewString(clone.Finalizers...); !finalizers.Has(DeleteFinalizerName) {
		finalizers.Insert(DeleteFinalizerName)
		return c.updateVsphereMachineClassFinalizers(clone, finalizers.List())
	}
	return nil
}

func (c *controller) deleteVsphereMachineClassFinalizers(class *v1alpha1.VsphereMachineClass) error {
	clone := class.DeepCopy()

	if finalizers := sets.NewString(clone.Finalizers...); finalizers.Has(DeleteFinalizerName) {
		finalizers.Delete(DeleteFinalizerName)
		return c.updateVsphereMachineClassFinalizers(clone, finalizers.List())
	}
	return nil
}

func (c *controller) updateVsphereMachineClassFinalizers(class *v1alpha1.VsphereMachineClass, finalizers []string) error {
	// Get the latest version of the class so that we can avoid conflicts
	class, err := c.controlMachineClient.VsphereMachineClasses(class.Namespace).Get(class.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	clone := class.DeepCopy()
	clone.Finalizers = finalizers
	_, err = c.controlMachineClient.VsphereMachineClasses(class.Namespace).Update(clone)
	if err != nil {
		klog.Warning("Updating VsphereMachineClass failed, retrying. ", class.Name, err)
		return err
	}
	klog.V(3).Infof("Successfully added/removed finalizer on the vspheremachineclass %q", class.Name)
	return err
}

func (c *controller) enqueueVsphereMachineClassAfter(obj interface{}, after time.Duration) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	c.openStackMachineClassQueue.AddAfter(key, after)
}
