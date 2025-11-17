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

	for _, edgeMetric := range edgeMetrics {
		xapp.Logger.Info("Edge metrics: ", edgeMetric)
		val, ok := edgeMetric.Cpu().AsDec().Unscaled()
		xapp.Logger.Info("CPU value: ", val)

		if ok && int(val) < offloadThreshold {
			return utils.OffloadTask(edgeClient.ClientSet, task)
		}
	}

	return utils.OffloadTask(cloudClient.ClientSet, task)
}
