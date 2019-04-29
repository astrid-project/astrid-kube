package graph

import (
	log "github.com/sirupsen/logrus"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// GraphManager manages all graphs (namespaces) inside the cluster
type GraphManager interface {
	Start()
}

// Manager is the implementation of the graph manager
type Manager struct {
	clientset kubernetes.Interface
	informer  cache.SharedIndexInformer
	stop      chan struct{}
}

// InitManager will initialize the graph manager
func InitManager(clientset kubernetes.Interface, stop chan struct{}) GraphManager {
	log.Infoln("Starting graph manager")

	manager := &Manager{
		clientset: clientset,
		stop:      stop,
	}

	informer := manager.getInformer()
	manager.informer = informer

	return manager
}

// Start starts the informer inside the graph manager.
func (manager *Manager) Start() {
	go manager.informer.Run(manager.stop)
}

func (manager *Manager) getInformer() cache.SharedIndexInformer {
	//	Get the informer
	informer := cache.NewSharedIndexInformer(&cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return manager.clientset.CoreV1().Namespaces().List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			return manager.clientset.CoreV1().Namespaces().Watch(options)
		},
	},
		&core_v1.Pod{},
		0, //Skip resync
		cache.Indexers{},
	)

	//	Set the events
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Infoln("New namespace!")
		},
		UpdateFunc: func(old, new interface{}) {
		},
		DeleteFunc: func(obj interface{}) {
		},
	})

	return informer
}
