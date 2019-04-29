package graph

import (
	"errors"

	core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Infrastructure interface {
}

type InfrastructureHandler struct {
	clientset kubernetes.Interface
	name      string
	labels    map[string]string
}

func new(clientset kubernetes.Interface, obj interface{}) (Infrastructure, error) {
	//	try to cast it
	ns, ok := obj.(*core_v1.Namespace)
	if !ok {
		return nil, errors.New("Error while trying to get the namespace")
	}

	//	the handler
	inf := &InfrastructureHandler{
		name:   ns.Name,
		labels: ns.Labels,
	}

	return inf, nil
}
