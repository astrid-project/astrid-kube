# ASTRID-kube

ASTRID-kube watches for changes in Kubernetes related to namespaces and pods and later sends the resulting infrastructure to all modules interested in knowing it.

## Installation

After cloning the repository, you can just run the provided executable file ``ASTRID-kube``.
For convenience, you can add its path to ``PATH``:

```bash
$ PATH=$PATH:<path-to-ASTRID-kube>
```

NOTE: make sure to make the appropriate edits to the ``conf.yaml`` in the ``settings`` folder.

### Configuration

Below is a brief explanation on the ``conf.yaml`` configuration file:

* ``fwInitTimer``: how many seconds to wait before creating the firewall when a pod is detected to be running. Unstable pods may compromise the stability of the rest of the graph, so this field must be set to a reasonable value to wait for any crashes to happen and to wait for all sidecars inside it to finit initializing.
* ``paths.kubeconfig``: if your kubeconfig file resides in the default folder, leave this empty. Otherwise, please fill this field accordingly.
* ``endpoints.verekube.infrastructure-info``: the endpoint where to send the resulting infrastructure. Usually, this is in the already provided format, you should only edit the provided ip with that of your machine running ``verekube``.
* ``endpoints.verekube.infrastructure-event`` (experimental): the endpoint where to send updates about the infrastructure.
* ``endpoints.cb.configuration``: the endpoint where the ``cb`` (the firewall rules pusher) is running.
* ``formats.infrastructure-info``: specify the format you want the infrastructure information to be sent as. Accepted values are ``xml``, ``yaml`` or ``json``.
* ``formats.infrastructure-event``: specify the format you want updates about the infrastructure to be sent as. Accepted values are ``xml``, ``yaml`` or ``json``.