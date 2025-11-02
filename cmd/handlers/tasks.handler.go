package handlers

import (
	"encoding/json"
	"net/http"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/models"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/utils"
)

type TaskHandler struct {
	sdlClient *xapp.SDLStorage
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{sdlClient: xapp.NewSdlStorage()}
}

func (t *TaskHandler) RegisterTask(w http.ResponseWriter, r *http.Request) {
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

	err = t.sdlClient.Store(utils.TaskNamespace, "deviceId", taskRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Failed to register task:"))
		_, _ = w.Write([]byte(err.Error()))

		return
	}
}

func (t *TaskHandler) RetrieveTask(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("RetrieveTaskHandler handler")
	xapp.Logger.Debug("Request body: %s", r.Body)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("[WIP] Task status: 0"))
}

func (t *TaskHandler) RetrieveAllTasks(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("RetrieveAllTasks handler")
	xapp.Logger.Debug("Request body: %s", r.Body)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("[WIP] Task status: 0"))
}

func (t *TaskHandler) StartTask(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("StartTaskHandler handler")
	xapp.Logger.Debug("Request body: %s", r.Body)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("[WIP] Task started"))
}

func (t *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("DeleteTaskHandler handler")
	xapp.Logger.Debug("Request body: %s", r.Body)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("[WIP] Task deleted"))
}
