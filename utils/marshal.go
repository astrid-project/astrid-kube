package utils

import (
	"encoding/json"
	"encoding/xml"
	"errors"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/SunSince90/ASTRID-kube/types"
)

func toYAML(target interface{}) ([]byte, string, error) {
	data, err := yaml.Marshal(&target)
	if err != nil {
		log.Errorln("Cannot marshal to yaml:", err)
		return nil, "", err
	}
	return data, "application/yaml", nil
}

func toXML(target interface{}) ([]byte, string, error) {
	data, err := xml.MarshalIndent(&target, "", "   ")
	if err != nil {
		log.Errorln("Cannot marshal to xml:", err)
		return nil, "", err
	}
	return data, "application/xml", nil
}

func toJSON(target interface{}) ([]byte, string, error) {
	data, err := json.MarshalIndent(&target, "", "   ")
	if err != nil {
		log.Errorln("Cannot marshal to json:", err)
		return nil, "", err
	}
	return data, "application/json", nil
}

func Marshal(to types.EncodingType, target interface{}) ([]byte, string, error) {
	switch to {
	case types.XML:
		return toXML(target)
	case types.YAML:
		return toYAML(target)
	case types.JSON:
		return toJSON(target)
	}

	return nil, "", errors.New("Unrecognized format")
}
