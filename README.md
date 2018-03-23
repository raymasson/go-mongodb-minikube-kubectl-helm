# Deployment of a persisted mongodb replica set and a golang api in kubernetes

An example project demonstrating the deployment of a mongodb replica set and a Golang API in Kubernetes.
Minikube is used to run k8s locally.
kubectl is used to deploy the mongodb replica set.
helm is used to deploy the golang API.

## 1 - Deploy the mongodb replica set

Follow the mongodb [readme](https://github.com/raymasson/go-mongodb-minikube-kubectl-helm/tree/master/mongodb)

## 2 - Deploy the golang API

Follow the api [readme](https://github.com/raymasson/go-mongodb-minikube-kubectl-helm/tree/master/api)

## 3 - Test the API

First of all, let's get the API service URL out of minikube.

```
$ minikube service  person-api --url
```

Use the outputted URL in your internet browser (e.g.: `http://192.168.99.100:31518`). You should see `Welcome to the person API!` displayed

Get persons out of the person API. The expected output is a `HTTP status 200 OK` and a `[]` of persons.

```
$ curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://192.168.99.100:31518/persons
```

Now, let's insert our first person. The expected output is a `HTTP status 201 Created`.

```
$ curl -i -H "Accept: application/json" -H "Content-Type: application/json" -d '{"first_name":"ray", "last_name":"masson"}'  -X POST http://192.168.99.100:31518/persons
```

Run again the first `curl` command. The output should now be something like `[{"id":"5ab505c424e8f10001f047b5","first_name":"ray","last_name":"masson"}]`
