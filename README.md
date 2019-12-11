@AUTHOR @SunSince90  Elis Lulja

# ASTRID-kube

ASTRID-kube watches for changes in Kubernetes related to namespaces and pods and later sends the resulting infrastructure to all modules interested in knowing it.

## Installation

After cloning the repository, you can just run the provided executable file ``ASTRID-kube``.
For convenience, you can add its path to ``PATH``:

```bash
$ PATH=$PATH:<path-to-ASTRID-kube>
```

NOTE: make sure to make the appropriate edits to the ``conf.yaml`` in the ``settings`` folder.

#### Configuration

Below is a brief explanation on the ``conf.yaml`` configuration file:

* ``fwInitTimer``: how many seconds to wait before creating the firewall when a pod is detected to be running. Unstable pods may compromise the stability of the rest of the graph, so this field must be set to a reasonable value to wait for any crashes to happen and to wait for all sidecars inside it to finit initializing.
* ``paths.kubeconfig``: if your kubeconfig file resides in the default folder, leave this empty. Otherwise, please fill this field accordingly.
* ``endpoints.verekube.infrastructure-info``: the endpoint where to send the resulting infrastructure. Usually, this is in the already provided format, you should only edit the provided ip with that of your machine running ``verekube``.
* ``endpoints.verekube.infrastructure-event`` (experimental): the endpoint where to send updates about the infrastructure.
* ``endpoints.cb.configuration``: the endpoint where the ``cb`` (the firewall rules pusher) is running.
* ``formats.infrastructure-info``: specify the format you want the infrastructure information to be sent as. Accepted values are ``xml``, ``yaml`` or ``json``.
* ``formats.infrastructure-event``: specify the format you want updates about the infrastructure to be sent as. Accepted values are ``xml``, ``yaml`` or ``json``.

## Usage 

Once running, ASTRID-kube will keep watching for changes in Kubernetes. To work properly, you need to make some adjustments to a couple of Kubernetes resources.

#### Graphs

Kubernetes does not have a concept of Graphs, but namespaces are the only thing that comes close to that definition. ASTRID-kube abstracts the idea of graphs with Kubernetes namespaces, and in order to be able to detect which namespaces you want to be managed by ASTRID-kube, you need to write an appropriate ``annotation`` under the namespace's ``metadata`` before deploying it containing a json list of all the deployments contained inside it:

```yaml
kind: Namespace
apiVersion: v1
metadata:
  name: mygraph
  labels:
    name: mygraph
  annotations:
    astrid.io/deployments: "[\"simple-service\", \"nodejs\", \"apache\"]"
```

Please make sure the names in the list match exactly the name of the corresponding deployment, otherwise ASTRID-kube will wait indefinitely for the applications to appear. Additionally, make sure the list is definitive, as this will signal ASTRID-kube the infrastructure you want to be built. All applications deployed later will be ignored.

#### Security Components

Once running, all applications will be protected with the appropriate security components, as specified in the ``astrid.io/security-components`` annotation. This is a json list of all security functions that the application needs. As of now, ``firewall`` is supported, but other components will be available soon. 

Take a look at the following deployment, which needs to be protected with a firewall.

```yaml
apiVersion: apps/v1 
kind: Deployment
metadata:
  name: simple-service
  namespace: mygraph
  annotations:
    astrid.io/security-components: "[\"firewall\"]"
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 1
  template:
    metadata:
      labels:
        app: simple-service
    spec:
      containers:
      - name: simple-service
        image: asimpleidea/simple-service:latest
        env:
        - name: APP_NAME
          value: "example"
        ports:
        - containerPort: 80
```

## Polycube 

ASTRID-kube relies on [Polycube](https://github.com/polycube-network/polycube) to instantiate all the proper network functions and, to do so, polycube must be injected as a sidecar in your applications.  

#### Automatic sidecar injection
The [polycube sidecar injector](https://github.com/SunSince90/polycube-sidecar-injector) is strongly recommended, as it will perform this job automatically for you. To know more about the polycube sidecar injector, please refer to the provided link. 

Once installed, add the following **label** to the example namespace provided above:

```yaml
template:
  metadata:
    polycube.network/sidecar: enabled
```

Additionally, the same key/pari must be added as **annotation** under ``spec.template.metadata`` of your deployments. Example: 

```yaml
polycube.network/sidecar: enabled
```

#### Manual sidecar injection

If you want to inject polycube manually, you have to add the following container to your deployments, under ``spec.template.containers``:

```yaml
- name: polycubed
  image: polycubenetwork/polycube:latest
  imagePullPolicy: Always
  command: ["polycubed", "--loglevel=DEBUG", "--addr=0.0.0.0", "--logfile=/host/var/log/pcn_k8s"]
  volumeMounts:
  - name: lib-modules
    mountPath: /lib/modules
  - name: usr-src
    mountPath: /usr/src
  - name: cni-path
    mountPath: /host/opt/cni/bin
  - name: etc-cni-netd
    mountPath: /host/etc/cni/net.d
  - name: var-log
    mountPath: /host/var/log
  securityContext:
    privileged: true
  ports:
    - name: polycubed
      containerPort: 9000
  terminationMessagePolicy: FallbackToLogsOnError
```

Add the following volumes, under ``spec.template.volumes`` (if not present, please add it)

```yaml
- name: lib-modules
  hostPath:
    path: /lib/modules
- name: usr-src
  hostPath:
    path: /usr/src
- name: cni-path
  hostPath:
    path: /opt/cni/bin
- name: etc-cni-netd
  hostPath:
    path: /etc/cni/net.d
- name: var-log
  hostPath:
    path: /var/log
- name: netns
  hostPath:
    path: /var/run/netns
- name: proc
  hostPath:
    path: /proc/
```    

## Examples 

The ``artifacts`` folder contains two examples that you can deploy after running ASTRID-kube. Deploy ``graph-psi.yaml`` if you are using ``polycube-sidecar-injector``; or ``graph.yaml`` if you are injecting polycube manually. 
