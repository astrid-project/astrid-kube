package graph

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/SunSince90/ASTRID-kube/informers"

	"github.com/SunSince90/ASTRID-kube/utils"

	"github.com/SunSince90/ASTRID-kube/settings"

	log "github.com/sirupsen/logrus"

	types "github.com/SunSince90/ASTRID-kube/types"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type InfrastructureInfo interface {
	PushService(string, *core_v1.ServiceSpec, []string)
	PushInstance(string, string, string)
	PopInstance(string)
	EnableSending()
	//Build(types.EncodingType)
}

type InfrastructureInfoBuilder struct {
	lock              sync.Mutex
	info              types.InfrastructureInfo
	deployedServices  map[string]*serviceOffset
	deployedInstances map[string]*instanceOffset
	clientset         kubernetes.Interface
	sendingMode       string
	mostRecentEvent   types.InfrastructureEvent
}

type serviceOffset struct {
	securityComponents []string
	position           int
}

type instanceOffset struct {
	value    string
	position int
	owner    string
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
		deployedServices:  map[string]*serviceOffset{},
		deployedInstances: map[string]*instanceOffset{},
		sendingMode:       "",
		mostRecentEvent:   types.InfrastructureEvent{},
	}
}

func (i *InfrastructureInfoBuilder) PushService(name string, spec *core_v1.ServiceSpec, securityComponents []string) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if _, exists := i.deployedServices[name]; exists {
		return
	}

	i.deployedServices[name] = &serviceOffset{
		position:           len(i.info.Spec.Services),
		securityComponents: securityComponents,
	}
	service := types.InfrastructureInfoService{
		Name: name,
	}

	//	Put security components
	for _, sc := range securityComponents {
		service.SecurityComponents = append(service.SecurityComponents, types.InfrastructureInfoSecurityComponent{
			Name: sc,
		})
	}

	for _, ports := range spec.Ports {
		if ports.Name == name+"-ambassador-port" {
			/*service.AmbassadorPort = types.InfrastructureInfoServicePort{
				Port:     9000,
				Exposed:  ports.NodePort,
				Protocol: types.TCP,
			}*/
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

	s, exists := i.deployedServices[service]
	if !exists {
		return
	}
	serviceOffset := s.position

	existingIP, exists := i.deployedInstances[uid]
	if exists {
		if existingIP.value == ip {
			return
		}
		existingIP.value = ip
	} else {
		i.deployedInstances[uid] = &instanceOffset{
			position: len(i.info.Spec.Services[serviceOffset].Instances),
			value:    ip,
			owner:    service,
		}
	}

	i.info.Spec.Services[serviceOffset].Instances = append(i.info.Spec.Services[serviceOffset].Instances, types.InfrastructureInfoServiceInstance{
		IP:  ip,
		UID: uid,
	})

	//	Put it in the most recent event
	i.mostRecentEvent = types.InfrastructureEvent{
		Type: types.New,
		EventData: types.InfrastructureEventResource{
			ResourceType: types.Pod,
			Name:         service,
			Ip:           ip,
			Uid:          uid,
		},
	}
	i.send()
}

func (i *InfrastructureInfoBuilder) PopInstance(uid string) {
	i.lock.Lock()
	defer i.lock.Unlock()

	instance, exists := i.deployedInstances[uid]
	if !exists {
		return
	}

	s, exists := i.deployedServices[instance.owner]
	if !exists {
		return
	}
	serviceOffset := s.position

	//	Only one?
	if len(i.info.Spec.Services[serviceOffset].Instances) == 1 {
		i.info.Spec.Services[serviceOffset].Instances = []types.InfrastructureInfoServiceInstance{}
	} else {
		//	swap
		t := instance.position
		i.info.Spec.Services[serviceOffset].Instances = append(i.info.Spec.Services[serviceOffset].Instances[:t], i.info.Spec.Services[serviceOffset].Instances[t+1:]...)
	}

	i.mostRecentEvent = types.InfrastructureEvent{
		Type: types.Delete,
		EventData: types.InfrastructureEventResource{
			ResourceType: types.Pod,
			Name:         instance.owner,
			Ip:           instance.value,
			Uid:          uid,
		},
	}
	i.send()
}

func (i *InfrastructureInfoBuilder) EnableSending() {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.sendingMode = "infrastructure-info"

	//i.demoDropAll()

	//	Send immediately
	i.send()
}

func (i *InfrastructureInfoBuilder) demoDropAll() {
	ips := map[string]string{}

	for name, instance := range i.deployedInstances {
		ips[instance.value] = name
	}

	utils.DemoFakeDropAll(ips)
}

func (i *InfrastructureInfoBuilder) generate() ([]byte, string, error) {

	infrastructureInfo := func() ([]byte, string, error) {
		i.info.Metadata.LastUpdate = time.Now().UTC()
		i.info.Spec.Nodes = informers.Nodes.Current()

		data, contentType, err := utils.Marshal(settings.Settings.Formats.InfrastructureInfo, i.info)
		if err == nil {
			log.Printf("# --- Infrastructure Info to send: --- #:\n%s\n\n# --- /Infrastructure Info to send --- #", string(data))
		}
		return data, contentType, err
	}

	infrastructureEvent := func() ([]byte, string, error) {
		i.mostRecentEvent.EventTime = time.Now().UTC()
		data, contentType, err := utils.Marshal(settings.Settings.Formats.InfrastructureEvent, i.mostRecentEvent)
		log.Printf("# --- Infrastructure Event to send: --- #:\n%s\n\n# --- /Infrastructure Event to send --- #", string(data))
		return data, contentType, err
	}

	switch i.sendingMode {
	case "infrastructure-info":
		return infrastructureInfo()
	case "infrastructure-event":
		return infrastructureEvent()
	}

	return nil, "", errors.New("Unrecognized sending mode")
}

func (i *InfrastructureInfoBuilder) send() {
	if len(i.sendingMode) < 1 {
		return
	}

	defer func() {
		i.sendingMode = "infrastructure-event"
	}()

	data, contentType, err := i.generate()
	if err != nil {
		return
	}

	i.sendRequest(data, contentType)
}

func (i *InfrastructureInfoBuilder) sendRequest(data []byte, contentType string) {
	endPoint := settings.Settings.EndPoints.Verekube.InfrastructureEvent
	if i.sendingMode == "infrastructure-info" {
		endPoint = settings.Settings.EndPoints.Verekube.InfrastructureInfo

		//	Stop and explain
		log.Println("\t\t\t--- STOP")
		text2 := ""
		fmt.Scanln(&text2)
		log.Println("\t\t\t--- Resuming...")
	}

	if len(endPoint) < 1 {
		log.Errorln("No endpoint has been specified. Will stop here.")
	}

	response, err := utils.Post(endPoint, contentType, data)
	if err != nil {
		return
	}
	log.WithFields(log.Fields{"GRAPH": i.info.Metadata.Name}).Println("Sent data and received", response.StatusCode)
	defer response.Body.Close()
	i.forward(response.Body)
}

func (i *InfrastructureInfoBuilder) forward(body io.ReadCloser) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Errorln("Error in decoding data")
	}
	endPoint := settings.Settings.EndPoints.FakeCB.Configuration
	req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/xml")
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Errorln("Error while forwarding:", err)
	}
}
