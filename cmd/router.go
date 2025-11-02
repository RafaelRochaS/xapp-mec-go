package main

import (
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/handlers"
)

func RegisterRoutes() {
	xapp.Resource.InjectRoute("/ric/v1/mec/start", handlers.StartTaskHandler, "POST")
	xapp.Resource.InjectRoute("/ric/v1/mec/register", handlers.RegisterTaskHandler, "POST")
	xapp.Resource.InjectRoute("/ric/v1/mec/jobs/{jobId}", handlers.RetrieveTaskHandler, "GET")
}
