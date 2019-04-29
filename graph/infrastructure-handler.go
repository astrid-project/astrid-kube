package graph

import (
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
	//	the handler
	inf := &InfrastructureHandler{}

	return inf, nil
}
