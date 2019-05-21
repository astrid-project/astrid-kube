package types

import "time"

type Settings struct {
	EndPoints   EndPoints     `yaml:"endpoints"`
	Formats     Formats       `yaml:"formats"`
	Paths       Paths         `yaml:"paths"`
	FwInitTimer time.Duration `yaml:"fwInitTimer"`
}

type EndPoints struct {
	Verekube VerekubeEndPoints `yaml:"verekube"`
}

type VerekubeEndPoints struct {
	InfrastructureInfo  string `yaml:"infrastructure-info"`
	InfrastructureEvent string `yaml:"infrastructure-event"`
}

type Formats struct {
	InfrastructureInfo  EncodingType `yaml:"infrastructureInfo"`
	InfrastructureEvent EncodingType `yaml:"infrastructureEvent"`
}

type Paths struct {
	Kubeconfig string `yaml:"kubeconfig"`
}
