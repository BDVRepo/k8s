package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/olekukonko/tablewriter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Resource represents a Kubernetes resource with its details.
type Resource struct {
	Namespace    string
	Name         string
	ResourceType string
	CreatedAt    time.Time
}

// getClientset creates a Kubernetes clientset from the kubeconfig file.
func getClientset() (*kubernetes.Clientset, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

// getResources fetches the resources with external access.
func getResources(clientset *kubernetes.Clientset) ([]Resource, error) {
	var resources []Resource

	// Get services with type NodePort and LoadBalancer
	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, svc := range services.Items {
		if svc.Spec.Type == "NodePort" || svc.Spec.Type == "LoadBalancer" {
			resources = append(resources, Resource{
				Namespace:    svc.Namespace,
				Name:         "svc/" + svc.Name,
				ResourceType: string(svc.Spec.Type),
				CreatedAt:    svc.CreationTimestamp.Time,
			})
		}
	}

	// Get ingresses
	ingresses, err := clientset.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, ing := range ingresses.Items {
		resources = append(resources, Resource{
			Namespace:    ing.Namespace,
			Name:         "ing/" + ing.Name,
			ResourceType: "Ingress",
			CreatedAt:    ing.CreationTimestamp.Time,
		})
	}

	// Get pods with host network
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, pod := range pods.Items {
		if pod.Spec.HostNetwork {
			resources = append(resources, Resource{
				Namespace:    pod.Namespace,
				Name:         "pod/" + pod.Name,
				ResourceType: "HostNetwork",
				CreatedAt:    pod.CreationTimestamp.Time,
			})
		}
	}

	return resources, nil
}

func main() {
	clientset, err := getClientset()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	resources, err := getResources(clientset)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching resources: %v\n", err)
		os.Exit(1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Namespace", "Resource/Name", "Resource Type", "Created At"})

	for _, resource := range resources {
		table.Append([]string{resource.Namespace, resource.Name, resource.ResourceType, resource.CreatedAt.Format(time.RFC3339)})
	}

	table.Render()
}
