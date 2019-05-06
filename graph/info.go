package graph

import (
	"sync"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	types "github.com/SunSince90/ASTRID-kube/types"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type InfrastructureInfo interface {
	PushService(string, *core_v1.ServiceSpec)
	PushInstance(string, string, string)
	Build()
}

type InfrastructureInfoBuilder struct {
	lock              sync.Mutex
	info              types.InfrastructureInfo
	deployedServices  map[string]int
	deployedInstances map[string]int
	clientset         kubernetes.Interface
}

func newBuilder(clientset kubernetes.Interface, name string) InfrastructureInfo {

	info := types.InfrastructureInfo{
		Kind: types.KIND,
		Metadata: types.InfrastructureInfoMetadata{
			Name:       name,
			LastUpdate: time.Now().UTC(),
		},
	}

	return &InfrastructureInfoBuilder{
		info:              info,
		clientset:         clientset,
		deployedServices:  map[string]int{},
		deployedInstances: map[string]int{},
	}
}

func (i *InfrastructureInfoBuilder) PushService(name string, spec *core_v1.ServiceSpec) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if _, exists := i.deployedServices[name]; exists {
		return
	}

	i.deployedServices[name] = len(i.info.Spec.Services)
	service := types.InfrastructureInfoService{
		Name: name,
	}

	for _, ports := range spec.Ports {
		if ports.Name == name+"-ambassador-port" {
			service.AmbassadorPort = ports.NodePort
		} else {
			var protocol types.InfrastructureInfoProtocol
			switch ports.Protocol {
			case core_v1.ProtocolTCP:
				protocol = types.TCP
			case core_v1.ProtocolUDP:
				protocol = types.UDP
			}

			service.ExposedPorts = append(service.ExposedPorts, ports.NodePort)
			service.InternalPorts = append(service.InternalPorts, types.InfrastructureInfoServiceInternalPort{
				Port:     ports.TargetPort.IntVal,
				Protocol: protocol,
			})
		}
	}

	i.info.Spec.Services = append(i.info.Spec.Services, service)
}

func (i *InfrastructureInfoBuilder) PushInstance(service, ip, uid string) {
	i.lock.Lock()
	defer i.lock.Unlock()

	serviceOffset, exists := i.deployedServices[service]
	if !exists {
		return
	}
	if _, exists := i.deployedInstances[uid]; exists {
		return
	}

	i.deployedInstances[uid] = len(i.info.Spec.Services[serviceOffset].Instances)
	i.info.Spec.Services[serviceOffset].Instances = append(i.info.Spec.Services[serviceOffset].Instances, types.InfrastructureInfoServiceInstance{
		IP:  ip,
		UID: uid,
	})
}

func (i *InfrastructureInfoBuilder) Build() {
	nodes, err := i.clientset.CoreV1().Nodes().List(meta_v1.ListOptions{})
	if err != nil {
		//	TODO: create error message
		return
	}

	for _, node := range nodes.Items {

		i.info.Spec.Nodes = append(i.info.Spec.Nodes, types.InfrastructureInfoNode{
			//	TODO: check this out
			IP: node.Status.Addresses[0].Address,
		})

	}
}
