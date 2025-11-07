package utils

import (
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type K8sClient struct {
	clientSet *kubernetes.Clientset
	metrics   *metrics.Clientset
}

type K8sClientType int

const (
	InCluster K8sClientType = iota
	OutCluster
)

func (t K8sClientType) String() string {
	return [...]string{"InCluster", "OutCluster"}[t]
}

func GetK8sClient(clientType K8sClientType) (*K8sClient, error) {
	xapp.Logger.Info("Getting k8s client for type: ", clientType)
	var config *rest.Config
	var err error

	if clientType == InCluster {
		config, err = getInClusterConfig()
	} else {
		config, err = getOutClusterConfig()
	}

	xapp.Logger.Info("Got config:", config)

	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	metricsClient, err := getMetricsClient(config)

	xapp.Logger.Info("Got k8s client set:", clientSet)

	return &K8sClient{clientSet: clientSet, metrics: metricsClient}, err
}

func getInClusterConfig() (*rest.Config, error) {
	return rest.InClusterConfig()
}

func getOutClusterConfig() (*rest.Config, error) {
	return clientcmd.BuildConfigFromFlags("", "./config/kubeconfig")
}

func getMetricsClient(config *rest.Config) (*metrics.Clientset, error) {
	return metrics.NewForConfig(config)
}
