package types

import "time"

type Settings struct {
	EndPoints   EndPoints     `yaml:"endpoints"`
	Formats     Formats       `yaml:"formats"`
	Paths       Paths         `yaml:"paths"`
	FwInitTimer time.Duration `yaml:"fwInitTimer"`
}

type EndPoints struct {
	Verekube string `yaml:"verekube"`
}

type Formats struct {
	InfrastructureInfo  FormatType `yaml:"infrastructureInfo"`
	InfrastructureEvent FormatType `yaml:"infrastructureEvent"`
}

type Paths struct {
	Kubeconfig string `yaml:"kubeconfig"`
}

type FormatType string

const (
	FormatXML  FormatType = "xml"
	FormatJSON FormatType = "json"
	FormatYAML FormatType = "yaml"
)
