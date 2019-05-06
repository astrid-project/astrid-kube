package types

type InformerType string

const (
	Deployments InformerType = "deployments"
	Services    InformerType = "services"
	Pods        InformerType = "pods"
	Nodes       InformerType = "nodes"
)
