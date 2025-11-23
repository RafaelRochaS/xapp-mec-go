package handlers

import (
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/models"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/utils"
	"golang.org/x/net/context"
)

func HandleOffload(edgeClient, cloudClient *utils.K8sClient, task models.Task) error {
	xapp.Logger.Info("Handling offload request")

	offloadThreshold := xapp.Config.GetInt("offload.threshold")
	xapp.Logger.Debug("HandleOffload :: threshold: ", offloadThreshold)

	edgeMetrics, err := utils.GetNodesResources(edgeClient.Metrics, context.TODO())

	if err != nil {
		return err
	}

	for _, edgeMetric := range *edgeMetrics {
		xapp.Logger.Debug("HandleOffload :: edgeMetric: ", edgeMetric)
		cpuValue := edgeMetric.Cpu().AsApproximateFloat64()
		cpuAsValue := edgeMetric.Cpu().Value()
		xapp.Logger.Debug("HandleOffload :: CPU value: ", cpuValue)
		xapp.Logger.Debug("HandleOffload :: CPU As Value: ", cpuAsValue)

		if int(cpuAsValue) < offloadThreshold {
			xapp.Logger.Info("HandleOffload :: edge server resources within threshold, offloading resource")
			ctx := context.WithValue(context.Background(), "executionSite", utils.EdgeExecutionSite)

			return utils.OffloadTask(edgeClient.ClientSet, task, ctx)
		}
	}

	xapp.Logger.Info("HandleOffload :: edge server resources above threshold, offloading to cloud")

	ctx := context.WithValue(context.Background(), "executionSite", utils.CloudExecutionSite)
	return utils.OffloadTask(cloudClient.ClientSet, task, ctx)
}
