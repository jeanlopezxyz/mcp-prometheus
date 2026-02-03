package kubernetes

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewK8SProxyClient creates an HTTP client and base URL for proxying through
// the Kubernetes API server to a service.
// URL pattern: {host}/api/v1/namespaces/{ns}/services/{scheme}:{svc}:{port}/proxy
// The scheme is required for services that use TLS (e.g., OpenShift monitoring).
func NewK8SProxyClient(kubeconfig, namespace, service, port, scheme string) (string, *http.Client, error) {
	config, err := getRESTConfig(kubeconfig)
	if err != nil {
		return "", nil, fmt.Errorf("kubernetes config: %w", err)
	}

	transport, err := rest.TransportFor(config)
	if err != nil {
		return "", nil, fmt.Errorf("kubernetes transport: %w", err)
	}

	host := config.Host
	baseURL := fmt.Sprintf("%s/api/v1/namespaces/%s/services/%s:%s:%s/proxy",
		host, namespace, scheme, service, port)

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

// bearerTokenTransport wraps an http.RoundTripper to add Authorization header.
type bearerTokenTransport struct {
	token string
	base  http.RoundTripper
}

func (t *bearerTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(req2)
}

// NewOpenShiftRouteClient creates an HTTP client and base URL using an OpenShift route.
// This is the preferred method for OpenShift clusters as routes provide authenticated access.
// Returns: routeURL, httpClient, error (nil error with empty URL if no route found)
func NewOpenShiftRouteClient(kubeconfig, namespace, routeName string) (string, *http.Client, error) {
	config, err := getRESTConfig(kubeconfig)
	if err != nil {
		return "", nil, fmt.Errorf("kubernetes config: %w", err)
	}

	// Get bearer token from config
	token := config.BearerToken
	if token == "" && config.BearerTokenFile != "" {
		data, err := os.ReadFile(config.BearerTokenFile)
		if err != nil {
			return "", nil, fmt.Errorf("reading bearer token file: %w", err)
		}
		token = strings.TrimSpace(string(data))
	}
	if token == "" {
		// No bearer token available, cannot use route authentication
		return "", nil, nil
	}

	// Create dynamic client to query routes
	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return "", nil, fmt.Errorf("creating dynamic client: %w", err)
	}

	// Query for the route
	routeGVR := schema.GroupVersionResource{
		Group:    "route.openshift.io",
		Version:  "v1",
		Resource: "routes",
	}

	route, err := dynClient.Resource(routeGVR).Namespace(namespace).Get(context.Background(), routeName, metav1.GetOptions{})
	if err != nil {
		// Route not found or not OpenShift - return empty to fall back to other methods
		return "", nil, nil
	}

	// Extract host from route spec
	host, found, err := unstructured.NestedString(route.Object, "spec", "host")
	if err != nil || !found || host == "" {
		return "", nil, nil
	}

	// Build route URL (routes always use HTTPS)
	routeURL := "https://" + host

	// Create HTTP client with bearer token and TLS (skip verify for self-signed certs)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: &bearerTokenTransport{
			token: token,
			base:  transport,
		},
	}

	return routeURL, httpClient, nil
}

// IsOpenShift checks if the cluster is OpenShift by checking for route.openshift.io API.
func IsOpenShift(kubeconfig string) bool {
	config, err := getRESTConfig(kubeconfig)
	if err != nil {
		return false
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return false
	}

	// Try to list routes - if this works, it's OpenShift
	routeGVR := schema.GroupVersionResource{
		Group:    "route.openshift.io",
		Version:  "v1",
		Resource: "routes",
	}

	_, err = dynClient.Resource(routeGVR).List(context.Background(), metav1.ListOptions{Limit: 1})
	return err == nil
}
