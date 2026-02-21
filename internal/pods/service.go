package pods

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Service handles pod-related operations against the Kubernetes API.
// It is intentionally decoupled from presentation so it can be reused
// by both the CLI layer and a future HTTP API layer.
type Service struct {
	client *kubernetes.Clientset
}

func NewService(client *kubernetes.Clientset) *Service {
	return &Service{client: client}
}

// List returns all pods in the given namespace. Pass an empty string to list
// across all namespaces.
func (s *Service) List(ctx context.Context, namespace string) ([]corev1.Pod, error) {
	pods, err := s.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

// Get returns a single pod by name within the given namespace.
func (s *Service) Get(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	return s.client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
}
