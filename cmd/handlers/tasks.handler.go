package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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

	taskId := uuid.New().String()
	task := models.Task{
		Id:              taskId,
		RegisterRequest: taskRequest,
	}

	serializedTask, err := json.Marshal(task)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to serialize task:"))
		_, _ = w.Write([]byte(err.Error()))
	}

	xapp.Logger.Debug("Serialized task: %s", string(serializedTask))

	err = t.sdlClient.Store(utils.TaskNamespace, taskId, string(serializedTask))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to register task:"))
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	output := models.RegisterTaskResponse{
		Id: taskId,
	}

	parsed, err := json.Marshal(output)

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(parsed)

	return
}

func (t *TaskHandler) RetrieveTask(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("RetrieveTaskHandler handler")
	xapp.Logger.Debug("Request body: %s", r.Body)

	path := strings.Split(r.URL.Path, "/")
	taskId := path[len(path)-1]

	if taskId == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Missing taskId parameter"))
		return
	}

	task, err := t.retrieveTask(taskId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Failed to retrieve task:"))
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	marshalledTask, err := json.Marshal(task)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to retrieve task:"))
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("Retrieved task:"))
	_, _ = w.Write(marshalledTask)
}

func (t *TaskHandler) StartTask(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("StartTaskHandler handler")
	xapp.Logger.Debug("Request body: ", r.Body)

	var taskRequest models.StartTaskRequest
	err := json.NewDecoder(r.Body).Decode(&taskRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Failed to parse request body:"))
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	parsedTask, err := t.retrieveTask(taskRequest.Id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to parse task: "))
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	err = HandleOffload(t.edgeClient, t.cloudClient, parsedTask)

	if err != nil {
		xapp.Logger.Error("Failed to offload task: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to start task:"))
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Task started"))

	err = t.sdlClient.Delete(utils.TaskNamespace, []string{parsedTask.Id})

	xapp.Logger.Debug("Deleted task from sdl storage: %v", err)
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

func (t *TaskHandler) retrieveTask(id string) (models.Task, error) {
	xapp.Logger.Debug("retrieveTask handler :: id: ", id)
	parsedId, err := uuid.Parse(id)

	if err != nil {
		return models.Task{}, err
	}

	xapp.Logger.Debug("retrieveTask handler :: parsedId: ", id)

	value, err := t.sdlClient.Read(utils.TaskNamespace, parsedId.String())

	if err != nil {
		return models.Task{}, err
	}

	xapp.Logger.Debug("Read value: %+v", value)

	task, ok := value[id].(string)
	xapp.Logger.Debug("Got task: %+v", task)

	if !ok {
		return models.Task{}, errors.New("task is not string")
	}

	var parsedTask models.Task
	err = json.Unmarshal([]byte(task), &parsedTask)

	if err != nil {
		return models.Task{}, err
	}

	return parsedTask, nil
}
