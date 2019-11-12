# Astrid-kube

### Description
>Astrid-kube is a plugin for Kubernetes that monitors events from it and, more specifically, about changes in the applications deployed.

When a new Kubernetes namespace is deployed, Astrid-kube reads its specifics to know more about it like the applications that will be part of it, and will instantiate and configure all the security components they need when it detects them to be running, such as firewalls.
When all the applications inside a namespace are fully running, it notifies astrid-config component about the resulting infrastructure and waits for it to build the security configuration. Once received, it sends it to the Firewall Handler and continues to monitor changes inside it, while still watching for events for other namespaces.


## Configuration
Before starting, the *conf.yaml* file must be filled with the correct endpoints for the Firewall Handler and astrid-config, along with other useful information like the infrastructure format, which can be sent as json, yaml or xml.

## Launch
After correctly configuring the conf.yaml file, the command sudo go run *.go in the root folder must be issued.

 
The Metadata contains information about the name of the graph and last update. Spec contains information about the Nodes – the machines – and their IPs. Finally, each Service is listed, with the security component they need, the port exposed and information about all instances.

## Firewall Handler
Description
The Firewall Handler receives firewall configuration and builds the according rules to push inside the applications.

## API Resources Design
| Resources | URLs | XML repr              | Meaning                               |
|-----------|------|-----------------------|---------------------------------------|
| ROOT      | /    | SecurityConfiguration | XML file with security configuration  |

### Launch 
To start, run sudo go run *.go in the root folder.

### Source code 
Available at https://github.com/SunSince90/ASTRID-kube