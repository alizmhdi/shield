package k8s

import (
	"context"
	"os"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Client wraps the Kubernetes clientset for shield operations.
type Client struct {
	Clientset *kubernetes.Clientset
}

// NewClient creates a new Kubernetes client using the local kubeconfig.
func NewClient() (*Client, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = os.ExpandEnv("$HOME/.kube/config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Client{Clientset: clientset}, nil
}

// ListIngresses lists all Ingress resources in the given namespace.
func (c *Client) ListIngresses(ctx context.Context, namespace string) ([]networkingv1.Ingress, error) {
	ingList, err := c.Clientset.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ingList.Items, nil
}

// UpdateIngress updates the given Ingress resource in the specified namespace.
func (c *Client) UpdateIngress(ctx context.Context, namespace string, ingress *networkingv1.Ingress) error {
	_, err := c.Clientset.NetworkingV1().Ingresses(namespace).Update(ctx, ingress, metav1.UpdateOptions{})
	return err
}
