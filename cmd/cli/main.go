package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"

	k8sclient "kubeapp/internal/k8s"
)

// globals shared across command files in package main
var (
	kubeconfig    string
	namespace     string
	allNamespaces bool
	client        *kubernetes.Clientset
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "kubeapp",
		Short: "CLI for interacting with Kubernetes clusters",
		Long: `kubeapp is a CLI tool for inspecting Kubernetes clusters.
It reads your kubeconfig (~/.kube/config by default) and works with
both local clusters (e.g. minikube) and remote clusters.`,
		// Initialise the k8s client before any subcommand runs.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			client, err = k8sclient.NewClient(kubeconfig)
			if err != nil {
				return fmt.Errorf("could not connect to cluster: %w", err)
			}
			return nil
		},
	}

	defaultKubeconfig := ""
	if home := homedir.HomeDir(); home != "" {
		defaultKubeconfig = filepath.Join(home, ".kube", "config")
	}

	root.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", defaultKubeconfig, "path to kubeconfig file")
	root.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "target namespace")

	root.AddCommand(newPodsCmd())
	root.AddCommand(newServicesCmd())

	return root
}
