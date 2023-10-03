/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package e2e

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	placementv1beta1 "go.goms.io/fleet/apis/placement/v1beta1"
	"go.goms.io/fleet/test/e2e/framework"
)

func validateWorkNamespaceOnCluster(cluster *framework.Cluster, name types.NamespacedName) error {
	ns := &corev1.Namespace{}
	if err := cluster.KubeClient.Get(ctx, name, ns); err != nil {
		return err
	}

	// Use the object created in the hub cluster as reference; this helps to avoid the trouble
	// of having to ignore default fields in the spec.
	wantNS := &corev1.Namespace{}
	if err := hubClient.Get(ctx, name, wantNS); err != nil {
		return err
	}

	if diff := cmp.Diff(
		ns, wantNS,
		ignoreNamespaceStatusField,
		ignoreObjectMetaAutoGeneratedFields,
		ignoreObjectMetaAnnotationField,
	); diff != "" {
		return fmt.Errorf("work namespace diff (-got, +want): %s", diff)
	}
	return nil
}

func validateConfigMapOnCluster(cluster *framework.Cluster, name types.NamespacedName) error {
	configMap := &corev1.ConfigMap{}
	if err := cluster.KubeClient.Get(ctx, name, configMap); err != nil {
		return err
	}

	// Use the object created in the hub cluster as reference.
	wantConfigMap := &corev1.ConfigMap{}
	if err := hubClient.Get(ctx, name, wantConfigMap); err != nil {
		return err
	}

	if diff := cmp.Diff(
		configMap, wantConfigMap,
		ignoreObjectMetaAutoGeneratedFields,
		ignoreObjectMetaAnnotationField,
	); diff != "" {
		return fmt.Errorf("app deployment diff (-got, +want): %s", diff)
	}

	return nil
}

func workNamespaceAndConfigMapPlacedOnClusterActual(cluster *framework.Cluster) func() error {
	workNamespaceName := fmt.Sprintf(workNamespaceNameTemplate, GinkgoParallelProcess())
	appConfigMapName := fmt.Sprintf(appConfigMapNameTemplate, GinkgoParallelProcess())

	return func() error {
		if err := validateWorkNamespaceOnCluster(cluster, types.NamespacedName{Name: workNamespaceName}); err != nil {
			return err
		}

		return validateConfigMapOnCluster(cluster, types.NamespacedName{Namespace: workNamespaceName, Name: appConfigMapName})
	}
}

func crpRolloutCompletedConditions(generation int64) []metav1.Condition {
	return []metav1.Condition{
		{
			Type:               string(placementv1beta1.ClusterResourcePlacementScheduledConditionType),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: generation,
		},
		{
			Type:               string(placementv1beta1.ClusterResourcePlacementSynchronizedConditionType),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: generation,
		},
		{
			Type:               string(placementv1beta1.ClusterResourcePlacementAppliedConditionType),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: generation,
		},
	}
}

func resourcePlacementRolloutCompletedConditions(generation int64) []metav1.Condition {
	return []metav1.Condition{
		{
			Type:               string(placementv1beta1.ResourceScheduledConditionType),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: generation,
		},
		{
			Type:               string(placementv1beta1.ResourceWorkSynchronizedConditionType),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: generation,
		},
		{
			Type:               string(placementv1beta1.ResourcesAppliedConditionType),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: generation,
		},
	}
}

func validateCRPStatus(name types.NamespacedName, wantSelectedResources []placementv1beta1.ResourceIdentifier) error {
	crp := &placementv1beta1.ClusterResourcePlacement{}
	if err := hubClient.Get(ctx, name, crp); err != nil {
		return err
	}

	wantCRPConditions := crpRolloutCompletedConditions(crp.Generation)
	wantPlacementStatus := []placementv1beta1.ResourcePlacementStatus{
		{
			ClusterName: memberCluster1Name,
			Conditions:  resourcePlacementRolloutCompletedConditions(crp.Generation),
		},
		//{
		//	ClusterName: memberCluster2Name,
		//	Conditions:  resourcePlacementRolloutCompletedConditions(crp.Generation),
		//},
		//{
		//	ClusterName: memberCluster3Name,
		//	Conditions:  resourcePlacementRolloutCompletedConditions(crp.Generation),
		//},
	}
	wantStatus := placementv1beta1.ClusterResourcePlacementStatus{
		Conditions:        wantCRPConditions,
		PlacementStatuses: wantPlacementStatus,
		SelectedResources: wantSelectedResources,
	}
	if diff := cmp.Diff(crp.Status, wantStatus, crpStatusCmpOptions...); diff != "" {
		return fmt.Errorf("CRP status diff (-got, +want): %s", diff)
	}
	return nil
}

func crpStatusUpdatedActual() func() error {
	crpName := fmt.Sprintf(crpNameTemplate, GinkgoParallelProcess())
	workNamespaceName := fmt.Sprintf(workNamespaceNameTemplate, GinkgoParallelProcess())
	appConfigMapName := fmt.Sprintf(appConfigMapNameTemplate, GinkgoParallelProcess())

	return func() error {
		wantSelectedResources := []placementv1beta1.ResourceIdentifier{
			{
				Kind:    "Namespace",
				Name:    workNamespaceName,
				Version: "v1",
			},
			{
				Kind:      "ConfigMap",
				Name:      appConfigMapName,
				Version:   "v1",
				Namespace: workNamespaceName,
			},
		}
		return validateCRPStatus(types.NamespacedName{Name: crpName}, wantSelectedResources)
	}
}

func workNamespaceRemovedFromClusterActual(cluster *framework.Cluster) func() error {
	client := cluster.KubeClient

	workNamespaceName := fmt.Sprintf(workNamespaceNameTemplate, GinkgoParallelProcess())
	return func() error {
		if err := client.Get(ctx, types.NamespacedName{Name: workNamespaceName}, &corev1.Namespace{}); !errors.IsNotFound(err) {
			return fmt.Errorf("work namespace %s still exists or an unexpected error occurred: %w", workNamespaceName, err)
		}
		return nil
	}
}

func allFinalizersExceptForCustomDeletionBlockerRemovedFromCRPActual() func() error {
	crpName := fmt.Sprintf(crpNameTemplate, GinkgoParallelProcess())

	return func() error {
		crp := &placementv1beta1.ClusterResourcePlacement{}
		if err := hubClient.Get(ctx, types.NamespacedName{Name: crpName}, crp); err != nil {
			return err
		}

		wantFinalizers := []string{customDeletionBlockerFinalizer}
		finalizer := crp.Finalizers
		if diff := cmp.Diff(finalizer, wantFinalizers); diff != "" {
			return fmt.Errorf("CRP finalizers diff (-got, +want): %s", diff)
		}

		return nil
	}
}

func crpRemovedActual() func() error {
	crpName := fmt.Sprintf(crpNameTemplate, GinkgoParallelProcess())

	return func() error {
		if err := hubClient.Get(ctx, types.NamespacedName{Name: crpName}, &placementv1beta1.ClusterResourcePlacement{}); !errors.IsNotFound(err) {
			return fmt.Errorf("CRP still exists or an unexpected error occurred: %w", err)
		}

		return nil
	}
}
