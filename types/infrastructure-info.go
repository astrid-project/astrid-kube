package types

import "time"

type InfrastructureInfo struct {
	Kind     string                     `yaml:"kind" xml:"kind"`
	Metadata InfrastructureInfoMetadata `yaml:"metadata" yaml:"metadata"`
	Spec     InfrastructureInfoSpec     `yaml:"spec" yaml:"spec"`
}

type InfrastructureInfoMetadata struct {
	Name       string    `yaml:"name" yaml:"name"`
	LastUpdate time.Time `yaml:"lastUpdate" yaml:"lastUpdate"`
}

type InfrastructureInfoSpec struct {
	Nodes    []InfrastructureInfoNode    `yaml:"nodes" yaml:"nodes" `
	Services []InfrastructureInfoService `yaml:"services" yaml:"services"`
}

type InfrastructureInfoNode struct {
	IP string `yaml:"ip" yaml:"ip"`
}

type InfrastructureInfoService struct {
	Name           string                              `yaml:"name" yaml:"name"`
	Ports          []InfrastructureInfoServicePort     `yaml:"ports" yaml:"ports"`
	AmbassadorPort InfrastructureInfoServicePort       `yaml:"ambassadorPort" yaml:"ambassadorPort"`
	Instances      []InfrastructureInfoServiceInstance `yaml:"instances" yaml:"instances"`
}

type InfrastructureInfoServicePort struct {
	Port     int32                      `yaml:"port" yaml:"port"`
	Protocol InfrastructureInfoProtocol `yaml:"protocol" yaml:"protocol"`
	Exposed  int32                      `yaml:"exposed" yaml:"exposed"`
}

type InfrastructureInfoProtocol string

const (
	TCP  InfrastructureInfoProtocol = "TCP"
	UDP  InfrastructureInfoProtocol = "UDP"
	ICMP InfrastructureInfoProtocol = "ICMP"
	KIND string                     = "InfrastructureInfo"
)

type InfrastructureInfoServiceInstance struct {
	IP  string `yaml:"ip" yaml:"ip"`
	UID string `yaml:"uid" yaml:"uid"`
}
