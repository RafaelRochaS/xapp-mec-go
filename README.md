# xApp MEC — Computional offloading platform for O-RAN environments 

NearRT RIC xApp to demonstrate MEC services in an O-RAN network, based on https://github.com/o-ran-sc/ric-app-hw-go.

## Usage
The xApp cannot be run by itself, instead relying on a NearRT RIC platform to reside. Afterward, the xApp can be onboarded to the platform,
and the offload endpoints will be available.

### Deploying The Platform
The O-RAN SC [RIC](https://github.com/o-ran-sc/ric-plt-ric-dep) is recommended, with a simple deployment
being sufficient to run the xApp. However, it is important to note that the [xApp Onboarder](https://github.com/o-ran-sc/ric-plt-ric-dep/tree/master/new-installer/helm/charts/nearrtric/xapp-onboarder) must be part of the RIC platform.

At this moment, the recommended way to deploy the xApp is to use the [new installer](https://github.com/o-ran-sc/ric-plt-ric-dep/tree/master/new-installer) of the RIC platform.

### Onboarding The xApp
To actually onboard the xApp, it is necessary to generate a [Helm](https://helm.sh/) chart, with the correct parameters and
pointing to the correct image. The built image on [Docker hub](https://hub.docker.com/repository/docker/rafaelrs94/xapp-mec/tags/mec-xapp/) can be used, 
or it can be built from source as well

**Note:** If built from source, the image must be pushed to a Docker registry accessible by the Kubernetes cluster.

### Offload Scenarios
If there are resources available, the xApp attempts to offload the tasks to the local Kubernetes cluster. Otherwise, it will offload to the remote Kubernetes cluster.
As such, for the cloud offload to work properly, a kubeconfig file must be injected into the xApp pod. The kubeconfig file can be placed into
the ```config/``` directory, and will automatically be picked up by the xApp.

At present, only CPU load is taken into consideration for offloading. The parameter can be set on the ```config/config-file.json```, on ```offload.threshold``` field.

The tasks are offloaded to the ```task-offload``` namespace, which must be created manually, both in the local and remote clusters.
Additionally, the xApp requires specific permissions to be able to read the metrics of the local cluster and to create jobs in both clusters.
For convenience, the provided cluster roles and bindings in the ```config\``` directory can be used.

### Tasks
Tasks are represented by the ```Task``` struct, which contains the following fields:
```
	DeviceId     int    `json:"deviceId"`
	Task         string `json:"task"`
	Image        string `json:"image"`
	CPU          string `json:"cpu"`
	Mem          int    `json:"mem"`
	DeadlineSecs int    `json:"deadlineSecs,omitempty"`
	Workload     int    `json:"workload"`
	CallbackUrl  string `json:"callbackUrl,omitempty"`
```

A task is referenced via its UUID, which is generated when registering the task. Tasks are stored in the [SDL storage](https://docs.o-ran-sc.org/projects/o-ran-sc-ric-plt-sdl/en/latest/), 
and can be retrieved via the ```RetrieveTask``` endpoint.

To start a task, the DeviceId used to register the task, as well as its UUID must be provided:

```
	Id       string `json:"id"`
	DeviceId int    `json:"deviceId"`
```

### Endpoints
The xApp exposes endpoints to register, retrieve and delete tasks, as well as start previously registered tasks.
The routes are defined as follows:

```
POST /ric/v1/mec/start
POST /ric/v1/mec/tasks
GET|DELETE /ric/v1/mec/tasks/{taskId}
```