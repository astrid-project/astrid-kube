package types

import "time"

type InfrastructureInfo struct {
	Kind     string                     `yaml:"kind"`
	Metadata InfrastructureInfoMetadata `yaml:"metadata"`
	Spec     InfrastructureInfoSpec     `yaml:"spec"`
}

type InfrastructureInfoMetadata struct {
	Name       string    `yaml:"name"`
	LastUpdate time.Time `yaml:"lastUpdate"`
}

type InfrastructureInfoSpec struct {
	Nodes    []InfrastructureInfoNode    `yaml:"nodes"`
	Services []InfrastructureInfoService `yaml:"services"`
}

type InfrastructureInfoNode struct {
	IP string `yaml:"ip"`
}

type InfrastructureInfoService struct {
	Name           string                              `yaml:"name"`
	Ports          []InfrastructureInfoServicePort     `yaml:"ports"`
	AmbassadorPort InfrastructureInfoServicePort       `yaml:"ambassadorPort"`
	Instances      []InfrastructureInfoServiceInstance `yaml:"instances"`
}

type InfrastructureInfoServicePort struct {
	Port     int32                      `yaml:"port"`
	Protocol InfrastructureInfoProtocol `yaml:"protocol"`
	Exposed  int32                      `yaml:"exposed"`
}

type InfrastructureInfoProtocol string

const (
	TCP  InfrastructureInfoProtocol = "TCP"
	UDP  InfrastructureInfoProtocol = "UDP"
	ICMP InfrastructureInfoProtocol = "ICMP"
	KIND string                     = "InfrastructureInfo"
)

type InfrastructureInfoServiceInstance struct {
	IP  string `yaml:"ip"`
	UID string `yaml:"uid"`
}
