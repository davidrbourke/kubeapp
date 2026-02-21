package services

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Service handles Kubernetes Service resource operations.
// It is intentionally decoupled from presentation so it can be reused
// by both the CLI layer and a future HTTP API layer.
type Service struct {
	client *kubernetes.Clientset
}

func NewService(client *kubernetes.Clientset) *Service {
	return &Service{client: client}
}

// List returns all Kubernetes services in the given namespace. Pass an empty
// string to list across all namespaces.
func (s *Service) List(ctx context.Context, namespace string) ([]corev1.Service, error) {
	svcs, err := s.client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return svcs.Items, nil
}

// Get returns a single Kubernetes service by name within the given namespace.
func (s *Service) Get(ctx context.Context, namespace, name string) (*corev1.Service, error) {
	return s.client.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
}
