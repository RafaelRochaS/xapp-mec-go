package handlers

import (
	"fmt"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/models"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/utils"
	"golang.org/x/net/context"
)

func HandleOffload(edgeClient, cloudClient *utils.K8sClient, task models.Task) error {
	xapp.Logger.Info("Handling offload request")

	edgeMetrics, err := utils.GetNodesResources(edgeClient.Metrics, context.TODO())

	if err != nil {
		return err
	}

	cpuOffThreshold, memOffThreshold := parseThresholds()

	for _, edgeMetric := range *edgeMetrics {
		xapp.Logger.Debug("HandleOffload :: edgeMetric: ", edgeMetric)

		cpuValue := edgeMetric.Cpu().AsApproximateFloat64()
		xapp.Logger.Debug("HandleOffload :: CPU value: ", cpuValue)

		memValue := edgeMetric.Memory().AsApproximateFloat64()
		xapp.Logger.Debug("HandleOffload :: Memory value: ", memValue)

		if int(cpuValue) < cpuOffThreshold && int(memValue) < memOffThreshold {
			xapp.Logger.Info("HandleOffload :: edge server resources within threshold, offloading resource")
			ctx := context.WithValue(context.Background(), "executionSite", utils.EdgeExecutionSite)

			return utils.OffloadTask(edgeClient.ClientSet, task, ctx)
		}
	}

	xapp.Logger.Info("HandleOffload :: edge server resources above threshold, offloading to cloud")

	ctx := context.WithValue(context.Background(), "executionSite", utils.CloudExecutionSite)
	return utils.OffloadTask(cloudClient.ClientSet, task, ctx)
}

func parseThresholds() (cpuThreshold, memThreshold int) {
	cpuThreshold = xapp.Config.GetInt("offload.threshold.cpu")
	memThreshold = xapp.Config.GetInt("offload.threshold.mem")

	if cpuThreshold <= 0 {
		cpuThreshold = 3
	}

	if memThreshold <= 0 {
		memThreshold = 35
	}

	xapp.Logger.Debug(fmt.Sprintf("HandleOffload :: threshold: CPU: %d\tMem: %d", cpuThreshold, memThreshold))

	return
}
