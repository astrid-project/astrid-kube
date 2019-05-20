package types

import "time"

type InfrastructureInfo struct {
	Kind     string                     `yaml:"kind" json:"metadata" xml:"-"`
	Metadata InfrastructureInfoMetadata `yaml:"metadata" json:"metadata" xml:"Metadata"`
	Spec     InfrastructureInfoSpec     `yaml:"spec" json:"spec" xml:"Spec"`
}

type InfrastructureInfoMetadata struct {
	Name       string    `yaml:"name" json:"name" xml:"name,attr"`
	LastUpdate time.Time `yaml:"lastUpdate" json:"lastUpdate" xml:"lastUpdate,attr"`
}

type InfrastructureInfoSpec struct {
	Nodes    []InfrastructureInfoNode    `yaml:"nodes" json:"nodes" xml:"Node" `
	Services []InfrastructureInfoService `yaml:"services"  json:"services" xml:"Service"`
}

type InfrastructureInfoNode struct {
	IP string `yaml:"ip"  json:"ip" xml:"ip,attr"`
}

type InfrastructureInfoService struct {
	Name               string                                `yaml:"name"  json:"name" xml:"name,attr"`
	SecurityComponents []InfrastructureInfoSecurityComponent `yaml:"securityComponents"  json:"securityComponents" xml:"SecurityComponent"`
	Ports              []InfrastructureInfoServicePort       `yaml:"ports"  json:"ports" xml:"Port"`
	AmbassadorPort     InfrastructureInfoServicePort         `yaml:"ambassadorPort"  json:"ambassadorPort" xml:"AmbassadorPort"`
	Instances          []InfrastructureInfoServiceInstance   `yaml:"instances"  json:"instances" xml:"Instance"`
}

type InfrastructureInfoSecurityComponent struct {
	Name string `yaml:"name"  json:"name" xml:"name,attr"`
}

type InfrastructureInfoServicePort struct {
	Port     int32                      `yaml:"port"  json:"port" xml:"internal,attr"`
	Protocol InfrastructureInfoProtocol `yaml:"protocol"  json:"protocol" xml:"protocol,attr"`
	Exposed  int32                      `yaml:"exposed"  json:"exposed" xml:"exposed,attr"`
}

type InfrastructureInfoProtocol string

const (
	TCP  InfrastructureInfoProtocol = "TCP"
	UDP  InfrastructureInfoProtocol = "UDP"
	ICMP InfrastructureInfoProtocol = "ICMP"
	KIND string                     = "InfrastructureInfo"
)

type InfrastructureInfoServiceInstance struct {
	IP  string `yaml:"ip"  json:"ip" xml:"ip,attr"`
	UID string `yaml:"uid"  json:"uid" xml:"uid,attr"`
}
