package handlers

import (
	"encoding/json"
	"net/http"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/models"
	"github.com/RafaelRochaS/xapp-mec-go/cmd/utils"
	"github.com/google/uuid"
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

	jobId := uuid.New().String()
	task := models.Task{
		Id:              jobId,
		RegisterRequest: taskRequest,
	}

	err = t.sdlClient.Store(utils.TaskNamespace, jobId, task)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to register task:"))
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("Registered job:"))
	_, _ = w.Write([]byte(jobId))

	return
}

func (t *TaskHandler) RetrieveTask(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("RetrieveTaskHandler handler")
	xapp.Logger.Debug("Request body: %s", r.Body)

	jobId := r.URL.Query().Get("jobId")

	if jobId == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Missing jobId parameter"))
		return
	}

	parsedJobId, err := uuid.Parse(jobId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid jobId parameter:"))
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	value, err := t.sdlClient.Read(utils.TaskNamespace, parsedJobId.String())

	parsedTask, err := json.Marshal(value)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to retrieve task:"))
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("Retrieved task:"))
	_, _ = w.Write(parsedTask)
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

	jobId := r.URL.Query().Get("jobId")

	if jobId == "" {
		w.WriteHeader(http.StatusNoContent)
		_, _ = w.Write([]byte("Missing jobId parameter"))
		return
	}

	parsedJobId, err := uuid.Parse(jobId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid jobId parameter:"))
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	err = t.sdlClient.Delete(utils.TaskNamespace, []string{parsedJobId.String()})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to delete task:"))
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
