package main

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"log"
	"os"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"go.goms.io/fleet/pkg/utils"
	"go.goms.io/fleet/test/e2e/framework"
)

const (
	hubClusterName           = "kind-hub"
	kubeConfigPathEnvVarName = "KUBECONFIG"
)

var (
	ctx    = context.Background()
	scheme = runtime.NewScheme()

	hubCluster *framework.Cluster
	hubClient  client.Client
)

func main() {
	// Add built-in APIs and extensions to the scheme.
	if err := k8sscheme.AddToScheme(scheme); err != nil {
		log.Fatalf("failed to add built-in APIs to the runtime scheme: %v", err)
	}
	if err := apiextensionsv1.AddToScheme(scheme); err != nil {
		log.Fatalf("failed to add API extensions to the runtime scheme: %v", err)
	}
	// Check if the required environment variable, which specifies the path to kubeconfig file, has been set.
	kubeConfigVal := os.Getenv(kubeConfigPathEnvVarName)
	if kubeConfigVal == "" {
		fmt.Println("Required environment variable KUBECONFIG is not set")
		os.Exit(0)
	}

	// Initialize the cluster objects and their clients.
	hubCluster = framework.NewCluster(hubClusterName, scheme)
	if hubCluster == nil {
		fmt.Println("hubCluster is nil")
		os.Exit(0)
	}
	clusterConfig := framework.GetClientConfig(hubCluster)
	restConfig, err := clusterConfig.ClientConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	hubCluster.KubeClient, err = client.New(restConfig, client.Options{Scheme: hubCluster.Scheme})
	hubClient = hubCluster.KubeClient
	var largeSecret v1.Secret
	utils.GetObjectFromManifest("./test/integration/manifests/resources/test-large-secret.yaml", &largeSecret)
	fmt.Println(largeSecret.Name)
	err = hubCluster.KubeClient.Create(ctx, &largeSecret)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	size, err := getObjectSize(largeSecret)
	if err != nil {
		fmt.Println("failed to get object size")
	}
	fmt.Println("object size is:", size)
}

func getObjectSize(obj interface{}) (int, error) {
	// Marshal the object to JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return 0, err
	}

	// Calculate the size in bytes
	size := len(data)
	return size, nil
}
