package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/jeanlopezxyz/mcp-prometheus/pkg/kubernetes"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/prometheus"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/toolsets"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/version"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type options struct {
	Version       bool
	LogLevel      int
	Port          string
	URL           string
	Namespace     string
	Service       string
	ServicePort   string
	ServiceScheme string
	Kubeconfig    string
}

func main() {
	flags := pflag.NewFlagSet(version.BinaryName, pflag.ExitOnError)
	pflag.CommandLine = flags

	root := newCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func newCommand() *cobra.Command {
	o := &options{}

	cmd := &cobra.Command{
		Use:   version.BinaryName,
		Short: "Prometheus Model Context Protocol (MCP) server",
		Long:  "A Model Context Protocol (MCP) server that provides Prometheus monitoring tools for AI assistants.",
		RunE: func(c *cobra.Command, args []string) error {
			if o.Version {
				fmt.Fprintln(c.OutOrStdout(), version.Version)
				return nil
			}
			return o.run()
		},
	}

	cmd.Flags().BoolVar(&o.Version, "version", false, "Print version information and quit")
	cmd.Flags().IntVar(&o.LogLevel, "log-level", 0, "Set the log level (from 0 to 9)")
	cmd.Flags().StringVar(&o.Port, "port", "", "Start a streamable HTTP server on the specified port (e.g. 8080)")
	cmd.Flags().StringVar(&o.URL, "url", "", "Direct Prometheus URL (e.g. http://localhost:9090). Overrides K8S auto-detect. Env: PROMETHEUS_URL")
	cmd.Flags().StringVar(&o.Namespace, "namespace", "", "Kubernetes namespace for Prometheus service (default: openshift-monitoring)")
	cmd.Flags().StringVar(&o.Service, "service", "", "Kubernetes service name for Prometheus (default: prometheus-operated)")
	cmd.Flags().StringVar(&o.ServicePort, "service-port", "", "Kubernetes service port for Prometheus (default: 9090)")
	cmd.Flags().StringVar(&o.ServiceScheme, "service-scheme", "", "Kubernetes service scheme: http or https (default: https)")
	cmd.Flags().StringVar(&o.Kubeconfig, "kubeconfig", "", "Path to kubeconfig file (default: auto-detect)")

	return cmd
}

func (o *options) run() error {
	o.initializeLogging()

	klog.V(1).Infof("Starting %s %s", version.BinaryName, version.Version)

	baseURL, httpClient, err := o.resolveConnection()
	if err != nil {
		return fmt.Errorf("failed to resolve Prometheus connection: %w", err)
	}

	client := prometheus.NewClient(baseURL, httpClient)

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:       version.BinaryName,
			Title:      version.BinaryName,
			Version:    version.Version,
			WebsiteURL: version.WebsiteURL,
		},
		&mcp.ServerOptions{
			Capabilities: &mcp.ServerCapabilities{
				Tools:   &mcp.ToolCapabilities{ListChanged: true},
				Logging: &mcp.LoggingCapabilities{},
			},
		},
	)

	toolsets.RegisterAll(server, client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown on signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		klog.V(1).Info("Shutting down...")
		cancel()
	}()

	if o.Port != "" {
		klog.V(1).Infof("Starting HTTP server on port %s", o.Port)
		handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
			return server
		}, &mcp.StreamableHTTPOptions{})
		httpServer := &http.Server{
			Addr:    ":" + o.Port,
			Handler: handler,
		}
		go func() {
			<-ctx.Done()
			httpServer.Close()
		}()
		return httpServer.ListenAndServe()
	}

	klog.V(1).Info("Starting stdio transport")
	return server.Run(ctx, &mcp.StdioTransport{})
}

func (o *options) initializeLogging() {
	flagSet := flag.NewFlagSet("klog", flag.ContinueOnError)
	klog.InitFlags(flagSet)
	if o.Port == "" {
		// Disable klog output for stdio mode to avoid breaking the protocol
		_ = flagSet.Parse([]string{"-logtostderr=false", "-alsologtostderr=false", "-stderrthreshold=FATAL"})
		return
	}
	if o.LogLevel >= 0 {
		_ = flagSet.Parse([]string{"--v", strconv.Itoa(o.LogLevel)})
	}
}

// resolveConnection determines how to connect to Prometheus.
// Priority: --url flag / PROMETHEUS_URL env → K8S auto-detect → error.
func (o *options) resolveConnection() (string, *http.Client, error) {
	// 1. Direct URL: flag takes precedence, then env var
	url := o.URL
	if url == "" {
		url = os.Getenv("PROMETHEUS_URL")
	}
	if url != "" {
		klog.V(1).Infof("Using direct Prometheus URL: %s", url)
		return url, nil, nil
	}

	// 2. K8S auto-detect via kubeconfig or in-cluster
	if kubernetes.CanConnectToCluster(o.Kubeconfig) {
		namespace := kubernetes.DetectNamespace(o.Namespace, "openshift-monitoring")
		service := o.Service
		if service == "" {
			service = "prometheus-operated"
		}
		port := o.ServicePort
		if port == "" {
			port = "9090"
		}
		scheme := o.ServiceScheme
		if scheme == "" {
			scheme = "https"
		}
		klog.V(1).Infof("Auto-detected Kubernetes cluster, connecting via API proxy: %s/%s:%s:%s", namespace, scheme, service, port)
		return kubernetes.NewK8SProxyClient(o.Kubeconfig, namespace, service, port, scheme)
	}

	// 3. Nothing available → error with helpful examples
	return "", nil, fmt.Errorf(`no Prometheus connection available

No direct URL provided and no Kubernetes cluster detected.

Configure one of the following:

  # Direct URL (testing/dev)
  %[1]s --url http://localhost:9090

  # Environment variable
  PROMETHEUS_URL=http://prometheus:9090 %[1]s

  # Kubernetes auto-detect with defaults (openshift-monitoring/prometheus-operated:9090)
  # Requires a valid kubeconfig or in-cluster service account
  %[1]s

  # Custom namespace/service
  %[1]s --namespace openshift-monitoring --service prometheus-operated

  # Explicit kubeconfig
  %[1]s --kubeconfig /path/to/kubeconfig`, version.BinaryName)
}
