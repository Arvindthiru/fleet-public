/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package workgenerator

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	fleetv1beta1 "go.goms.io/fleet/apis/placement/v1beta1"
	"go.goms.io/fleet/pkg/controllers/work"
	"go.goms.io/fleet/pkg/utils/condition"
	"go.goms.io/fleet/pkg/utils/controller"
)

func TestGetWorkNamePrefixFromSnapshotName(t *testing.T) {
	tests := map[string]struct {
		resourceSnapshot *fleetv1beta1.ClusterResourceSnapshot
		wantErr          error
		wantedName       string
	}{
		"the work name is crp name + \"work\", if there is only one resource snapshot": {
			resourceSnapshot: &fleetv1beta1.ClusterResourceSnapshot{
				ObjectMeta: metav1.ObjectMeta{
					Name: "placement-2",
					Labels: map[string]string{
						fleetv1beta1.CRPTrackingLabel: "placement",
					},
				},
			},
			wantErr:    nil,
			wantedName: "placement-work",
		},
		"should return error if the resource snapshot has negative subindex": {
			resourceSnapshot: &fleetv1beta1.ClusterResourceSnapshot{
				ObjectMeta: metav1.ObjectMeta{
					Name: "placement-1-2",
					Labels: map[string]string{
						fleetv1beta1.CRPTrackingLabel: "placement",
					},
					Annotations: map[string]string{
						fleetv1beta1.SubindexOfResourceSnapshotAnnotation: "-1",
					},
				},
			},
			wantErr:    controller.ErrUnexpectedBehavior,
			wantedName: "",
		},
		"the work name is the concatenation of the crp name and subindex start at 0": {
			resourceSnapshot: &fleetv1beta1.ClusterResourceSnapshot{
				ObjectMeta: metav1.ObjectMeta{
					Name: "placement-1-2",
					Labels: map[string]string{
						fleetv1beta1.CRPTrackingLabel: "placement",
					},
					Annotations: map[string]string{
						fleetv1beta1.SubindexOfResourceSnapshotAnnotation: "0",
					},
				},
			},
			wantErr:    nil,
			wantedName: "placement-0",
		},
		"the work name is the concatenation of the crp name and subindex": {
			resourceSnapshot: &fleetv1beta1.ClusterResourceSnapshot{
				ObjectMeta: metav1.ObjectMeta{
					Name: "placement-1-2",
					Labels: map[string]string{
						fleetv1beta1.CRPTrackingLabel: "placement",
					},
					Annotations: map[string]string{
						fleetv1beta1.SubindexOfResourceSnapshotAnnotation: "2",
					},
				},
			},
			wantErr:    nil,
			wantedName: "placement-2",
		},
		"test return error if the resource snapshot has invalid subindex": {
			resourceSnapshot: &fleetv1beta1.ClusterResourceSnapshot{
				ObjectMeta: metav1.ObjectMeta{
					Name: "placement-1-2",
					Labels: map[string]string{
						fleetv1beta1.CRPTrackingLabel: "placement",
					},
					Annotations: map[string]string{
						fleetv1beta1.SubindexOfResourceSnapshotAnnotation: "what?",
					},
				},
			},
			wantErr:    controller.ErrUnexpectedBehavior,
			wantedName: "",
		},
		"test return error if the resource snapshot does not have CRP track": {
			resourceSnapshot: &fleetv1beta1.ClusterResourceSnapshot{
				ObjectMeta: metav1.ObjectMeta{
					Name: "placement-1-2",
					Annotations: map[string]string{
						fleetv1beta1.SubindexOfResourceSnapshotAnnotation: "what?",
					},
				},
			},
			wantErr:    controller.ErrUnexpectedBehavior,
			wantedName: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			workName, err := getWorkNamePrefixFromSnapshotName(tt.resourceSnapshot)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("failed getWorkNamePrefixFromSnapshotName test `%s` error = %v, wantErr %v", name, err, tt.wantErr)
				return
			}
			if workName != tt.wantedName {
				t.Errorf("getWorkNamePrefixFromSnapshotName test `%s` workName = `%v`, wantedName `%v`", name, workName, tt.wantedName)
			}
		})
	}
}

func TestBuildAllWorkAppliedCondition(t *testing.T) {
	tests := map[string]struct {
		works      map[string]*fleetv1beta1.Work
		generation int64
		want       metav1.Condition
	}{
		"applied should be true if all work applied": {
			works: map[string]*fleetv1beta1.Work{
				"appliedWork1": {
					ObjectMeta: metav1.ObjectMeta{
						Name:       "work1",
						Generation: 123,
					},
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:               fleetv1beta1.WorkConditionTypeApplied,
								Status:             metav1.ConditionTrue,
								ObservedGeneration: 123,
							},
						},
					},
				},
				"appliedWork2": {
					ObjectMeta: metav1.ObjectMeta{
						Name:       "work2",
						Generation: 12,
					},
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:               fleetv1beta1.WorkConditionTypeApplied,
								Status:             metav1.ConditionTrue,
								ObservedGeneration: 12,
							},
						},
					},
				},
			},
			generation: 1,
			want: metav1.Condition{
				Status:             metav1.ConditionTrue,
				Type:               string(fleetv1beta1.ResourceBindingApplied),
				Reason:             condition.AllWorkAppliedReason,
				ObservedGeneration: 1,
			},
		},
		"applied should be false if not all work applied to the latest generation": {
			works: map[string]*fleetv1beta1.Work{
				"notAppliedWork1": {
					ObjectMeta: metav1.ObjectMeta{
						Name:       "work1",
						Generation: 123,
					},
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:               fleetv1beta1.WorkConditionTypeApplied,
								Status:             metav1.ConditionTrue,
								ObservedGeneration: 122, // not the latest generation
							},
						},
					},
				},
				"appliedWork2": {
					ObjectMeta: metav1.ObjectMeta{
						Name:       "work2",
						Generation: 12,
					},
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:               fleetv1beta1.WorkConditionTypeApplied,
								Status:             metav1.ConditionTrue,
								ObservedGeneration: 12,
							},
						},
					},
				},
			},
			generation: 1,
			want: metav1.Condition{
				Status:             metav1.ConditionFalse,
				Type:               string(fleetv1beta1.ResourceBindingApplied),
				Reason:             condition.WorkNotAppliedReason,
				ObservedGeneration: 1,
			},
		},
		"applied should be false if not all work has applied": {
			works: map[string]*fleetv1beta1.Work{
				"appliedWork1": {
					ObjectMeta: metav1.ObjectMeta{
						Name:       "work1",
						Generation: 123,
					},
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:               fleetv1beta1.WorkConditionTypeApplied,
								Status:             metav1.ConditionTrue,
								ObservedGeneration: 122, // not the latest generation
							},
						},
					},
				},
				"notAppliedWork2": {
					ObjectMeta: metav1.ObjectMeta{
						Name:       "work2",
						Generation: 12,
					},
				},
			},
			generation: 1,
			want: metav1.Condition{
				Status:             metav1.ConditionFalse,
				Type:               string(fleetv1beta1.ResourceBindingApplied),
				Reason:             condition.WorkNotAppliedReason,
				ObservedGeneration: 1,
			},
		},
		"applied should be false if some work applied condition is unknown": {
			works: map[string]*fleetv1beta1.Work{
				"appliedWork1": {
					ObjectMeta: metav1.ObjectMeta{
						Name:       "work1",
						Generation: 123,
					},
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:               fleetv1beta1.WorkConditionTypeApplied,
								Status:             metav1.ConditionUnknown,
								ObservedGeneration: 123,
							},
						},
					},
				},
				"appliedWork2": {
					ObjectMeta: metav1.ObjectMeta{
						Name:       "work2",
						Generation: 12,
					},
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:               fleetv1beta1.WorkConditionTypeApplied,
								Status:             metav1.ConditionTrue,
								ObservedGeneration: 12,
							},
						},
					},
				},
			},
			generation: 1,
			want: metav1.Condition{
				Status:             metav1.ConditionFalse,
				Type:               string(fleetv1beta1.ResourceBindingApplied),
				Reason:             condition.WorkNotAppliedReason,
				ObservedGeneration: 1,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			binding := &fleetv1beta1.ClusterResourceBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "test",
					Generation: tt.generation,
				},
			}
			got := buildAllWorkAppliedCondition(tt.works, binding)
			if diff := cmp.Diff(got, tt.want, cmpConditionOption); diff != "" {
				t.Errorf("buildAllWorkAppliedCondition test `%s` mismatch (-got +want):\n%s", name, diff)
			}
		})
	}
}

func TestBuildAllWorkAvailableCondition(t *testing.T) {
	tests := map[string]struct {
		works   map[string]*fleetv1beta1.Work
		binding *fleetv1beta1.ClusterResourceBinding
		want    metav1.Condition
	}{
		"All works are available": {
			works: map[string]*fleetv1beta1.Work{
				"work1": {
					ObjectMeta: metav1.ObjectMeta{
						Name: "work1",
					},
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Reason: "any",
								Status: metav1.ConditionTrue,
							},
						},
					},
				},
				"work2": {
					ObjectMeta: metav1.ObjectMeta{
						Name: "work2",
					},
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Reason: "any",
								Status: metav1.ConditionTrue,
							},
						},
					},
				},
			},
			binding: &fleetv1beta1.ClusterResourceBinding{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 1,
				},
			},
			want: metav1.Condition{
				Status:             metav1.ConditionTrue,
				Type:               string(fleetv1beta1.ResourceBindingAvailable),
				Reason:             condition.AllWorkAvailableReason,
				ObservedGeneration: 1,
			},
		},
		"All works are available but one of them is not trackable": {
			works: map[string]*fleetv1beta1.Work{
				"work1": {
					ObjectMeta: metav1.ObjectMeta{
						Name: "work1",
					},
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Reason: work.WorkNotTrackableReason,
								Status: metav1.ConditionTrue,
							},
						},
					},
				},
				"work2": {
					ObjectMeta: metav1.ObjectMeta{
						Name: "work2",
					},
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Reason: "any",
								Status: metav1.ConditionTrue,
							},
						},
					},
				},
			},
			binding: &fleetv1beta1.ClusterResourceBinding{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 1,
				},
			},
			want: metav1.Condition{
				Status:             metav1.ConditionTrue,
				Type:               string(fleetv1beta1.ResourceBindingAvailable),
				Reason:             work.WorkNotTrackableReason,
				ObservedGeneration: 1,
			},
		},
		"Not all works are available": {
			works: map[string]*fleetv1beta1.Work{
				"work1": {
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Status: metav1.ConditionTrue,
							},
						},
					},
				},
				"work2": {
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Status: metav1.ConditionFalse,
							},
						},
					},
				},
			},
			binding: &fleetv1beta1.ClusterResourceBinding{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 1,
				},
			},
			want: metav1.Condition{
				Status:             metav1.ConditionFalse,
				Type:               string(fleetv1beta1.ResourceBindingAvailable),
				Reason:             condition.WorkNotAvailableReason,
				Message:            "work object work2 is not available",
				ObservedGeneration: 1,
			},
		},
		"Available condition of one work is unknown": {
			works: map[string]*fleetv1beta1.Work{
				"work1": {
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Status: metav1.ConditionTrue,
							},
						},
					},
				},
				"work2": {
					Status: fleetv1beta1.WorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Status: metav1.ConditionUnknown,
							},
						},
					},
				},
			},
			binding: &fleetv1beta1.ClusterResourceBinding{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 1,
				},
			},
			want: metav1.Condition{
				Status:             metav1.ConditionFalse,
				Type:               string(fleetv1beta1.ResourceBindingAvailable),
				Reason:             condition.WorkNotAvailableReason,
				Message:            "work object work2 is not available",
				ObservedGeneration: 1,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := buildAllWorkAvailableCondition(tt.works, tt.binding)
			if diff := cmp.Diff(got, tt.want, cmpConditionOption); diff != "" {
				t.Errorf("buildAllWorkAvailableCondition test `%s` mismatch (-got +want):\n%s", name, diff)
			}
		})
	}
}

func TestSetBindingStatus(t *testing.T) {
	tests := map[string]struct {
		works                           map[string]*fleetv1beta1.Work
		maxFailedResourcePlacementLimit *int
		want                            []fleetv1beta1.FailedResourcePlacement
	}{
		"NoWorks": {
			works: map[string]*fleetv1beta1.Work{},
			want:  nil,
		},
		"both work are available": {
			works: map[string]*fleetv1beta1.Work{
				"work1": {
					Status: fleetv1beta1.WorkStatus{
						ManifestConditions: []fleetv1beta1.ManifestCondition{
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   1,
									Group:     "",
									Version:   "v1",
									Kind:      "Service",
									Name:      "svc-name",
									Namespace: "svc-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionTrue,
									},
									{
										Type:   fleetv1beta1.WorkConditionTypeAvailable,
										Status: metav1.ConditionTrue,
									},
								},
							},
						},
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeApplied,
								Status: metav1.ConditionTrue,
							},
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Status: metav1.ConditionTrue,
							},
						},
					},
				},
				"work2": {
					Status: fleetv1beta1.WorkStatus{
						ManifestConditions: []fleetv1beta1.ManifestCondition{
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   0,
									Group:     "",
									Version:   "v1",
									Kind:      "ConfigMap",
									Name:      "config-name",
									Namespace: "config-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionTrue,
									},
									{
										Type:   fleetv1beta1.WorkConditionTypeAvailable,
										Status: metav1.ConditionTrue,
									},
								},
							},
						},
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeApplied,
								Status: metav1.ConditionTrue,
							},
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Status: metav1.ConditionTrue,
							},
						},
					},
				},
			},
			want: nil,
		},
		"One work has one not available and one work has one not applied": {
			works: map[string]*fleetv1beta1.Work{
				"work1": {
					Status: fleetv1beta1.WorkStatus{
						ManifestConditions: []fleetv1beta1.ManifestCondition{
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   0,
									Group:     "",
									Version:   "v1",
									Kind:      "ConfigMap",
									Name:      "config-name",
									Namespace: "config-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionTrue,
									},
									{
										Type:   fleetv1beta1.WorkConditionTypeAvailable,
										Status: metav1.ConditionFalse,
									},
								},
							},
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   1,
									Group:     "",
									Version:   "v1",
									Kind:      "Service",
									Name:      "svc-name",
									Namespace: "svc-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionTrue,
									},
									{
										Type:   fleetv1beta1.WorkConditionTypeAvailable,
										Status: metav1.ConditionTrue,
									},
								},
							},
						},
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeApplied,
								Status: metav1.ConditionTrue,
							},
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Status: metav1.ConditionFalse,
							},
						},
					},
				},
				"work2": {
					Status: fleetv1beta1.WorkStatus{
						ManifestConditions: []fleetv1beta1.ManifestCondition{
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   0,
									Group:     "",
									Version:   "v1",
									Kind:      "ConfigMap",
									Name:      "config-name-1",
									Namespace: "config-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionTrue,
									},
									{
										Type:   fleetv1beta1.WorkConditionTypeAvailable,
										Status: metav1.ConditionTrue,
									},
								},
							},
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   1,
									Group:     "",
									Version:   "v1",
									Kind:      "Service",
									Name:      "svc-name-1",
									Namespace: "svc-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionFalse,
									},
								},
							},
						},
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeApplied,
								Status: metav1.ConditionFalse,
							},
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Status: metav1.ConditionFalse,
							},
						},
					},
				},
			},
			want: []fleetv1beta1.FailedResourcePlacement{
				{
					ResourceIdentifier: fleetv1beta1.ResourceIdentifier{
						Group:     "",
						Version:   "v1",
						Kind:      "ConfigMap",
						Name:      "config-name",
						Namespace: "config-namespace",
					},
					Condition: metav1.Condition{
						Type:   fleetv1beta1.WorkConditionTypeAvailable,
						Status: metav1.ConditionFalse,
					},
				},
				{
					ResourceIdentifier: fleetv1beta1.ResourceIdentifier{
						Group:     "",
						Version:   "v1",
						Kind:      "Service",
						Name:      "svc-name-1",
						Namespace: "svc-namespace",
					},
					Condition: metav1.Condition{
						Type:   fleetv1beta1.WorkConditionTypeApplied,
						Status: metav1.ConditionFalse,
					},
				},
			},
		},
		"One work has one not available and one work has one not applied (exceed the maxFailedResourcePlacementLimit)": {
			works: map[string]*fleetv1beta1.Work{
				"work1": {
					Status: fleetv1beta1.WorkStatus{
						ManifestConditions: []fleetv1beta1.ManifestCondition{
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   0,
									Group:     "",
									Version:   "v1",
									Kind:      "ConfigMap",
									Name:      "config-name",
									Namespace: "config-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionTrue,
									},
									{
										Type:   fleetv1beta1.WorkConditionTypeAvailable,
										Status: metav1.ConditionFalse,
									},
								},
							},
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   1,
									Group:     "",
									Version:   "v1",
									Kind:      "Service",
									Name:      "svc-name",
									Namespace: "svc-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionTrue,
									},
									{
										Type:   fleetv1beta1.WorkConditionTypeAvailable,
										Status: metav1.ConditionTrue,
									},
								},
							},
						},
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeApplied,
								Status: metav1.ConditionTrue,
							},
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Status: metav1.ConditionFalse,
							},
						},
					},
				},
				"work2": {
					Status: fleetv1beta1.WorkStatus{
						ManifestConditions: []fleetv1beta1.ManifestCondition{
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   0,
									Group:     "",
									Version:   "v1",
									Kind:      "ConfigMap",
									Name:      "config-name-1",
									Namespace: "config-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionTrue,
									},
									{
										Type:   fleetv1beta1.WorkConditionTypeAvailable,
										Status: metav1.ConditionTrue,
									},
								},
							},
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   1,
									Group:     "",
									Version:   "v1",
									Kind:      "Service",
									Name:      "svc-name-1",
									Namespace: "svc-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionFalse,
									},
								},
							},
						},
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeApplied,
								Status: metav1.ConditionFalse,
							},
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Status: metav1.ConditionFalse,
							},
						},
					},
				},
			},
			maxFailedResourcePlacementLimit: ptr.To(1),
			want: []fleetv1beta1.FailedResourcePlacement{
				{
					ResourceIdentifier: fleetv1beta1.ResourceIdentifier{
						Group:     "",
						Version:   "v1",
						Kind:      "ConfigMap",
						Name:      "config-name",
						Namespace: "config-namespace",
					},
					Condition: metav1.Condition{
						Type:   fleetv1beta1.WorkConditionTypeAvailable,
						Status: metav1.ConditionFalse,
					},
				},
			},
		},
		"One work has one not available and one work all available": {
			works: map[string]*fleetv1beta1.Work{
				"work1": {
					Status: fleetv1beta1.WorkStatus{
						ManifestConditions: []fleetv1beta1.ManifestCondition{
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   0,
									Group:     "",
									Version:   "v1",
									Kind:      "ConfigMap",
									Name:      "config-name",
									Namespace: "config-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionTrue,
									},
									{
										Type:   fleetv1beta1.WorkConditionTypeAvailable,
										Status: metav1.ConditionTrue,
									},
								},
							},
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   1,
									Group:     "",
									Version:   "v1",
									Kind:      "Service",
									Name:      "svc-name",
									Namespace: "svc-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionTrue,
									},
									{
										Type:   fleetv1beta1.WorkConditionTypeAvailable,
										Status: metav1.ConditionTrue,
									},
								},
							},
						},
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeApplied,
								Status: metav1.ConditionTrue,
							},
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Status: metav1.ConditionTrue,
							},
						},
					},
				},
				"all available work": {
					Status: fleetv1beta1.WorkStatus{
						ManifestConditions: []fleetv1beta1.ManifestCondition{
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   0,
									Group:     "",
									Version:   "v1",
									Kind:      "ConfigMap",
									Name:      "config-name",
									Namespace: "config-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionTrue,
									},
									{
										Type:   fleetv1beta1.WorkConditionTypeAvailable,
										Status: metav1.ConditionTrue,
									},
								},
							},
							{
								Identifier: fleetv1beta1.WorkResourceIdentifier{
									Ordinal:   1,
									Group:     "",
									Version:   "v1",
									Kind:      "Service",
									Name:      "svc-name",
									Namespace: "svc-namespace",
								},
								Conditions: []metav1.Condition{
									{
										Type:   fleetv1beta1.WorkConditionTypeApplied,
										Status: metav1.ConditionFalse,
									},
								},
							},
						},
						Conditions: []metav1.Condition{
							{
								Type:   fleetv1beta1.WorkConditionTypeApplied,
								Status: metav1.ConditionFalse,
							},
							{
								Type:   fleetv1beta1.WorkConditionTypeAvailable,
								Status: metav1.ConditionFalse,
							},
						},
					},
				},
			},
			want: []fleetv1beta1.FailedResourcePlacement{
				{
					ResourceIdentifier: fleetv1beta1.ResourceIdentifier{
						Group:     "",
						Version:   "v1",
						Kind:      "Service",
						Name:      "svc-name",
						Namespace: "svc-namespace",
					},
					Condition: metav1.Condition{
						Type:   fleetv1beta1.WorkConditionTypeApplied,
						Status: metav1.ConditionFalse,
					},
				},
			},
		},
	}

	originalMaxFailedResourcePlacementLimit := maxFailedResourcePlacementLimit
	defer func() {
		maxFailedResourcePlacementLimit = originalMaxFailedResourcePlacementLimit
	}()
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.maxFailedResourcePlacementLimit != nil {
				maxFailedResourcePlacementLimit = *tt.maxFailedResourcePlacementLimit
			} else {
				maxFailedResourcePlacementLimit = originalMaxFailedResourcePlacementLimit
			}

			binding := &fleetv1beta1.ClusterResourceBinding{}
			setBindingStatus(tt.works, binding)
			got := binding.Status.FailedPlacements
			// setBindingStatus is using map to populate the failedResourcePlacement.
			// There is no default order in traversing the map.
			// When the result of failedResourcePlacement exceeds the limit, the result will be truncated and cannot be
			// guaranteed.
			if maxFailedResourcePlacementLimit == len(tt.want) {
				opt := cmp.Comparer(func(x, y fleetv1beta1.FailedResourcePlacement) bool {
					return x.Condition.Status == y.Condition.Status // condition should be set as false
				})
				if diff := cmp.Diff(got, tt.want, opt); diff != "" {
					t.Errorf("setBindingStatus got FailedPlacements mismatch (-got +want):\n%s", diff)
				}
				return
			}

			statusCmpOptions := []cmp.Option{
				cmpopts.SortSlices(func(i, j fleetv1beta1.FailedResourcePlacement) bool {
					if i.Group < j.Group {
						return true
					}
					if i.Kind < j.Kind {
						return true
					}
					return i.Name < j.Name
				}),
			}
			if diff := cmp.Diff(got, tt.want, statusCmpOptions...); diff != "" {
				t.Errorf("setBindingStatus got FailedPlacements mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestExtractFailedResourcePlacementsFromWork(t *testing.T) {
	var statusCmpOptions = []cmp.Option{
		// ignore the message as we may change the message in the future
		cmpopts.IgnoreFields(metav1.Condition{}, "Message"),
		cmpopts.SortSlices(func(c1, c2 metav1.Condition) bool {
			return c1.Type < c2.Type
		}),
		cmpopts.SortSlices(func(s1, s2 string) bool {
			return s1 < s2
		}),
		cmpopts.SortSlices(func(n1, n2 fleetv1beta1.NamespacedName) bool {
			if n1.Namespace == n2.Namespace {
				return n1.Name < n2.Name
			}
			return n1.Namespace < n2.Namespace
		}),
		cmpopts.SortSlices(func(f1, f2 fleetv1beta1.FailedResourcePlacement) bool {
			return f1.ResourceIdentifier.Kind < f2.ResourceIdentifier.Kind
		}),
		cmp.Comparer(func(t1, t2 metav1.Time) bool {
			if t1.Time.IsZero() || t2.Time.IsZero() {
				return true // treat them as equal
			}
			if t1.Time.After(t2.Time) {
				t1, t2 = t2, t1 // ensure t1 is always before t2
			}
			// we're within the margin (10s) if x + margin >= y
			return !t1.Time.Add(10 * time.Second).Before(t2.Time)
		}),
	}
	workGeneration := int64(12)
	tests := []struct {
		name string
		work fleetv1beta1.Work
		want []fleetv1beta1.FailedResourcePlacement
	}{
		{
			name: "apply is true and available is false",
			work: fleetv1beta1.Work{
				ObjectMeta: metav1.ObjectMeta{
					Generation: workGeneration,
				},
				Status: fleetv1beta1.WorkStatus{
					ManifestConditions: []fleetv1beta1.ManifestCondition{
						{
							Identifier: fleetv1beta1.WorkResourceIdentifier{
								Ordinal:   0,
								Group:     "",
								Version:   "v1",
								Kind:      "ConfigMap",
								Name:      "config-name",
								Namespace: "config-namespace",
							},
							Conditions: []metav1.Condition{
								{
									Type:   fleetv1beta1.WorkConditionTypeApplied,
									Status: metav1.ConditionTrue,
								},
								{
									Type:   fleetv1beta1.WorkConditionTypeAvailable,
									Status: metav1.ConditionFalse,
								},
							},
						},
						{
							Identifier: fleetv1beta1.WorkResourceIdentifier{
								Ordinal:   1,
								Group:     "",
								Version:   "v1",
								Kind:      "Service",
								Name:      "svc-name",
								Namespace: "svc-namespace",
							},
							Conditions: []metav1.Condition{
								{
									Type:   fleetv1beta1.WorkConditionTypeApplied,
									Status: metav1.ConditionTrue,
								},
								{
									Type:   fleetv1beta1.WorkConditionTypeAvailable,
									Status: metav1.ConditionTrue,
								},
							},
						},
					},
					Conditions: []metav1.Condition{
						{
							Type:               fleetv1beta1.WorkConditionTypeApplied,
							Status:             metav1.ConditionTrue,
							ObservedGeneration: workGeneration,
						},
						{
							Type:               fleetv1beta1.WorkConditionTypeAvailable,
							Status:             metav1.ConditionFalse,
							ObservedGeneration: workGeneration,
						},
					},
				},
			},
			want: []fleetv1beta1.FailedResourcePlacement{
				{
					ResourceIdentifier: fleetv1beta1.ResourceIdentifier{
						Group:     "",
						Version:   "v1",
						Kind:      "ConfigMap",
						Name:      "config-name",
						Namespace: "config-namespace",
					},
					Condition: metav1.Condition{
						Type:   fleetv1beta1.WorkConditionTypeAvailable,
						Status: metav1.ConditionFalse,
					},
				},
			},
		},
		{
			name: "apply is true and available is false for enveloped object",
			work: fleetv1beta1.Work{
				ObjectMeta: metav1.ObjectMeta{
					Generation: workGeneration,
					Labels: map[string]string{
						fleetv1beta1.EnvelopeNameLabel:      "test-env",
						fleetv1beta1.EnvelopeNamespaceLabel: "test-env-ns",
						fleetv1beta1.EnvelopeTypeLabel:      "pod",
					},
				},
				Status: fleetv1beta1.WorkStatus{
					ManifestConditions: []fleetv1beta1.ManifestCondition{
						{
							Identifier: fleetv1beta1.WorkResourceIdentifier{
								Ordinal:   0,
								Group:     "",
								Version:   "v1",
								Kind:      "ConfigMap",
								Name:      "config-name",
								Namespace: "config-namespace",
							},
							Conditions: []metav1.Condition{
								{
									Type:   fleetv1beta1.WorkConditionTypeApplied,
									Status: metav1.ConditionTrue,
								},
								{
									Type:   fleetv1beta1.WorkConditionTypeAvailable,
									Status: metav1.ConditionFalse,
								},
							},
						},
					},
					Conditions: []metav1.Condition{
						{
							Type:               fleetv1beta1.WorkConditionTypeApplied,
							Status:             metav1.ConditionTrue,
							ObservedGeneration: workGeneration,
						},
						{
							Type:               fleetv1beta1.WorkConditionTypeAvailable,
							Status:             metav1.ConditionFalse,
							ObservedGeneration: workGeneration,
						},
					},
				},
			},
			want: []fleetv1beta1.FailedResourcePlacement{
				{
					ResourceIdentifier: fleetv1beta1.ResourceIdentifier{
						Group:     "",
						Version:   "v1",
						Kind:      "ConfigMap",
						Name:      "config-name",
						Namespace: "config-namespace",
						Envelope: &fleetv1beta1.EnvelopeIdentifier{
							Name:      "test-env",
							Namespace: "test-env-ns",
							Type:      "pod",
						},
					},
					Condition: metav1.Condition{
						Type:   fleetv1beta1.WorkConditionTypeAvailable,
						Status: metav1.ConditionFalse,
					},
				},
			},
		},
		{
			name: "both conditions are true",
			work: fleetv1beta1.Work{
				ObjectMeta: metav1.ObjectMeta{
					Generation: workGeneration,
				},
				Status: fleetv1beta1.WorkStatus{
					Conditions: []metav1.Condition{
						{
							Type:               fleetv1beta1.WorkConditionTypeApplied,
							Status:             metav1.ConditionTrue,
							ObservedGeneration: workGeneration,
						},
						{
							Type:               fleetv1beta1.WorkConditionTypeAvailable,
							Status:             metav1.ConditionTrue,
							ObservedGeneration: workGeneration,
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "apply is true and available is unknown",
			work: fleetv1beta1.Work{
				ObjectMeta: metav1.ObjectMeta{
					Generation: workGeneration,
				},
				Status: fleetv1beta1.WorkStatus{
					Conditions: []metav1.Condition{
						{
							Type:               fleetv1beta1.WorkConditionTypeApplied,
							Status:             metav1.ConditionTrue,
							ObservedGeneration: workGeneration,
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "applied is false but not for the latest work",
			work: fleetv1beta1.Work{
				ObjectMeta: metav1.ObjectMeta{
					Generation: workGeneration,
				},
				Status: fleetv1beta1.WorkStatus{
					Conditions: []metav1.Condition{
						{
							Type:               fleetv1beta1.WorkConditionTypeApplied,
							Status:             metav1.ConditionFalse,
							ObservedGeneration: workGeneration - 1,
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "apply is false",
			work: fleetv1beta1.Work{
				ObjectMeta: metav1.ObjectMeta{
					Generation: workGeneration,
				},
				Status: fleetv1beta1.WorkStatus{
					ManifestConditions: []fleetv1beta1.ManifestCondition{
						{
							Identifier: fleetv1beta1.WorkResourceIdentifier{
								Ordinal:   0,
								Group:     "",
								Version:   "v1",
								Kind:      "ConfigMap",
								Name:      "config-name",
								Namespace: "config-namespace",
							},
							Conditions: []metav1.Condition{
								{
									Type:   fleetv1beta1.WorkConditionTypeApplied,
									Status: metav1.ConditionFalse,
								},
							},
						},
					},
					Conditions: []metav1.Condition{
						{
							Type:               fleetv1beta1.WorkConditionTypeApplied,
							Status:             metav1.ConditionFalse,
							ObservedGeneration: workGeneration,
						},
					},
				},
			},
			want: []fleetv1beta1.FailedResourcePlacement{
				{
					ResourceIdentifier: fleetv1beta1.ResourceIdentifier{
						Group:     "",
						Version:   "v1",
						Kind:      "ConfigMap",
						Name:      "config-name",
						Namespace: "config-namespace",
					},
					Condition: metav1.Condition{
						Type:   fleetv1beta1.WorkConditionTypeApplied,
						Status: metav1.ConditionFalse,
					},
				},
			},
		},
		{
			name: "apply is false for enveloped object",
			work: fleetv1beta1.Work{
				ObjectMeta: metav1.ObjectMeta{
					Generation: workGeneration,
					Labels: map[string]string{
						fleetv1beta1.EnvelopeNameLabel:      "test-env",
						fleetv1beta1.EnvelopeNamespaceLabel: "test-env-ns",
						fleetv1beta1.EnvelopeTypeLabel:      "pod",
					},
				},
				Status: fleetv1beta1.WorkStatus{
					ManifestConditions: []fleetv1beta1.ManifestCondition{
						{
							Identifier: fleetv1beta1.WorkResourceIdentifier{
								Ordinal:   0,
								Group:     "",
								Version:   "v1",
								Kind:      "ConfigMap",
								Name:      "config-name",
								Namespace: "config-namespace",
							},
							Conditions: []metav1.Condition{
								{
									Type:   fleetv1beta1.WorkConditionTypeApplied,
									Status: metav1.ConditionFalse,
								},
							},
						},
					},
					Conditions: []metav1.Condition{
						{
							Type:               fleetv1beta1.WorkConditionTypeApplied,
							Status:             metav1.ConditionFalse,
							ObservedGeneration: workGeneration,
						},
					},
				},
			},
			want: []fleetv1beta1.FailedResourcePlacement{
				{
					ResourceIdentifier: fleetv1beta1.ResourceIdentifier{
						Group:     "",
						Version:   "v1",
						Kind:      "ConfigMap",
						Name:      "config-name",
						Namespace: "config-namespace",
						Envelope: &fleetv1beta1.EnvelopeIdentifier{
							Name:      "test-env",
							Namespace: "test-env-ns",
							Type:      "pod",
						},
					},
					Condition: metav1.Condition{
						Type:   fleetv1beta1.WorkConditionTypeApplied,
						Status: metav1.ConditionFalse,
					},
				},
			},
		},
		{
			name: "apply condition is unknown",
			work: fleetv1beta1.Work{
				ObjectMeta: metav1.ObjectMeta{
					Generation: workGeneration,
				},
				Status: fleetv1beta1.WorkStatus{
					Conditions: []metav1.Condition{
						{
							Type:               fleetv1beta1.WorkConditionTypeApplied,
							Status:             metav1.ConditionUnknown,
							ObservedGeneration: workGeneration,
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "multiple manifests in the failed work",
			work: fleetv1beta1.Work{
				ObjectMeta: metav1.ObjectMeta{
					Generation: workGeneration,
				},
				Status: fleetv1beta1.WorkStatus{
					ManifestConditions: []fleetv1beta1.ManifestCondition{
						{
							Identifier: fleetv1beta1.WorkResourceIdentifier{
								Ordinal:   0,
								Group:     "",
								Version:   "v1",
								Kind:      "ConfigMap",
								Name:      "config-name",
								Namespace: "config-namespace",
							},
							Conditions: []metav1.Condition{
								{
									Type:   fleetv1beta1.WorkConditionTypeApplied,
									Status: metav1.ConditionTrue,
								},
								{
									Type:   fleetv1beta1.WorkConditionTypeAvailable,
									Status: metav1.ConditionFalse,
								},
							},
						},
						{
							Identifier: fleetv1beta1.WorkResourceIdentifier{
								Ordinal:   1,
								Group:     "",
								Version:   "v1",
								Kind:      "Service",
								Name:      "svc-name",
								Namespace: "svc-namespace",
							},
							Conditions: []metav1.Condition{
								{
									Type:   fleetv1beta1.WorkConditionTypeApplied,
									Status: metav1.ConditionTrue,
								},
								{
									Type:   fleetv1beta1.WorkConditionTypeAvailable,
									Status: metav1.ConditionFalse,
								},
							},
						},
					},
					Conditions: []metav1.Condition{
						{
							Type:               fleetv1beta1.WorkConditionTypeApplied,
							Status:             metav1.ConditionTrue,
							ObservedGeneration: workGeneration,
						},
						{
							Type:               fleetv1beta1.WorkConditionTypeAvailable,
							Status:             metav1.ConditionFalse,
							ObservedGeneration: workGeneration,
						},
					},
				},
			},
			want: []fleetv1beta1.FailedResourcePlacement{
				{
					ResourceIdentifier: fleetv1beta1.ResourceIdentifier{
						Group:     "",
						Version:   "v1",
						Kind:      "ConfigMap",
						Name:      "config-name",
						Namespace: "config-namespace",
					},
					Condition: metav1.Condition{
						Type:   fleetv1beta1.WorkConditionTypeAvailable,
						Status: metav1.ConditionFalse,
					},
				},
				{
					ResourceIdentifier: fleetv1beta1.ResourceIdentifier{
						Group:     "",
						Version:   "v1",
						Kind:      "Service",
						Name:      "svc-name",
						Namespace: "svc-namespace",
					},
					Condition: metav1.Condition{
						Type:   fleetv1beta1.WorkConditionTypeAvailable,
						Status: metav1.ConditionFalse,
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractFailedResourcePlacementsFromWork(&tc.work)
			if diff := cmp.Diff(tc.want, got, statusCmpOptions...); diff != "" {
				t.Errorf("extractFailedResourcePlacementsFromWork() status mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestIsFailedResourcePlacementsEqual(t *testing.T) {
	time1 := metav1.NewTime(time.Now())
	time2 := metav1.NewTime(time1.Add(-1 * time.Hour))
	tests := []struct {
		name string
		old  []fleetv1beta1.FailedResourcePlacement
		new  []fleetv1beta1.FailedResourcePlacement
		want bool
	}{
		{
			name: "compare two failed resource placement",
			old: []fleetv1beta1.FailedResourcePlacement{
				{
					ResourceIdentifier: fleetv1beta1.ResourceIdentifier{
						Group:     "core",
						Version:   "v1",
						Kind:      "Pod",
						Name:      "pod2",
						Namespace: "default",
						Envelope:  nil,
					},
					Condition: metav1.Condition{
						Type:               condition.ApplyFailedReason,
						Status:             metav1.ConditionTrue,
						ObservedGeneration: 0,
						LastTransitionTime: time1,
						Reason:             work.ManifestApplyFailedReason,
						Message:            "message1",
					},
				},
			},
			new: []fleetv1beta1.FailedResourcePlacement{
				{
					ResourceIdentifier: fleetv1beta1.ResourceIdentifier{
						Group:     "core",
						Version:   "v1",
						Kind:      "Pod",
						Name:      "pod2",
						Namespace: "default",
						Envelope:  nil,
					},
					Condition: metav1.Condition{
						Type:               condition.ApplyFailedReason,
						Status:             metav1.ConditionTrue,
						ObservedGeneration: 0,
						LastTransitionTime: time2,
						Reason:             work.ManifestApplyFailedReason,
						Message:            "message2",
					},
				},
			},
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isFailedResourcePlacementsEqual(tc.old, tc.new)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("isFailedResourcePlacementsEqual() status mismatch (-want, +got):\n%s", diff)
			}
		})
	}

}
