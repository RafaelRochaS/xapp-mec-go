package handlers

import (
	"net/http"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
)

func StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("StartTaskHandler handler")
	xapp.Logger.Debug("Request body: %s", r.Body)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("[WIP] Task started"))
}
