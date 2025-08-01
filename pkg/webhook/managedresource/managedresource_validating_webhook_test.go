/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package managedresource

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	admissionv1 "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func Test_managedResourceValidator_Handle(t *testing.T) {
	const fleet1p = "fleet1p"
	validator := &managedResourceValidator{
		whiteListedUsers: []string{fleet1p},
	}

	tests := []struct {
		name           string
		username       string
		operation      admissionv1.Operation
		oldLabels      map[string]string
		oldAnnotations map[string]string
		newLabels      map[string]string
		newAnnotations map[string]string
		expectedResp   admission.Response
		modReq         func(*admission.Request)
	}{
		{
			name:           "allowed when not managed by arm",
			operation:      admissionv1.Update,
			oldLabels:      map[string]string{"foo": "bar"},
			oldAnnotations: map[string]string{"baz": "qux"},
			newLabels:      map[string]string{"foo": "bar"},
			newAnnotations: map[string]string{"baz": "qux"},
			expectedResp:   admission.Allowed(""),
		},
		{
			name:         "denied - error on getLabelsAndAnnotations failure",
			operation:    admissionv1.Create,
			expectedResp: admission.Errored(http.StatusInternalServerError, fmt.Errorf("error decoding string from json: unexpected trailing data at offset 9")),
			modReq: func(req *admission.Request) {
				req.Object = runtime.RawExtension{Raw: []byte(`"invalid"}`)} // Invalid object without labels or annotations
			},
		},
		{
			name:           "denied - managed by arm in labels, not whitelisted",
			operation:      admissionv1.Create,
			oldLabels:      nil,
			oldAnnotations: nil,
			newLabels:      map[string]string{ManagedByArmKey: ManagedByArmValue},
			newAnnotations: nil,
			expectedResp:   admission.Denied(fmt.Sprintf(resourceDeniedFormat, metav1.GroupVersionKind{Kind: "TestKind"}, "test-resource", "default")),
		},
		{
			name:           "denied - managed by arm in annotations, not whitelisted",
			operation:      admissionv1.Update,
			oldLabels:      nil,
			oldAnnotations: nil,
			newLabels:      nil,
			newAnnotations: map[string]string{ManagedByArmKey: ManagedByArmValue},
			expectedResp:   admission.Denied(fmt.Sprintf(resourceDeniedFormat, metav1.GroupVersionKind{Kind: "TestKind"}, "test-resource", "default")),
		},
		{
			name:           "allowed for other operations",
			operation:      admissionv1.Connect,
			oldLabels:      nil,
			oldAnnotations: nil,
			newLabels:      nil,
			newAnnotations: nil,
			expectedResp:   admission.Allowed(""),
		},
		{
			name:           "allowed for other operations - managed by arm, but user whitelisted",
			username:       "fleet1p",
			operation:      admissionv1.Update,
			oldLabels:      map[string]string{"managedBy": ManagedByArmValue},
			oldAnnotations: nil,
			newLabels:      nil,
			newAnnotations: nil,
			expectedResp:   admission.Allowed(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldObj := makeUnstructured(tt.oldLabels, tt.oldAnnotations)
			newObj := makeUnstructured(tt.newLabels, tt.newAnnotations)
			req := admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Operation:   tt.operation,
					Name:        "test-resource",
					Namespace:   "default",
					OldObject:   runtime.RawExtension{Object: oldObj},
					Object:      runtime.RawExtension{Object: newObj},
					Kind:        metav1.GroupVersionKind{Kind: "TestKind"},
					RequestKind: &metav1.GroupVersionKind{Kind: "TestKind"},
				},
			}
			req.UserInfo = authenticationv1.UserInfo{
				Username: tt.username,
			}
			if tt.modReq != nil {
				tt.modReq(&req)
			}
			resp := validator.Handle(context.Background(), req)
			if diff := cmp.Diff(tt.expectedResp.Result, resp.Result); diff != "" {
				t.Errorf("managedResourceValidator Handle response (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_getLabelsAndAnnotations(t *testing.T) {
	tests := []struct {
		name            string
		obj             runtime.RawExtension
		wantLabels      map[string]string
		wantAnnotations map[string]string
		expectError     bool
	}{
		{
			name: "object with labels and annotations",
			obj: runtime.RawExtension{
				Object: &metav1.PartialObjectMetadata{
					ObjectMeta: metav1.ObjectMeta{
						Labels:      map[string]string{"foo": "bar"},
						Annotations: map[string]string{"baz": "qux"},
					},
				},
			},
			wantLabels:      map[string]string{"foo": "bar"},
			wantAnnotations: map[string]string{"baz": "qux"},
			expectError:     false,
		},
		{
			name: "object with no labels or annotations",
			obj: runtime.RawExtension{
				Object: &metav1.PartialObjectMetadata{
					ObjectMeta: metav1.ObjectMeta{},
				},
			},
			wantLabels:      nil,
			wantAnnotations: nil,
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels, annotations, err := getLabelsAndAnnotations(tt.obj)
			if err != nil && !tt.expectError {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.wantLabels, labels); diff != "" {
				t.Errorf("labels mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantAnnotations, annotations); diff != "" {
				t.Errorf("annotations mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_managedByArm(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]string
		want bool
	}{
		{
			name: "nil map",
			m:    nil,
			want: false,
		},
		{
			name: "empty map",
			m:    map[string]string{},
			want: false,
		},
		{
			name: "key missing",
			m:    map[string]string{"foo": "bar"},
			want: false,
		},
		{
			name: "key present, not managed key",
			m:    map[string]string{"managingBy": ManagedByArmValue},
			want: false,
		},
		{
			name: "key present, not managed value",
			m:    map[string]string{ManagedByArmKey: "not-arm"},
			want: false,
		},
		{
			name: "key present, managed key and value",
			m:    map[string]string{ManagedByArmKey: ManagedByArmValue},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, managedByArm(tt.m)); diff != "" {
				t.Errorf("managedByArm result (-want +got):\n%s", diff)
			}
		})
	}
}

func makeUnstructured(labels, annotations map[string]string) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetLabels(labels)
	obj.SetAnnotations(annotations)
	return obj
}
