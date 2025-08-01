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

// Package membercluster features a controller that enqueues CRPs on member cluster changes.
package membercluster

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	clusterv1beta1 "go.goms.io/fleet/apis/cluster/v1beta1"
	placementv1beta1 "go.goms.io/fleet/apis/placement/v1beta1"
	"go.goms.io/fleet/pkg/scheduler/clustereligibilitychecker"
	"go.goms.io/fleet/pkg/scheduler/queue"
	"go.goms.io/fleet/pkg/utils/controller"
)

// Reconciler is the member cluster controller reconciler.
type Reconciler struct {
	// Client is a (cached) client for accessing the Kubernetes API server.
	Client client.Client

	// SchedulerWorkQueue is the work queue for the scheduler.
	SchedulerWorkQueue queue.PlacementSchedulingQueueWriter

	// clusterEligibilityCheck helps check if a cluster is eligible for resource replacement.
	ClusterEligibilityChecker *clustereligibilitychecker.ClusterEligibilityChecker
}

// Reconcile reconciles a member cluster.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	memberClusterRef := klog.KRef(req.Namespace, req.Name)
	startTime := time.Now()
	klog.V(2).InfoS("Reconciliation starts", "memberCluster", memberClusterRef)
	defer func() {
		latency := time.Since(startTime).Milliseconds()
		klog.V(2).InfoS("Reconciliation ends", "memberCluster", memberClusterRef, "latency", latency)
	}()

	// Generally speaking, in a fleet there might be the following two types of cluster changes:
	//
	//  1. a cluster, originally ineligible for resource placement for some reason, becomes eligible;
	//
	//     It may happen for 2 reasons:
	//
	//     a) the cluster's setup (e.g., its labels) or status (e.g., resource/non-resource properties),
	//        has changed; and/or
	//     b) an unexpected development which originally leads the scheduler to disregard the cluster
	//     (e.g., agents not joining, network partition, etc.) has been resolved.
	//     c) the cluster, has a taint removed from it and now is eligible for scheduling.
	//
	//  2. a cluster, originally eligible for resource placement, becomes ineligible for some reason.
	//
	//     Similarly, it may happen for 2 reasons:
	//
	//     a) the cluster's setup (e.g., its labels) or status (e.g., resource/non-resource properties),
	//        has changed; and/or
	//     b) an unexpected development (e.g., agents failing, network partition, etc.) has occurred.
	//     c) the cluster, which may or may not have resources placed on it, has left the fleet (deleted).
	//
	// Among the cases,
	//
	// * 1a), 1b) and 1c) require attention on the scheduler's end, specifically:
	//   - CRPs of the PickAll placement type may be able to select this cluster now;
	//   - CRPs of the PickN placement type, which have not been fully scheduled yet, may be
	//     able to select this cluster, and gets a step closer to being fully scheduled;
	//   - CRPs of the PickFixed placement type, which have not been fully scheduled yet, may
	//     be able to select this cluster, and gets a step closer to being fully scheduled, 1c)
	//     doesn't apply to this scenario since taints are not honored for PickFixed CRPs.
	//
	// * 2a) and 2b) require no attention on the scheduler's end, specifically:
	//   - CRPs which have already selected this cluster, regardless of its placement type, cannot
	//     deselect it; resources are only removed from the cluster if the user explicitly requires
	//     so, for example, by specifying a new scheduling policy, this helps reduce fluctuations
	//     in the system (i.e., ignoredDuringExecution semantics).
	//   - CRPs which have not selected this cluster, regardless of its placement type, cannot
	//     select it either, as it does not meet the requirement.
	//
	// (Note also that from this controller's perspective, we cannot reliably tell the difference
	//  between 1a) and 2a).)
	//
	// * 2c) requires attention on the scheduler's end, specifically:
	//   - All CRPs, which have already selected this cluster, regardless of its placement type,
	//     must deselect it, as the binding is no longer valid (dangling). CRPs of the PickN type
	//     may further need to pick another cluster as replacement.
	//
	// This controller is set to handle cases 1a), 1b), and 2c). Note that it is only guaranteed
	// that this controller will not emit false negatives, i.e., all the changes that require
	// the scheduler's attention will be captured; in other words, false positives may still
	// happen, i.e., this controller may trigger the scheduler to run a scheduling loop even though
	// there is no change needed. In such situations, the scheduler will be the final judge and
	// process the CRPs accordingly.

	// Retrieve the member cluster.
	memberCluster := &clusterv1beta1.MemberCluster{}
	memberClusterKey := types.NamespacedName{Name: req.Name}
	isMemberClusterMissing := false

	memberClusterGetErr := r.Client.Get(ctx, memberClusterKey, memberCluster)
	switch {
	case errors.IsNotFound(memberClusterGetErr):
		// This could happen when the member cluster is deleted. In this case, controller will request the scheduler to check
		// all CRPs just in case.
		isMemberClusterMissing = true
	case memberClusterGetErr != nil:
		klog.ErrorS(memberClusterGetErr, "Failed to get member cluster", "memberCluster", memberClusterRef)
		return ctrl.Result{}, controller.NewAPIServerError(true, memberClusterGetErr)
		// Do nothing if there is no error returned.
	}

	// List all CRPs.
	//
	// Note that this controller reads CRPs from the same cache as the scheduler.
	crpList := &placementv1beta1.ClusterResourcePlacementList{}
	if err := r.Client.List(ctx, crpList); err != nil {
		klog.ErrorS(err, "Failed to list CRPs", "memberCluster", memberClusterRef)
		return ctrl.Result{}, controller.NewAPIServerError(true, err)
	}

	crps := crpList.Items
	if !isMemberClusterMissing && memberCluster.GetDeletionTimestamp().IsZero() {
		// If the member cluster is set to the left state, the scheduler needs to process all
		// CRPs (case 2c)); otherwise, only CRPs of the PickAll type + CRPs of the PickN type,
		// which have not been fully scheduled, need to be processed (case 1a) and 1b)).
		crps = classifyCRPs(crpList.Items)
	}

	// Enqueue the CRPs.
	//
	// Note that all the CRPs in the system are enqueued; technically speaking, for situation
	// 1a), 1b) and 1c), PickN CRPs that have been fully scheduled needs no further processing, however,
	// for simplicity reasons, this controller will not distinguish between the cases.
	for idx := range crps {
		crp := &crps[idx]
		klog.V(2).InfoS(
			"Enqueueing CRP for scheduler processing",
			"memberCluster", memberClusterRef,
			"clusterResourcePlacement", klog.KObj(crp))
		r.SchedulerWorkQueue.Add(queue.PlacementKey(crp.Name))
	}

	// The reconciliation loop completes.
	return ctrl.Result{}, nil
}

// SetupWithManager builds a controller with Reconciler and sets it up with a controller manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	customPredicate := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			// Normally it is safe to ignore newly created cluster objects, as they are not yet
			// ready for scheduling; when the clusters do become ready, the controller will catch
			// an update event.
			//
			// Furthermore, when the controller restarts, the scheduler is set to reconcile all
			// CRPs anyway, which will account for any missing updates on the cluster side
			// during the downtime; in other words, notifications from this controller is not
			// necessary.
			klog.V(3).InfoS("Ignoring create events for member cluster objects", "eventObject", klog.KObj(e.Object))
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// We only notify the scheduler when a member cluster is deleted which means the member agent has finished the leaving process.
			klog.V(2).InfoS("Member cluster object is deleted", "eventObject", klog.KObj(e.Object))
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Check if the update event is valid.
			if e.ObjectOld == nil || e.ObjectNew == nil {
				err := controller.NewUnexpectedBehaviorError(fmt.Errorf("update event is invalid"))
				klog.ErrorS(err, "Failed to process update event")
				return false
			}

			oldCluster, oldOk := e.ObjectOld.(*clusterv1beta1.MemberCluster)
			newCluster, newOk := e.ObjectNew.(*clusterv1beta1.MemberCluster)
			if !oldOk || !newOk {
				err := controller.NewUnexpectedBehaviorError(fmt.Errorf("failed to cast runtime objects in update event to member cluster objects"))
				klog.ErrorS(err, "Failed to process update event")
				return false
			}

			clusterKObj := klog.KObj(newCluster)
			// The cluster is being deleted.
			if !newCluster.GetDeletionTimestamp().IsZero() {
				klog.V(2).InfoS("A member cluster is leaving the fleet, ignore until it is left", "memberCluster", clusterKObj)
				return false
			}

			// Capture label changes.
			//
			// Note that the controller runs only when label changes happen on joined clusters.
			if !reflect.DeepEqual(oldCluster.Labels, newCluster.Labels) {
				klog.V(2).InfoS("A member cluster label change has been detected", "memberCluster", clusterKObj)
				return true
			}

			// Capture taint update/delete changes.
			if isTaintsUpdatedOrDeleted(oldCluster.Spec.Taints, newCluster.Spec.Taints) {
				klog.V(2).InfoS("A member cluster taint update/delete has been detected", "memberCluster", clusterKObj)
				return true
			}

			// Capture non-resource property changes.
			//
			// Observation time refreshes is not considered as a change.
			oldProperties := oldCluster.Status.Properties
			newProperties := newCluster.Status.Properties
			if len(oldProperties) != len(newProperties) {
				return true
			}
			for oldK, oldV := range oldProperties {
				newV, ok := newProperties[oldK]
				if !ok || oldV.Value != newV.Value {
					return true
				}
			}

			// Capture resource usage changes.
			oldCapacity := oldCluster.Status.ResourceUsage.Capacity
			newCapacity := newCluster.Status.ResourceUsage.Capacity
			if !equality.Semantic.DeepEqual(oldCapacity, newCapacity) {
				return true
			}

			oldAllocatable := oldCluster.Status.ResourceUsage.Allocatable
			newAllocatable := newCluster.Status.ResourceUsage.Allocatable
			if !equality.Semantic.DeepEqual(oldAllocatable, newAllocatable) {
				return true
			}

			oldAvailable := oldCluster.Status.ResourceUsage.Available
			newAvailable := newCluster.Status.ResourceUsage.Available
			if !equality.Semantic.DeepEqual(oldAvailable, newAvailable) {
				return true
			}

			// Check the resource placement eligibility for the old and new cluster object.
			oldEligible, _ := r.ClusterEligibilityChecker.IsEligible(oldCluster)
			newEligible, _ := r.ClusterEligibilityChecker.IsEligible(newCluster)

			if !oldEligible && newEligible {
				// The cluster becomes eligible for resource placement, i.e., match for case 1b).
				//
				// The reverse, i.e., eligible -> ineligible, is ignored (case 2b)).
				klog.V(2).InfoS("A member cluster may become eligible for resource placement", "memberCluster", clusterKObj)
				return true
			}

			// All the other changes are ignored.
			klog.V(3).InfoS("Ignoring update events that are irrelevant to the scheduler", "memberCluster", clusterKObj)
			return false
		},
	}

	return ctrl.NewControllerManagedBy(mgr).Named("membercluster-scheduler-watcher").
		For(&clusterv1beta1.MemberCluster{}).
		WithEventFilter(customPredicate).
		Complete(r)
}

func isTaintsUpdatedOrDeleted(oldTaints []clusterv1beta1.Taint, newTaints []clusterv1beta1.Taint) bool {
	newTaintsMap := make(map[clusterv1beta1.Taint]bool)
	for _, newTaint := range newTaints {
		newTaintsMap[newTaint] = true
	}
	for _, oldTaint := range oldTaints {
		if !newTaintsMap[oldTaint] {
			return true
		}
	}
	return false
}
