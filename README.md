# Kubernetes lab using Vagrant

This is a semi-scripted tutorial for setting up a Kubernetes lab using Vagrant.

```
vagrant up
vagrant ssh master
sudo cat /root/kubeadm.log
```

Run the commands you see at the end of kubeadm.log.

Also, copy the `kubeadm join` command and use it to join the worker nodes (`vagrant ssh worker1 ` and `vagrant ssh worker2`)

After you've joined both of the workers to the cluster, run this from the master:
```
vagrant@master:~$ kubectl get nodes
NAME      STATUS   ROLES    AGE   VERSION
master    Ready    master   11m   v1.15.3
worker1   Ready    <none>   99s   v1.15.3
worker2   Ready    <none>   59s   v1.15.3
```

From your Mac, build and upload the Go server:
```
./build-and-upload
```

Next, from `vagrant ssh master`:
```
docker import test-go-server.tar
kubectl run test-go-server --image test-go-server:latest --image-pull-policy IfNotPresent --port 8080 --replicas 2
kubectl expose deployment test-go-server --port 8080
kubectl get services
```

Test (kubectl needs to be run on the master, but the curl command should work from any node):
```
vagrant@master:~$ kubectl get services
NAME             TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)    AGE
kubernetes       ClusterIP   10.96.0.1     <none>        443/TCP    12m
test-go-server   ClusterIP   10.97.105.5   <none>        8080/TCP   4s
vagrant@master:~$ curl http://10.97.105.5:8080/
hello from kubernetes and vagrant! HOSTNAME=test-go-server-576f77b8b6-b7fn9
```


To expose the port to your Mac:
```
kubectl expose service test-go-server --type NodePort --name test-go-server-nodeport
kubectl get services
```

* Get the nodePort number from the `get services` output (e.g. 32004)
* From your Mac, add port forwarding for the nodePort for the NAT interface in the Virtual Box GUI for ANY one of the machines in the cluster.
* From your Mac, run:
```
curl http://localhost:32004/
```
* Notice you can do the port forwarding on ANY of the machines and it still works. 

# Misc craziness

To get a shell on a container inside the cluster
https://kubernetes.io/blog/2015/10/some-things-you-didnt-know-about-kubectl_28/
```
kubectl run -i --tty busybox --image=busybox --restart=Never -- sh
```

To generate a new command to join the cluster:
```
sudo kubeadm token create --print-join-command
```

## Create kubernetes objects

This shows how a ReplicaSet can take "ownership" of pods even if those pods were created before the ReplicaSet was created:

1. Create `pod.yaml`:
    ```
    cat <<EOF >pod.yaml
    apiVersion: v1
    kind: Pod
    metadata:
        name: frontend
        labels:
            app: guestbook
            woot: dan
    spec:
        containers:
        - name: dan-container
          image: test-go-server
          imagePullPolicy: IfNotPresent
    EOF
    ```
2. Start the pod:
    ```
    kubectl create -f pod.yaml
    ```
3. See the pod:
    ```
    kubectl get pods
    ```
4. Create `replicaset.yaml`:
    ```
    cat <<EOF >replicaset.yaml
    apiVersion: apps/v1
    kind: ReplicaSet
    metadata:
      name: frontend
      labels:
        app: guestbook
        woot: dan
    spec:
      replicas: 2
      selector:
        matchLabels:
          woot: dan
      template:
        metadata:
          labels:
            woot: dan
        spec:
          containers:
          - name: dan-container
            image: test-go-server
            imagePullPolicy: IfNotPresent
    EOF
    ```
5. Create the ReplicaSet:
    ```
    kubectl create -f replicaset.yaml
    ```
6. Repeat step 3 to see the pods again
7. Delete the ReplicaSet:
    ```
    kubectl delete replicaset frontend
    ```
8. Set that the original pod is gone by repeating step 3

# Tear down and start over

To [tear down a cluster](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#tear-down):
```
kubeadm reset
iptables -F && iptables -t nat -F && iptables -t mangle -F && iptables -X
```

# Notes

kubectl seems driven off of `~/.kube/config`

Fun command: `strace kubectl cluster-info 2> >(grep .kube)`


# References
* https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/
* https://github.com/cloudnativelabs/kube-router/blob/master/docs/kubeadm.md
* https://blog.laputa.io/kubernetes-flannel-networking-6a1cb1f8ec7c (great article on flannel)
