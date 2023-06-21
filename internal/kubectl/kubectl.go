package kubectl

import (
	"errors"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Kubectl struct {
	client *kubernetes.Clientset
}

func New() (*Kubectl, error) {
	if client, err := NewFromInCluster(); err == nil {
		return client, nil
	}

	if client, err := NewFromOutOfCluster(); err == nil {
		return client, nil
	}

	return nil, errors.New("unable to initialize kubectl")
}

func NewFromInCluster() (*Kubectl, error) {
	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	return &Kubectl{
		client: clientset,
	}, nil
}

func NewFromOutOfCluster() (*Kubectl, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	return &Kubectl{
		client: clientset,
	}, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
