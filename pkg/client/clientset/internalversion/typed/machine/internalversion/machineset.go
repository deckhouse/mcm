/*
Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

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

// Code generated by client-gen. DO NOT EDIT.

package internalversion

import (
	"time"

	machine "github.com/gardener/machine-controller-manager/pkg/apis/machine"
	scheme "github.com/gardener/machine-controller-manager/pkg/client/clientset/internalversion/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MachineSetsGetter has a method to return a MachineSetInterface.
// A group's client should implement this interface.
type MachineSetsGetter interface {
	MachineSets(namespace string) MachineSetInterface
}

// MachineSetInterface has methods to work with MachineSet resources.
type MachineSetInterface interface {
	Create(*machine.MachineSet) (*machine.MachineSet, error)
	Update(*machine.MachineSet) (*machine.MachineSet, error)
	UpdateStatus(*machine.MachineSet) (*machine.MachineSet, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*machine.MachineSet, error)
	List(opts v1.ListOptions) (*machine.MachineSetList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *machine.MachineSet, err error)
	MachineSetExpansion
}

// machineSets implements MachineSetInterface
type machineSets struct {
	client rest.Interface
	ns     string
}

// newMachineSets returns a MachineSets
func newMachineSets(c *MachineClient, namespace string) *machineSets {
	return &machineSets{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the machineSet, and returns the corresponding machineSet object, and an error if there is any.
func (c *machineSets) Get(name string, options v1.GetOptions) (result *machine.MachineSet, err error) {
	result = &machine.MachineSet{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("machinesets").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MachineSets that match those selectors.
func (c *machineSets) List(opts v1.ListOptions) (result *machine.MachineSetList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &machine.MachineSetList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("machinesets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested machineSets.
func (c *machineSets) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("machinesets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a machineSet and creates it.  Returns the server's representation of the machineSet, and an error, if there is any.
func (c *machineSets) Create(machineSet *machine.MachineSet) (result *machine.MachineSet, err error) {
	result = &machine.MachineSet{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("machinesets").
		Body(machineSet).
		Do().
		Into(result)
	return
}

// Update takes the representation of a machineSet and updates it. Returns the server's representation of the machineSet, and an error, if there is any.
func (c *machineSets) Update(machineSet *machine.MachineSet) (result *machine.MachineSet, err error) {
	result = &machine.MachineSet{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("machinesets").
		Name(machineSet.Name).
		Body(machineSet).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *machineSets) UpdateStatus(machineSet *machine.MachineSet) (result *machine.MachineSet, err error) {
	result = &machine.MachineSet{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("machinesets").
		Name(machineSet.Name).
		SubResource("status").
		Body(machineSet).
		Do().
		Into(result)
	return
}

// Delete takes name of the machineSet and deletes it. Returns an error if one occurs.
func (c *machineSets) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("machinesets").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *machineSets) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("machinesets").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched machineSet.
func (c *machineSets) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *machine.MachineSet, err error) {
	result = &machine.MachineSet{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("machinesets").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
