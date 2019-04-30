package informers

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// TODO: absolutely change this
	kubeconfig = "/home/elis/.kube/config"
)

var clientset kubernetes.Interface

type Informer interface {
	initInformer()
	Start()
	Stop()
	AddEventHandler(func(interface{}), func(interface{}, interface{}), func(interface{}))
}

func init() {
	//	Use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	//	Get the clientset
	_clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	clientset = _clientset
}
