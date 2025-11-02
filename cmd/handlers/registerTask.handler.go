package handlers

import (
	"encoding/json"
	"net/http"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/models"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/utils"
)

func RegisterTaskHandler(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("Register task handler")
	xapp.Logger.Debug("Request body: %s", r.Body)

	var taskRequest models.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&taskRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Failed to parse request body:"))
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	sdlTaskRegister := xapp.NewSdlStorage()
	err = sdlTaskRegister.Store(utils.TaskNamespace, "deviceId", taskRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Failed to register task:"))
		_, _ = w.Write([]byte(err.Error()))

		return
	}
}
