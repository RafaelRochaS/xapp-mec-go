package utils

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/models"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type K8sClient struct {
	ClientSet *kubernetes.Clientset
	Metrics   *metrics.Clientset
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

	return &K8sClient{ClientSet: clientSet, Metrics: metricsClient}, err
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

func GetNodesResources(m *metrics.Clientset, ctx context.Context) (*[]v1.ResourceList, error) {
	nodeList, err := m.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	resources := make([]v1.ResourceList, len(nodeList.Items))

	for i, node := range nodeList.Items {
		log.Printf("Node resource: %+v", node.Usage)
		resources[i] = node.Usage
	}

	log.Printf("Got resources: %+v", resources)

	return &resources, nil
}

func OffloadTask(c *kubernetes.Clientset, task models.Task, ctx context.Context) error {
	var ttlSecondsAfterFinished int32 = 5

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: task.Id,
			Labels: map[string]string{
				"offload":  "true",
				"deviceId": strconv.Itoa(task.DeviceId),
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  fmt.Sprintf("task-container-%s", task.Id),
					Image: task.Image,
					Env: []v1.EnvVar{
						{
							Name:  "WORKLOAD_SIZE",
							Value: strconv.Itoa(task.Workload),
						},
						{
							Name:  "DEVICE_ID",
							Value: strconv.Itoa(task.DeviceId),
						},
						{
							Name:  "EXECUTION_SITE",
							Value: ctx.Value("executionSite").(string),
						},
						{
							Name:  "TASK_ID",
							Value: task.Id,
						},
						{
							Name:  "CALLBACK_ADDR",
							Value: task.CallbackUrl,
						},
					},
				},
			},
			RestartPolicy: v1.RestartPolicyOnFailure,
		},
	}

	jobsClient := c.BatchV1().Jobs("task-offload")
	job := &batchv1.Job{
		ObjectMeta: pod.ObjectMeta,
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: v1.PodTemplateSpec{
				Spec: pod.Spec,
			},
		},
	}

	log.Println("Offloading task: ", task.Id)

	_, err := jobsClient.Create(context.TODO(), job, metav1.CreateOptions{})

	return err
}
