package main

import (
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/handlers"
)

func RegisterRoutes() {
	t := handlers.NewTaskHandler()

	xapp.Resource.InjectRoute("/ric/v1/mec/start", t.StartTask, "POST")
	xapp.Resource.InjectRoute("/ric/v1/mec/tasks", t.RegisterTask, "POST")
	xapp.Resource.InjectRoute("/ric/v1/mec/tasks/{taskId}", t.RetrieveTask, "GET")
	xapp.Resource.InjectRoute("/ric/v1/mec/tasks/{taskId}", t.DeleteTask, "DELETE")
}
