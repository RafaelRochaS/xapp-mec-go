package handlers

import (
	"net/http"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
)

func RetrieveAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("RetrieveAllTasksHandler handler")
	xapp.Logger.Debug("Request body: %s", r.Body)
}
