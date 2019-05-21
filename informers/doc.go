package informers

import (
	"github.com/SunSince90/ASTRID-kube/settings"
	astrid_types "github.com/SunSince90/ASTRID-kube/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clientset kubernetes.Interface
	Nodes     *NodeInformer
)

type Informer interface {
	initInformer()
	Start()
	Stop()
	AddEventHandler(func(interface{}), func(interface{}, interface{}), func(interface{}))
}

func init() {
	//	Use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", settings.Settings.Paths.Kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	//	Get the clientset
	_clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	clientset = _clientset

	//	Set up the node informer
	Nodes = New(astrid_types.Nodes, "").(*NodeInformer)
	Nodes.Start()
}
