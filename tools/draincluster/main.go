/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package main

import (
	"context"
	"flag"
	"log"

	"k8s.io/apimachinery/pkg/runtime"

	clusterv1beta1 "go.goms.io/fleet/apis/cluster/v1beta1"
	placementv1beta1 "go.goms.io/fleet/apis/placement/v1beta1"
	"go.goms.io/fleet/tools/draincluster/drain"
	toolsutils "go.goms.io/fleet/tools/utils"
)

func main() {
	//TODO (arvindth): add flags for timeout, help for program.
	scheme := runtime.NewScheme()
	ctx := context.Background()

	hubClusterContext := flag.String("hubClusterContext", "", "the kubectl context for the hub cluster")
	clusterName := flag.String("clusterName", "", "name of the cluster to cordon")
	flag.Parse()

	if *hubClusterContext == "" {
		log.Fatalf("hub cluster context for kubectl cannot be empty")
	}
	if *clusterName == "" {
		log.Fatalf("cluster name to cordon cannot be empty")
	}

	if err := clusterv1beta1.AddToScheme(scheme); err != nil {
		log.Fatalf("failed to add custom APIs (cluster) to the runtime scheme: %v", err)
	}
	if err := placementv1beta1.AddToScheme(scheme); err != nil {
		log.Fatalf("failed to add custom APIs (placement) to the runtime scheme: %v", err)
	}

	hubClient, err := toolsutils.GetClusterClientFromClusterContext(*hubClusterContext, scheme)
	if err != nil {
		log.Fatalf("failed to create hub cluster client: %v", err)
	}

	drainClusterHelper := drain.Helper{
		HubClient:                            hubClient,
		ClusterName:                          *clusterName,
		ClusterResourcePlacementResourcesMap: make(map[string][]placementv1beta1.ResourceIdentifier),
	}

	isDrainSuccessful, err := drainClusterHelper.Drain(ctx)
	if err != nil {
		log.Fatalf("failed to drain member cluster %s: %v", drainClusterHelper.ClusterName, err)
	}

	if isDrainSuccessful {
		log.Printf("drain was successful for cluster %s", *clusterName)
	} else {
		log.Printf("drain was not successful for cluster %s", *clusterName)
	}

	log.Printf("retrying drain to ensure all resources propagated from hub cluster are evicted")

	// reset ClusterResourcePlacementResourcesMap for retry.
	drainClusterHelper.ClusterResourcePlacementResourcesMap = map[string][]placementv1beta1.ResourceIdentifier{}
	isDrainRetrySuccessful, err := drainClusterHelper.Drain(ctx)
	if err != nil {
		log.Fatalf("failed to drain cluster on retry %s: %v", drainClusterHelper.ClusterName, err)
	}
	if isDrainRetrySuccessful {
		log.Printf("drain retry was successful for cluster %s", *clusterName)
	} else {
		log.Printf("drain retry was not successful for cluster %s", *clusterName)
	}
}
