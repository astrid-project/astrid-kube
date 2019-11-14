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
	CB       CBEndPoints       `yaml:"cb"`
}

type CBEndPoints struct {
	Configuration string `yaml:"configuration"`
}

type VerekubeEndPoints struct {
	InfrastructureInfo  string `yaml:"infrastructure-info"`
	InfrastructureEvent string `yaml:"infrastructure-event"`
}

type Formats struct {
	InfrastructureInfo  EncodingType `yaml:"infrastructure-info"`
	InfrastructureEvent EncodingType `yaml:"infrastructure-event"`
}

type Paths struct {
	Kubeconfig string `yaml:"kubeconfig"`
}
