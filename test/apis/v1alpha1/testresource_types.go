/*
Copyright 2025 The KubeFleet Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestResourceSpec defines the desired state of TestResource
type TestResourceSpec struct {
	// +optional
	Foo string `json:"foo,omitempty"`

	// +optional
	Bar string `json:"bar,omitempty"`

	// optional
	NewField string `json:"newField,omitempty"`

	// +optional
	Items []string `json:"items,omitempty"`

	// +optional
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty"`
}

// TestResourceStatus defines the observed state of TestResource
type TestResourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TestResource is the Schema for the testresources API
type TestResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +required
	Spec TestResourceSpec `json:"spec"`
	// +optional
	Status TestResourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TestResourceList contains a list of TestResource
type TestResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TestResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TestResource{}, &TestResourceList{})
}
