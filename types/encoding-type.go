package types

type EncodingType string

const (
	XML             EncodingType = "xml"
	YAML            EncodingType = "yaml"
	JSON            EncodingType = "json"
	ContentTypeXML  string       = "application/xml"
	ContentTypeJSON string       = "application/json"
	ContentTypeYAML string       = "application/yaml"
)
