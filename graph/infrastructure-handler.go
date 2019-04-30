package graph

import (
	"errors"
	"strings"
	"sync"

	informer "github.com/SunSince90/ASTRID-kube/informers"
	astrid_types "github.com/SunSince90/ASTRID-kube/types"
	log "github.com/sirupsen/logrus"
	apps_v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type Infrastructure interface {
}

type InfrastructureHandler struct {
	clientset           kubernetes.Interface
	name                string
	log                 *log.Entry
	labels              map[string]string
	deploymentsInformer cache.SharedIndexInformer
	servicesInformer    cache.SharedIndexInformer
	stopWatching        chan struct{}
	resources           map[string]bool
	deployments         map[string]*count
	services            map[string]*core_v1.ServiceSpec
	lock                sync.Mutex
}

type count struct {
	needed  int32
	current int32
	sync.Mutex
}

type serviceInfo struct {
	nodePort   int32
	targetPort int32
}

func new(clientset kubernetes.Interface, namespace *core_v1.Namespace) (Infrastructure, error) {
	//	the handler
	inf := &InfrastructureHandler{
		name:        namespace.Name,
		labels:      namespace.Labels,
		clientset:   clientset,
		deployments: map[string]*count{},
		services:    map[string]*core_v1.ServiceSpec{},
		resources:   map[string]bool{},
		log:         log.New().WithFields(log.Fields{"GRAPH": namespace.Name}),
	}

	inf.log.Infoln("Starting graph handler for graph", namespace.Name)

	if len(namespace.Annotations) < 1 {
		inf.log.Errorln("Namespace has no annotations. Will stop here.")
		return nil, errors.New("Namespace has no annotations. Will stop here")
	}

	//	Get all deployments needed
	for name := range namespace.Annotations {
		if strings.HasPrefix(name, "astrid.io/") {
			actualName := strings.Split(name, "/")[1]
			inf.resources[actualName] = true
		}
	}

	//	First let's look at deployments
	deploymentsInformer := informer.New(astrid_types.Deployments, namespace.Name)
	deploymentsInformer.AddEventHandler(func(obj interface{}) {
		d := obj.(*apps_v1.Deployment)
		inf.handleNewDeployment(d)
	}, nil, nil)
	deploymentsInformer.Start()

	//	and then at services
	servInformer := informer.New(astrid_types.Services, namespace.Name)
	servInformer.AddEventHandler(func(obj interface{}) {
		s := obj.(*core_v1.Service)
		inf.handleNewService(s)
	}, nil, nil)
	servInformer.Start()

	return inf, nil
}

func (handler *InfrastructureHandler) handleNewDeployment(deployment *apps_v1.Deployment) {
	handler.lock.Lock()
	defer handler.lock.Unlock()

	handler.log.Infoln("Detected deployment with name:", deployment.Name)

	//	Get replicas
	handler.deployments[deployment.Name] = &count{
		needed:  *deployment.Spec.Replicas,
		current: 0,
	}

	//	Do we have all deployments?
	if len(handler.deployments) != len(handler.resources) {
		return
	}
	for deployment := range handler.resources {
		if _, exists := handler.deployments[deployment]; !exists {
			return
		}
	}

	//	ok we can start listening for pods
	handler.log.Infoln("Found all deployments needed for this graph")
}

func (handler *InfrastructureHandler) handleNewService(service *core_v1.Service) {
	handler.lock.Lock()
	defer handler.lock.Unlock()

	handler.log.Infoln("Detected service with name:", service.Name)

	handler.services[service.Name] = &service.Spec

	//	Do we have all services?
	if len(handler.services) != len(handler.resources) {
		return
	}
	for deployment := range handler.resources {
		if _, exists := handler.services[deployment]; !exists {
			return
		}
	}

	//	ok we can start listening for pods
	handler.log.Infoln("Found all services needed for this graph")
}
