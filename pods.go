package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func getDefaultKubeconfigPath() string {
	if os.Getenv("KUBECONFIG") != "" {
		return os.Getenv("KUBECONFIG")
	}
	if home := homedir.HomeDir(); home != "" {
		path := filepath.Join(home, ".kube", "config")
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func getClientset(kubeconfig, context string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if kubeconfig == "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("error creating in-cluster config: %v", err)
		}
	} else {
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
			&clientcmd.ConfigOverrides{
				CurrentContext: context,
			}).ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("error building config: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error building clientset: %w", err)
	}

	return clientset, nil
}

func getImagesFromAllPods() ([]string, error) {
	clientset, err := getClientset(getDefaultKubeconfigPath(), "")
	if err != nil {
		return nil, fmt.Errorf("error getting clientset: %v", err)
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing pods: %v", err)
	}

	var images []string
	imageSet := make(map[string]struct{})
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if _, exists := imageSet[container.Image]; !exists {
				images = append(images, container.Image)
				imageSet[container.Image] = struct{}{}
			}
		}
	}

	return images, nil
}
