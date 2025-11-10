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
	sdlClient   *xapp.SDLStorage
	edgeClient  *utils.K8sClient
	cloudClient *utils.K8sClient
}

func NewTaskHandler() *TaskHandler {
	inClusterClient, err := utils.GetK8sClient(utils.InCluster)

	if err != nil {
		xapp.Logger.Error("Error in getting in cluster k8s client: %v", err)
		panic(err)
	}

	outClusterClient, err := utils.GetK8sClient(utils.OutCluster)

	if err != nil {
		xapp.Logger.Error("Error in getting out of cluster k8s client: %v", err)
		panic(err)
	}
	return &TaskHandler{
		sdlClient:   xapp.NewSdlStorage(),
		edgeClient:  inClusterClient,
		cloudClient: outClusterClient}
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

	serializedTask, err := json.Marshal(task)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to serialize task:"))
		_, _ = w.Write([]byte(err.Error()))
	}

	xapp.Logger.Debug("Serialized task: %s", string(serializedTask))

	err = t.sdlClient.Store(utils.TaskNamespace, jobId, serializedTask)

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

	task := value[jobId].([]byte)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to retrieve task:"))
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("Retrieved task:"))
	_, _ = w.Write(task)
}

func (t *TaskHandler) StartTask(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("StartTaskHandler handler")
	xapp.Logger.Debug("Request body: %s", r.Body)

	var taskRequest models.StartTaskRequest
	err := json.NewDecoder(r.Body).Decode(&taskRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Failed to parse request body:"))
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	value, err := t.sdlClient.Read(utils.TaskNamespace, taskRequest.Id)

	jsonTask := value[taskRequest.Id].([]byte)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to retrieve task:"))
		_, _ = w.Write([]byte(err.Error()))
	}

	var task models.Task

	if err = json.Unmarshal(jsonTask, &task); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to parse task:"))
		_, _ = w.Write([]byte(err.Error()))
	}

	err = HandleOffload(t.edgeClient, t.cloudClient, task)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to start task:"))
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Task started"))
}

func (t *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("DeleteTaskHandler handler")
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

	err = t.sdlClient.Delete(utils.TaskNamespace, []string{parsedJobId.String()})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to delete task:"))
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
