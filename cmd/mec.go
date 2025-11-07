package main

import (
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/utils"
)

type MECApp struct {
	stats               map[string]xapp.Counter
	inClusterClientSet  *utils.K8sClient
	outClusterClientSet *utils.K8sClient
}

func main() {
	metrics := []xapp.CounterOpts{
		{
			Name: "RICIndicationRx",
			Help: "Total number of RIC Indication message received",
		},
	}

	inClusterClient, err := utils.GetK8sClient(utils.InCluster)

	if err != nil {
		xapp.Logger.Error("Error in getting in cluster k8s client: %v", err)
		return
	}

	outClusterClient, err := utils.GetK8sClient(utils.OutCluster)

	if err != nil {
		xapp.Logger.Error("Error in getting out of cluster k8s client: %v", err)
		return
	}

	mec := MECApp{
		stats:               xapp.Metric.RegisterCounterGroup(metrics, "mec_app"),
		inClusterClientSet:  inClusterClient,
		outClusterClientSet: outClusterClient,
	}

	RegisterRoutes()

	mec.Run()
}

func (e *MECApp) Run() {

	xapp.Logger.SetMdc("MECApp", "0.0.4")
	xapp.AddConfigChangeListener(e.ConfigChangeHandler)
	xapp.SetReadyCB(e.xAppStartCB, true)
	waitForSdl := xapp.Config.GetBool("db.waitForSdl")

	xapp.RunWithParams(e, waitForSdl)
}

func (e *MECApp) ConfigChangeHandler(f string) {
	xapp.Logger.Info("Config file changed", f)
}

func (e *MECApp) xAppStartCB(d interface{}) {
	xapp.Logger.Info("xApp ready call back received", d)
}

func (e *MECApp) handleRICIndication(ranName string, r *xapp.RMRParams) {
	xapp.Logger.Info("handleRICIndication", ranName, "\tparams: ", *r)
	e.stats["RICIndicationRx"].Inc()
}

func (e *MECApp) Consume(msg *xapp.RMRParams) (err error) {
	id := xapp.Rmr.GetRicMessageName(msg.Mtype)

	xapp.Logger.Info("Message received: name=%s meid=%s subId=%d txid=%s len=%d", id, msg.Meid.RanName, msg.SubId, msg.Xid, msg.PayloadLen)

	switch id {
	case "RIC_HEALTH_CHECK_REQ":
		xapp.Logger.Info("Received health check request")

	case "RIC_INDICATION":
		xapp.Logger.Info("Received RIC Indication message")
		e.handleRICIndication(msg.Meid.RanName, msg)

	default:
		xapp.Logger.Info("Unknown message type '%d', discarding", msg.Mtype)
	}

	defer func() {
		xapp.Rmr.Free(msg.Mbuf)
		msg.Mbuf = nil
	}()
	return
}
