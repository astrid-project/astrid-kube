package graph

import (
	"sync"
	"time"

	types "github.com/SunSince90/ASTRID-kube/types"
	core_v1 "k8s.io/api/core/v1"
)

type InfrastructureInfo interface {
	PushService(string, *core_v1.ServiceSpec)
}

type InfrastructureInfoBuilder struct {
	lock             sync.Mutex
	info             types.InfrastructureInfo
	deployedServices map[string]bool
}

func newBuilder(name string) InfrastructureInfo {

	info := types.InfrastructureInfo{
		Kind: types.KIND,
		Metadata: types.InfrastructureInfoMetadata{
			Name:       name,
			LastUpdate: time.Now().UTC(),
		},
	}

	return &InfrastructureInfoBuilder{
		info:             info,
		deployedServices: map[string]bool{},
	}
}

func (i *InfrastructureInfoBuilder) PushService(name string, spec *core_v1.ServiceSpec) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if _, exists := i.deployedServices[name]; exists {
		return
	}

	i.deployedServices[name] = true
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
}

func (i *InfrastructureInfoBuilder) Build() {
}
