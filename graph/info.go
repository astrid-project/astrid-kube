package graph

import (
	"sync"
	"time"

	"encoding/json"
	"encoding/xml"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	types "github.com/SunSince90/ASTRID-kube/types"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type InfrastructureInfo interface {
	PushService(string, *core_v1.ServiceSpec)
	PushInstance(string, string, string)
	PopInstance(string)
	ToggleSending()
	Build(types.EncodingType)
}

type InfrastructureInfoBuilder struct {
	lock              sync.Mutex
	info              types.InfrastructureInfo
	deployedServices  map[string]int
	deployedInstances map[string]*offset
	clientset         kubernetes.Interface
	canSend           bool
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
		deployedInstances: map[string]*offset{},
		canSend:           false,
	}
}

type offset struct {
	value    string
	position int
	owner    string
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
			service.AmbassadorPort = types.InfrastructureInfoServicePort{
				Port:     9000,
				Exposed:  ports.NodePort,
				Protocol: types.TCP,
			}
		} else {
			var protocol types.InfrastructureInfoProtocol
			switch ports.Protocol {
			case core_v1.ProtocolTCP:
				protocol = types.TCP
			case core_v1.ProtocolUDP:
				protocol = types.UDP
			}

			service.Ports = append(service.Ports, types.InfrastructureInfoServicePort{
				Port:     ports.TargetPort.IntVal,
				Exposed:  ports.NodePort,
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

	existingIP, exists := i.deployedInstances[uid]
	if exists {
		if existingIP.value == ip {
			return
		}
		existingIP.value = ip
	} else {
		i.deployedInstances[uid] = &offset{
			position: len(i.info.Spec.Services[serviceOffset].Instances),
			value:    ip,
			owner:    service,
		}
	}

	i.info.Spec.Services[serviceOffset].Instances = append(i.info.Spec.Services[serviceOffset].Instances, types.InfrastructureInfoServiceInstance{
		IP:  ip,
		UID: uid,
	})

	i.send(types.XML)

}

func (i *InfrastructureInfoBuilder) PopInstance(uid string) {
	i.lock.Lock()
	defer i.lock.Unlock()

	instance, exists := i.deployedInstances[uid]
	if !exists {
		return
	}

	serviceOffset, exists := i.deployedServices[instance.owner]
	if !exists {
		return
	}

	//	Only one?
	if len(i.info.Spec.Services[serviceOffset].Instances) == 1 {
		i.info.Spec.Services[serviceOffset].Instances = []types.InfrastructureInfoServiceInstance{}
	} else {
		//	swap
		t := instance.position
		i.info.Spec.Services[serviceOffset].Instances = append(i.info.Spec.Services[serviceOffset].Instances[:t], i.info.Spec.Services[serviceOffset].Instances[t+1:]...)
	}

	i.send(types.XML)

}

func (i *InfrastructureInfoBuilder) ToggleSending() {
	i.canSend = !i.canSend
}

func (i *InfrastructureInfoBuilder) Build(to types.EncodingType) {
	nodes, err := i.clientset.CoreV1().Nodes().List(meta_v1.ListOptions{})
	if err != nil {
		log.Errorln("Cannot get nodes:", err)
		return
	}

	if len(i.info.Spec.Nodes) < 1 {
		for _, node := range nodes.Items {
			i.info.Spec.Nodes = append(i.info.Spec.Nodes, types.InfrastructureInfoNode{
				//	TODO: check this out
				IP: node.Status.Addresses[0].Address,
			})
		}
	}

	yaml := func() {
		data, err := yaml.Marshal(&i.info)
		if err != nil {
			log.Errorln("Cannot marshal to yaml:", err)
			return
		}
		log.Printf("--- t dump:\n%s\n\n", string(data))
	}

	xml := func() {
		data, err := xml.MarshalIndent(&i.info, "", "   ")
		if err != nil {
			log.Errorln("Cannot marshal to xml:", err)
			return
		}
		log.Printf("--- t dump:\n%s\n\n", string(data))
	}

	json := func() {
		data, err := json.MarshalIndent(&i.info, "", "   ")
		if err != nil {
			log.Errorln("Cannot marshal to json:", err)
			return
		}
		log.Printf("--- t dump:\n%s\n\n", string(data))
	}

	switch to {
	case types.XML:
		xml()
	case types.YAML:
		yaml()
	case types.JSON:
		json()
	}
}

func (i *InfrastructureInfoBuilder) send(to types.EncodingType) {
	if !i.canSend {
		return
	}
	nodes, err := i.clientset.CoreV1().Nodes().List(meta_v1.ListOptions{})
	if err != nil {
		log.Errorln("Cannot get nodes:", err)
		return
	}

	if len(i.info.Spec.Nodes) < 1 {
		for _, node := range nodes.Items {
			i.info.Spec.Nodes = append(i.info.Spec.Nodes, types.InfrastructureInfoNode{
				//	TODO: check this out
				IP: node.Status.Addresses[0].Address,
			})
		}
	}

	yaml := func() {
		data, err := yaml.Marshal(&i.info)
		if err != nil {
			log.Errorln("Cannot marshal to yaml:", err)
			return
		}
		log.Printf("--- t dump:\n%s\n\n", string(data))
	}

	xml := func() {
		data, err := xml.MarshalIndent(&i.info, "", "   ")
		if err != nil {
			log.Errorln("Cannot marshal to xml:", err)
			return
		}
		log.Printf("--- t dump:\n%s\n\n", string(data))
	}

	json := func() {
		data, err := json.MarshalIndent(&i.info, "", "   ")
		if err != nil {
			log.Errorln("Cannot marshal to json:", err)
			return
		}
		log.Printf("--- t dump:\n%s\n\n", string(data))
	}

	switch to {
	case types.XML:
		xml()
	case types.YAML:
		yaml()
	case types.JSON:
		json()
	}
}
