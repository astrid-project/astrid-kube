package graph

import (
	log "github.com/sirupsen/logrus"
	apps_v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type Infrastructure interface {
}

type InfrastructureHandler struct {
	clientset           kubernetes.Interface
	name                string
	labels              map[string]string
	deploymentsInformer cache.SharedIndexInformer
	servicesInformer    cache.SharedIndexInformer
	stopWatching        chan struct{}
}

func new(clientset kubernetes.Interface, namespace *core_v1.Namespace) (Infrastructure, error) {
	//	the handler
	inf := &InfrastructureHandler{
		name:      namespace.Name,
		labels:    namespace.Labels,
		clientset: clientset,
	}

	deploymentsInformer := inf.getDeploymentsInformer()
	inf.deploymentsInformer = deploymentsInformer

	servicesInformer := inf.getServicesInformer()
	inf.servicesInformer = servicesInformer

	stopWatching := make(chan struct{})

	go deploymentsInformer.Run(stopWatching)
	go servicesInformer.Run(stopWatching)
	inf.stopWatching = stopWatching

	return inf, nil
}

func (handler *InfrastructureHandler) getDeploymentsInformer() cache.SharedIndexInformer {
	//	Get the informer
	informer := cache.NewSharedIndexInformer(&cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return handler.clientset.AppsV1().Deployments(handler.name).List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			return handler.clientset.AppsV1().Deployments(handler.name).Watch(options)
		},
	},
		&apps_v1.Deployment{},
		0, //Skip resync
		cache.Indexers{},
	)

	//	Set the events
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Infoln("new deployment!")
		},
		UpdateFunc: func(old, new interface{}) {
		},
		DeleteFunc: func(obj interface{}) {
		},
	})
	return informer
}

func (handler *InfrastructureHandler) getServicesInformer() cache.SharedIndexInformer {
	//	Get the informer
	informer := cache.NewSharedIndexInformer(&cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return handler.clientset.CoreV1().Services(handler.name).List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			return handler.clientset.CoreV1().Services(handler.name).Watch(options)
		},
	},
		&core_v1.Service{},
		0, //Skip resync
		cache.Indexers{},
	)

	//	Set the events
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Infoln("new service!")
		},
		UpdateFunc: func(old, new interface{}) {
		},
		DeleteFunc: func(obj interface{}) {
		},
	})
	return informer
}
