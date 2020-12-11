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

// Package validation is used to validate all the machine CRD objects
package validation

import (
	"regexp"

	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	utilvalidation "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/gardener/machine-controller-manager/pkg/apis/machine"
)

const vsphereNameFmt string = `[-a-z0-9]+`
const vsphereNameMaxLength int = 63

var vsphereNameRegexp = regexp.MustCompile("^" + vsphereNameFmt + "$")

// vsphereValidateName is the validation function for object names.
func vsphereValidateName(value string, prefix bool) []string {
	var errs []string
	if len(value) > nameMaxLength {
		errs = append(errs, utilvalidation.MaxLenError(vsphereNameMaxLength))
	}
	if !vsphereNameRegexp.MatchString(value) {
		errs = append(errs, utilvalidation.RegexError(vsphereNameFmt, "name-40d-0983-1b89"))
	}

	return errs
}

// ValidateVsphereMachineClass validates a VsphereMachineClass and returns a list of errors.
func ValidateVsphereMachineClass(VsphereMachineClass *machine.VsphereMachineClass) field.ErrorList {
	return internalValidateVsphereMachineClass(VsphereMachineClass)
}

func internalValidateVsphereMachineClass(VsphereMachineClass *machine.VsphereMachineClass) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, apivalidation.ValidateObjectMeta(&VsphereMachineClass.ObjectMeta, true, /*namespace*/
		vsphereValidateName,
		field.NewPath("metadata"))...)

	allErrs = append(allErrs, validateVsphereMachineClassSpec(&VsphereMachineClass.Spec, field.NewPath("spec"))...)
	return allErrs
}

func validateVsphereMachineClassSpec(spec *machine.VsphereMachineClassSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// TODO: Validation
	allErrs = append(allErrs, validateSecretRef(spec.SecretRef, field.NewPath("spec.secretRef"))...)

	return allErrs
}
