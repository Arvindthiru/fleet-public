package clusterresourcebindingwatcher

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	fleetv1beta1 "go.goms.io/fleet/apis/placement/v1beta1"
	"go.goms.io/fleet/pkg/utils/controller"
)

// Reconciler reconciles updates to clusterResourceBinding.
type Reconciler struct {
	// Client is the client the controller uses to access the hub cluster.
	client.Client
	// PlacementController maintains a rate limited queue which used to store
	// the name of the clusterResourcePlacement and a reconcile function to consume the items in queue.
	PlacementController controller.Controller
}

// Reconcile reconciles the clusterResourceBinding.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	bindingRef := klog.KRef("", req.Name)
	var binding fleetv1beta1.ClusterResourceBinding
	if err := r.Client.Get(ctx, req.NamespacedName, &binding); err != nil {
		klog.ErrorS(err, "Failed to get cluster resource binding", "clusterResourceBinding", bindingRef)
		return ctrl.Result{}, controller.NewAPIServerError(true, client.IgnoreNotFound(err))
	}

	// Check if the cluster resource binding has been deleted.
	// Normally this would not happen as the event filter is set to filter out all deletion events.
	if binding.DeletionTimestamp != nil {
		// The cluster resource binding has been deleted; ignore it.
		return ctrl.Result{}, nil
	}

	// Verify if the policy snapshot is currently active.
	crpName, ok := binding.Labels[fleetv1beta1.CRPTrackingLabel]
	if !ok {
		// The CRPTrackingLabel label is not present; normally this should never occur.
		klog.ErrorS(controller.NewUnexpectedBehaviorError(fmt.Errorf("CRPTrackingLabel is missing")),
			"CRPTrackingLabel is not present",
			"clusterResourceBinding", bindingRef)
		// This is not a situation that the controller can recover by itself. Should the label
		// value be corrected, the controller will be triggered again.
		return ctrl.Result{}, nil
	}

	// Enqueue the CRP name for reconciling.
	r.PlacementController.Enqueue(crpName)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	customPredicate := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			// Ignore creation events.
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// Ignore deletion events.
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Check if the update event is valid.
			if e.ObjectOld == nil || e.ObjectNew == nil {
				err := controller.NewUnexpectedBehaviorError(fmt.Errorf("update event is invalid"))
				klog.ErrorS(err, "Failed to process update event")
				return false
			}
			oldBinding, oldOk := e.ObjectOld.(*fleetv1beta1.ClusterResourceBinding)
			newBinding, newOk := e.ObjectNew.(*fleetv1beta1.ClusterResourceBinding)
			if !oldOk || !newOk {
				err := controller.NewUnexpectedBehaviorError(fmt.Errorf("failed to cast runtime objects in update event to cluster resource binding objects"))
				klog.ErrorS(err, "Failed to process update event")
				return false
			}

			return isGenerationUpdated(oldBinding, newBinding) ||
				isConditionUpdated(oldBinding.GetCondition(string(fleetv1beta1.ResourceBindingRolloutStarted)), newBinding.GetCondition(string(fleetv1beta1.ResourceBindingRolloutStarted))) ||
				isConditionUpdated(oldBinding.GetCondition(string(fleetv1beta1.ResourceBindingOverridden)), newBinding.GetCondition(string(fleetv1beta1.ResourceBindingOverridden))) ||
				isConditionUpdated(oldBinding.GetCondition(string(fleetv1beta1.ResourceBindingWorkCreated)), newBinding.GetCondition(string(fleetv1beta1.ResourceBindingWorkCreated))) ||
				isConditionUpdated(oldBinding.GetCondition(string(fleetv1beta1.ResourceBindingApplied)), newBinding.GetCondition(string(fleetv1beta1.ResourceBindingApplied))) ||
				isConditionUpdated(oldBinding.GetCondition(string(fleetv1beta1.ResourceBindingAvailable)), newBinding.GetCondition(string(fleetv1beta1.ResourceBindingAvailable)))
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&fleetv1beta1.ClusterResourceBinding{}).
		WithEventFilter(customPredicate).
		Complete(r)
}

func isGenerationUpdated(oldBinding, newBinding *fleetv1beta1.ClusterResourceBinding) bool {
	return oldBinding.Generation != newBinding.Generation
}

func isConditionUpdated(oldCondition, newCondition *metav1.Condition) bool {
	if oldCondition == nil && newCondition == nil {
		return false
	}
	if oldCondition == nil || newCondition == nil {
		return true
	}
	return oldCondition.Status != newCondition.Status
}