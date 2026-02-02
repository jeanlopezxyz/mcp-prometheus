package kubernetes

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewK8SProxyClient creates an HTTP client and base URL for proxying through
// the Kubernetes API server to a service.
// URL pattern: {host}/api/v1/namespaces/{ns}/services/{svc}:{port}/proxy
func NewK8SProxyClient(kubeconfig, namespace, service, port string) (string, *http.Client, error) {
	config, err := getRESTConfig(kubeconfig)
	if err != nil {
		return "", nil, fmt.Errorf("kubernetes config: %w", err)
	}

	transport, err := rest.TransportFor(config)
	if err != nil {
		return "", nil, fmt.Errorf("kubernetes transport: %w", err)
	}

	host := config.Host
	baseURL := fmt.Sprintf("%s/api/v1/namespaces/%s/services/%s:%s/proxy",
		host, namespace, service, port)

	httpClient := &http.Client{Transport: transport}
	return baseURL, httpClient, nil
}

// DetectNamespace returns the best namespace to use.
// Priority: explicit flag → in-cluster namespace file → defaultNS.
func DetectNamespace(explicit, defaultNS string) string {
	if explicit != "" {
		return explicit
	}

	// In-cluster: read namespace from service account
	if data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if ns := strings.TrimSpace(string(data)); ns != "" {
			return ns
		}
	}

	return defaultNS
}

// CanConnectToCluster returns true if a Kubernetes REST config can be loaded
// from the given kubeconfig path, in-cluster config, or default kubeconfig rules.
func CanConnectToCluster(kubeconfig string) bool {
	_, err := getRESTConfig(kubeconfig)
	return err == nil
}

// getRESTConfig attempts to load Kubernetes config.
// Strategy: explicit kubeconfig path → in-cluster config → default kubeconfig rules.
func getRESTConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
		configOverrides := &clientcmd.ConfigOverrides{}
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
		return kubeConfig.ClientConfig()
	}

	// Try in-cluster config first (when running inside a pod)
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Fall back to default kubeconfig rules
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	return kubeConfig.ClientConfig()
}
