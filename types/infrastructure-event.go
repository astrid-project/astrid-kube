package types

type InfrastructureEvent struct {
	Type      InfrastructureEventType
	EventData InfrastructureEventResource
}

type InfrastructureEventResource struct {
	ResourceType InfrastructureEventResourceType
	Name         string
	Ip           string
	Uid          string
}

type InfrastructureEventType string
type InfrastructureEventResourceType string

const (
	New    InfrastructureEventType         = "new"
	Delete InfrastructureEventType         = "delete"
	Pod    InfrastructureEventResourceType = "pod"
	Node   InfrastructureEventResourceType = "node"
)
