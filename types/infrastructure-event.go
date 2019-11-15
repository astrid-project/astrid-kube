package types

import "time"

type InfrastructureEvent struct {
	GraphName string                      `yaml:"graph-name"  json:"graphName" xml:"graphName,attr"`
	Type      InfrastructureEventType     `yaml:"type"  json:"type" xml:"type,attr"`
	EventData InfrastructureEventResource `yaml:"event-data"  json:"eventData" xml:"EventData"`
	EventTime time.Time                   `yaml:"event-time"  json:"eventTime" xml:"EventTime"`
}

type InfrastructureEventResource struct {
	ResourceType InfrastructureEventResourceType `yaml:"resource-type"  json:"resourceType" xml:"resourceType,attr"`
	Name         string                          `yaml:"name"  json:"name" xml:"name,attr"`
	Ip           string                          `yaml:"ip"  json:"ip" xml:"ip,attr"`
	Uid          string                          `yaml:"uid"  json:"uid" xml:"uid,attr"`
}

type InfrastructureEventType string
type InfrastructureEventResourceType string

const (
	New    InfrastructureEventType         = "new"
	Delete InfrastructureEventType         = "delete"
	Pod    InfrastructureEventResourceType = "pod"
	Node   InfrastructureEventResourceType = "node"
)
