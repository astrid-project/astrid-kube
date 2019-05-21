package informers

var (
	Nodes *NodeInformer
)

type Informer interface {
	initInformer()
	Start()
	Stop()
	AddEventHandler(func(interface{}), func(interface{}, interface{}), func(interface{}))
}

/*func init() {
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
}*/
