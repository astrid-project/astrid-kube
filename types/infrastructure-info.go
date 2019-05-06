package types

import "time"

type InfrastructureInfo struct {
	Kind     string                     `yaml:"kind" xml:"-"`
	Metadata InfrastructureInfoMetadata `yaml:"metadata" xml:"Metadata"`
	Spec     InfrastructureInfoSpec     `yaml:"spec" xml:"Spec"`
}

type InfrastructureInfoMetadata struct {
	Name       string    `yaml:"name" xml:"name,attr"`
	LastUpdate time.Time `yaml:"lastUpdate" xml:"lastUpdate,attr"`
}

type InfrastructureInfoSpec struct {
	Nodes    []InfrastructureInfoNode    `yaml:"nodes" xml:"Node" `
	Services []InfrastructureInfoService `yaml:"services" xml:"Service"`
}

type InfrastructureInfoNode struct {
	IP string `yaml:"ip" xml:"ip,attr"`
}

type InfrastructureInfoService struct {
	Name           string                              `yaml:"name" xml:"name,attr"`
	Ports          []InfrastructureInfoServicePort     `yaml:"ports" xml:"Port"`
	AmbassadorPort InfrastructureInfoServicePort       `yaml:"ambassadorPort" xml:"AmbassadorPort"`
	Instances      []InfrastructureInfoServiceInstance `yaml:"instances" xml:"Instance"`
}

type InfrastructureInfoServicePort struct {
	Port     int32                      `yaml:"port" xml:"internal,attr"`
	Protocol InfrastructureInfoProtocol `yaml:"protocol" xml:"protocol,attr"`
	Exposed  int32                      `yaml:"exposed" xml:"exposed,attr"`
}

type InfrastructureInfoProtocol string

const (
	TCP  InfrastructureInfoProtocol = "TCP"
	UDP  InfrastructureInfoProtocol = "UDP"
	ICMP InfrastructureInfoProtocol = "ICMP"
	KIND string                     = "InfrastructureInfo"
)

type InfrastructureInfoServiceInstance struct {
	IP  string `yaml:"ip" xml:"ip,attr"`
	UID string `yaml:"uid" xml:"uid,attr"`
}
