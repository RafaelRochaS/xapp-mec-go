# xApp MEC — Computional offloading platform for O-RAN environments 

NearRT RIC xApp to demonstrate MEC services in an O-RAN network, based on https://github.com/o-ran-sc/ric-app-hw-go.

## Usage
The xApp cannot be run by itself, instead relying on a NearRT RIC platform to reside. Afterward, the xApp can be onboarded to the platform,
and the offload endpoints will be available.

### Deploying The Platform
The O-RAN SC [RIC](https://github.com/o-ran-sc/ric-plt-ric-dep) is recommended, with a simple deployment
being sufficient to run the xApp. However, it is important to note that the [xApp Onboarder](https://github.com/o-ran-sc/ric-plt-ric-dep/tree/master/new-installer/helm/charts/nearrtric/xapp-onboarder) must be part of the RIC platform.

At this moment, the recommended way to deploy the xApp is to use the [new installer](https://github.com/o-ran-sc/ric-plt-ric-dep/tree/master/new-installer) of the RIC platform.

## Onboarding The xApp
To actually onboard the xApp, it is necessary to generate a [Helm](https://helm.sh/) chart, with the correct parameters and
pointing to the correct image. The built image on [Docker hub](https://hub.docker.com/repository/docker/rafaelrs94/xapp-mec/tags/mec-xapp/) can be used, 
or it can be built from source as well

**Note:** If built from source, the image must be pushed to a Docker registry accessible by the Kubernetes cluster.