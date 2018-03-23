# MongoDB Deployment for Kubernetes on Minikube (i.e. running on local workstation)

An example project demonstrating the deployment of a MongoDB Replica Set via Kubernetes on Minikube (Kubernetes running locally on a workstation). Contains example Kubernetes YAML deployment file (in the 'deployment' folder), associated Kubernetes based Bash scripts (in the 'scripts' folder) to configure the environment and deploy a MongoDB Replica Set and a self-signed SSL certificate (in the '
ssl' folder) for SSL support on the deployed mongodb cluster.

## 1 - Prerequisites

Ensure the following dependencies are already fulfilled on your host Linux/Windows/Mac Workstation/Laptop:

1.  The [VirtualBox](https://www.virtualbox.org/wiki/Downloads) hypervisor has been installed.
2.  The [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) command-line tool for Kubernetes has been installed.
3.  The [Minikube](https://github.com/kubernetes/minikube/releases) tool for running Kubernetes locally has been installed.
4.  The Minikube cluster has been started, inside a local Virtual Machine, using the following command (also includes commands to check that kubectl is configured correctly to see the running minikube pod):

    ```
    $ minikube start
    $ kubectl get nodes
    $ kubectl describe nodes
    $ kubectl get services
    ```

## 2 - Main Deployment Steps

1.  To deploy the MongoDB Service (including the StatefulSet running "mongod" containers), via a command-line terminal/shell, execute the following:

    ```
    $ cd scripts
    $ ./generate.sh
    ```

2.  Re-run the following command, until all 3 “mongod” pods (and their containers) have been successfully started (“Status=Running”; usually takes a minute or two).

    ```
    $ kubectl get all
    ```

3.  Execute the following script which connects to the first Mongod instance running in a container of the Kubernetes StatefulSet, via the Mongo Shell, to (1) initialise the MongoDB Replica Set, and (2) create a MongoDB admin user (specify the password you want as the argument to the script, replacing 'abc123').

    ```
    $ ./configure_repset_auth.sh abc123
    ```

    If you have an error like `Error: couldn't add user: not master`, execute the following command again on the PRIMARY mongodb replica (e.g.: mongod-1 is PRIMARY):

    ```
    $ kubectl exec mongod-1 -c mongod-container -- mongo --eval 'db.getSiblingDB("admin").createUser({user:"main_admin",pwd:"abc123",roles:[{role:"root",db:"admin"}]});'
    ```

You should now have a MongoDB Replica Set initialised, secured and running in a Kubernetes StatefulSet.

You can also view the state of the deployed environment, via the Kubernetes dashboard, which can be launched in a browser with the following command: `$ minikube dashboard`

## 3 - Example Tests To Run To Check Things Are Working

Use this section to prove:

1.  Data is being replicated between members of the containerised replica set.
2.  Data is retained even when the MongoDB Service/StatefulSet is removed and then re-created (by virtue of re-using the same Persistent Volume Claims).

#### 3.1 Replication Test

Connect to the container running the first "mongod" replica, then use the Mongo Shell to authenticate and add some test data to a database (e.g.: mongod-0 is the `PRIMARY`):

    $ kubectl exec -it mongod-0 -c mongod-container bash
    $ mongo
    > db.getSiblingDB('admin').auth("main_admin", "abc123");
    > use test;
    > db.testcoll.insert({a:1});
    > db.testcoll.insert({b:2});
    > db.testcoll.find();

Exit out of the shell and exit out of the first container (“mongod-0”). Then connect to the second container (“mongod-1”), run the Mongo Shell again and see if the previously inserted data is visible to the second "mongod" replica (e.g.: mongod-1 is `SECONDARY`):

    $ kubectl exec -it mongod-1 -c mongod-container bash
    $ mongo
    > db.getSiblingDB('admin').auth("main_admin", "abc123");
    > use test;
    > db.setSlaveOk(1);
    > db.testcoll.find();

You should see that the two records inserted via the first replica, are visible to the second replica.

#### 3.2 Redeployment Without Data Loss Test

`Warning :` If the version of your kubectl CLI client is under 1.9.x, you will get the following error: no kind "GetOptions" is registered for version "apps/v1". You have to manually delete the StatefulSet via the k8s dashboard.

To see if Persistent Volume Claims really are working, run a script to drop the Service & StatefulSet (thus stopping the pods and their “mongod” containers) and then a script to re-create them again:

    $ ./delete_service.sh
    $ ./recreate_service.sh
    $ kubectl get all

As before, keep re-running the last command above, until you can see that all 3 “mongod” pods and their containers have been successfully started again. Then connect to the first container, run the Mongo Shell and query to see if the data we’d inserted into the old containerised replica-set is still present in the re-instantiated replica set:

    $ kubectl exec -it mongod-0 -c mongod-container bash
    $ mongo
    > db.getSiblingDB('admin').auth("main_admin", "abc123");
    > use test;
    > db.testcoll.find();

You should see that the two records inserted earlier, are still present.

## 4 - Undeploying & Cleaning Down the Kubernetes Environment

`Warning :` If the version of your kubectl CLI client is under 1.9.x, you will get the following error: no kind "GetOptions" is registered for version "apps/v1". You have to manually delete the StatefulSet via the k8s dashboard.

Run the following script to undeploy the MongoDB Service & StatefulSet.

    $ ./teardown.sh

If you want, you can shutdown the Minikube virtual machine with the following command.

    $ minikube stop
