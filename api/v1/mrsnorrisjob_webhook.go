/*
Copyright 2023.

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
	"github.com/robfig/cron"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var mrsnorrisjoblog = logf.Log.WithName("mrsnorrisjob-resource")

func (r *MrsNorrisJob) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-filch-caretaker-sh-v1-mrsnorrisjob,mutating=true,failurePolicy=fail,sideEffects=None,groups=filch.caretaker.sh,resources=mrsnorrisjobs,verbs=create;update,versions=v1,name=mmrsnorrisjob.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &MrsNorrisJob{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *MrsNorrisJob) Default() {
	mrsnorrisjoblog.Info("default", "name", r.Name)

	if r.Spec.ConcurrencyPolicy == "" {
		r.Spec.ConcurrencyPolicy = AllowConcurrent
	}
	if r.Spec.Suspend == nil {
		r.Spec.Suspend = new(bool)
	}
	if r.Spec.SuccessfulJobsHistoryLimit == nil {
		r.Spec.SuccessfulJobsHistoryLimit = new(int32)
		*r.Spec.SuccessfulJobsHistoryLimit = 3
	}
	if r.Spec.FailedJobsHistoryLimit == nil {
		r.Spec.FailedJobsHistoryLimit = new(int32)
		*r.Spec.FailedJobsHistoryLimit = 1
	}
}

//+kubebuilder:webhook:verbs=create;update;delete,path=/validate-filch-caretaker-sh-v1-mrsnorrisjob,mutating=false,failurePolicy=fail,groups=filch.caretaker.sh,resources=mrsnorrisjobs,versions=v1,name=vmrsnorrisjob.kb.io,sideEffects=None,admissionReviewVersions=v1

var _ webhook.Validator = &MrsNorrisJob{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MrsNorrisJob) ValidateCreate() error {
	mrsnorrisjoblog.Info("validate create", "name", r.Name)

	return r.validateMrsNorrisJob()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *MrsNorrisJob) ValidateUpdate(old runtime.Object) error {
	mrsnorrisjoblog.Info("validate update", "name", r.Name)

	return r.validateMrsNorrisJob()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *MrsNorrisJob) ValidateDelete() error {
	mrsnorrisjoblog.Info("validate delete", "name", r.Name)

	return nil
}

func (r *MrsNorrisJob) validateMrsNorrisJob() error {
	var allErrs field.ErrorList
	if err := r.validateMrsNorrisJobName(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateMrsNorrisJobSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "filch.caretaker.sh", Kind: "MrsNorrisJob"},
		r.Name, allErrs)
}

func (r *MrsNorrisJob) validateMrsNorrisJobSpec() *field.Error {
	// The field helpers from the kubernetes API machinery help us return nicely
	// structured validation errors.
	return validateScheduleFormat(
		r.Spec.Schedule,
		field.NewPath("spec").Child("schedule"))
}

func validateScheduleFormat(schedule string, fldPath *field.Path) *field.Error {
	if _, err := cron.ParseStandard(schedule); err != nil {
		return field.Invalid(fldPath, schedule, err.Error())
	}
	return nil
}

func (r *MrsNorrisJob) validateMrsNorrisJobName() *field.Error {
	if len(r.ObjectMeta.Name) > validationutils.DNS1035LabelMaxLength-11 {
		// The job name length is 63 character like all Kubernetes objects
		// (which must fit in a DNS subdomain). The mrsnorrisjob controller appends
		// a 11-character suffix to the mrsnorrisjob (`-$TIMESTAMP`) when creating
		// a job. The job name length limit is 63 characters. Therefore mrsnorrisjob
		// names must have length <= 63-11=52. If we don't validate this here,
		// then job creation will fail later.
		return field.Invalid(field.NewPath("metadata").Child("name"), r.Name, "must be no more than 52 characters")
	}
	return nil
}
