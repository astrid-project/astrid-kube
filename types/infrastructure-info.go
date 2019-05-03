package types

import "time"

type InfrastructureInfo struct {
	Kind     string
	Metadata InfrastructureInfoMetadata
	Spec     InfrastructureInfoSpec
}

type InfrastructureInfoMetadata struct {
	Name       string
	LastUpdate time.Time
}

type InfrastructureInfoSpec struct {
	Nodes    []InfrastructureInfoNode
	Services []InfrastructureInfoService
}

type InfrastructureInfoNode struct {
	IP string
}

type InfrastructureInfoService struct {
	Name           string
	ExposedPorts   []int32
	AmbassadorPort int32
	InternalPorts  []InfrastructureInfoServiceInternalPort
	Instances      []InfrastructureInfoServiceInstance
}

type InfrastructureInfoServiceInternalPort struct {
	Port     int32
	Protocol InfrastructureInfoProtocol
}

type InfrastructureInfoProtocol string

const (
	TCP  InfrastructureInfoProtocol = "TCP"
	UDP  InfrastructureInfoProtocol = "UDP"
	ICMP InfrastructureInfoProtocol = "ICMP"
	KIND string                     = "InfrastructureInfo"
)

type InfrastructureInfoServiceInstance struct {
	IP  string
	UID string
}
