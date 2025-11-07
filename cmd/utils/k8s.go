package utils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClientType int

const (
	InCluster K8sClientType = iota
	OutCluster
)

func GetK8sClient(clientType K8sClientType) (error, *kubernetes.Clientset) {
	var config *rest.Config
	var err error

	if clientType == InCluster {
		config, err = getInClusterConfig()
	} else {
		config, err = getOutClusterConfig()
	}

	if err != nil {
		return err, nil
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err, nil
	}

	return nil, clientSet
}

func getInClusterConfig() (*rest.Config, error) {
	return rest.InClusterConfig()
}

func getOutClusterConfig() (*rest.Config, error) {
	return clientcmd.BuildConfigFromFlags("", "./config/kubeconfig")
}
