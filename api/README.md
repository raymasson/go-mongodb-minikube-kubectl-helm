# Golang API Deployment for Kubernetes on Minikube using helm

An example project demonstrating the deployment of a Golang API via Kubernetes on Minikube (Kubernetes running locally on a workstation). Helm is used for managing the deployment configuration.

## 1 - Prerequisites

Ensure the following dependencies are already fulfilled on your host Linux/Windows/Mac Workstation/Laptop:

1.  The [VirtualBox](https://www.virtualbox.org/wiki/Downloads) hypervisor has been installed.
2.  [Docker](https://docs.docker.com/install/) has been installed.
3.  The [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) command-line tool for Kubernetes has been installed.
4.  The [Minikube](https://github.com/kubernetes/minikube/releases) tool for running Kubernetes locally has been installed.
5.  The Minikube cluster has been started, inside a local Virtual Machine, using the following command (also includes commands to check that kubectl is configured correctly to see the running minikube pod):

    ```
    $ minikube start
    $ kubectl get nodes
    $ kubectl describe nodes
    $ kubectl get services
    ```

6.  [Helm/Tiller](https://docs.helm.sh/using_helm/) have been installed and initialized

## 2 - Docker Image

#### 2.1 Build the docker image

Kubernetes needs a docker image so that the api can be deployed.
I use dockerhub to store my docker image. Any other docker registry can be used though.
Let's build the image with the following command:

```
$ make USERNAME="<USERNAME>" build
```

`<USERNAME>` must be replaced by your Docker ID. (e.g. : raymasson)

Check that the image has been created in your local repository by executing the command below:

```
$ docker images | grep person-api
```

#### 2.2 Publish the docker image

Now that you've created a Docker image which contains your api code, you can upload it into a public registry. This project uses Docker Hub, but you can select one of your own, such as:

* [Google Container Registry](https://cloud.google.com/container-registry/)
* [Amazon EC2 Container Registry](https://aws.amazon.com/ecr/)
* [Azure container Registry](https://azure.microsoft.com/en-us/services/container-registry/)
* [Quay](https://quay.io/)

To upload the image to Docker Hub, follow the steps below:

* Log in to Docker Hub:

```
$ docker login
```

* Push the image to your Docker Hub account. Replace the USERNAME placeholder with your Docker ID:

```
$ docker push USERNAME/person-api:1.0.0
```

## 3 - Deployment

#### 3.1 Edit the values.yaml file

Edit the values.yaml file (helm-chart folder) to replace the current image repository name with the name of the docker hub repository.

#### 3.2 Deploy the api in Kubernetes

Make sure that you can to connect to your Kubernetes cluster by executing the command below:

```
$ kubectl cluster-info
```

Deploy the Helm chart by executing the following (`person` is an example name, replace it with the name you want to give to your api):

```
$ helm init
$ helm upgrade --install person helm-chart/
```

Check the status and logs of the deployed pod thanks to the kubernetes dashboard:

```
$ minikube dashboard
```
