package handlers

import (
	"net/http"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
)

func RetrieveTaskHandler(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("RetrieveTaskHandler handler")
	xapp.Logger.Debug("Request body: %s", r.Body)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("[WIP] Task status: 0"))
}
