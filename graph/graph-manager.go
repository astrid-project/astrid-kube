package graph

import (
	"strings"
	"sync"

	"github.com/SunSince90/ASTRID-kube/informers"
	"github.com/SunSince90/ASTRID-kube/types"

	log "github.com/sirupsen/logrus"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// Manager manages all graphs (namespaces) inside the cluster
type Manager interface {
	Start()
}

// GraphManager is the implementation of the graph manager
type graphManager struct {
	clientset       kubernetes.Interface
	informer        cache.SharedIndexInformer
	stop            chan struct{}
	lock            sync.Mutex
	infrastructures map[string]Infrastructure
	nodeInformer    informers.Informer
	nodesList       map[string]bool
}

// InitManager will initialize the graph manager
func InitManager(clientset kubernetes.Interface, stop chan struct{}) Manager {
	log.Infoln("Starting graph manager")

	manager := &graphManager{
		clientset:       clientset,
		stop:            stop,
		infrastructures: map[string]Infrastructure{},
		nodeInformer:    informers.New(types.Nodes, ""),
		nodesList:       map[string]bool{},
	}

	informer := manager.getInformer()
	manager.informer = informer

	manager.nodeInformer.AddEventHandler(func(obj interface{}) {
		manager.lock.Lock()
		defer manager.lock.Unlock()

		n := obj.(*core_v1.Node)
		manager.nodesList[n.Name] = true
	}, nil, nil)
	manager.nodeInformer.Start()

	return manager
}

// Start starts the informer inside the graph manager.
func (manager *graphManager) Start() {
	go manager.informer.Run(manager.stop)
}

func (manager *graphManager) getInformer() cache.SharedIndexInformer {
	//	Get the informer
	informer := cache.NewSharedIndexInformer(&cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return manager.clientset.CoreV1().Namespaces().List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			return manager.clientset.CoreV1().Namespaces().Watch(options)
		},
	},
		&core_v1.Namespace{},
		0, //Skip resync
		cache.Indexers{},
	)

	//	Set the events
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			manager.doPreliminaryChecks(obj)
		},
		UpdateFunc: func(old, new interface{}) {
		},
		DeleteFunc: func(obj interface{}) {
		},
	})

	return informer
}

func (manager *graphManager) doPreliminaryChecks(obj interface{}) {
	//------------------------------------
	//	Try to get it
	//------------------------------------

	//	get the key
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Errorln("Error while trying to parse a graph:", err)
		return
	}

	//	try to get the object
	_ns, _, err := manager.informer.GetIndexer().GetByKey(key)
	//	Errors?
	if err != nil {
		log.Errorf("An error occurred: cannot find cache element with key %s from store %v", key, err)
		return
	}

	var ns *core_v1.Namespace

	//	Get the namespace or try to recover it (this is a very improbable case, as we're doing this just for a new event).
	ns, ok := _ns.(*core_v1.Namespace)
	if !ok {
		ns, ok = obj.(*core_v1.Namespace)
		if !ok {
			tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
			if !ok {
				log.Errorln("error decoding object, invalid type")
				return
			}
			ns, ok = tombstone.Obj.(*core_v1.Namespace)
			if !ok {
				log.Errorln("error decoding object tombstone, invalid type")
				return
			}
			log.Infof("Recovered deleted object '%s' from tombstone", ns.Name)
		}
	}

	if strings.HasPrefix(ns.Name, "kube-") || ns.Name == "default" {
		return
	}

	//------------------------------------
	//	Add it
	//------------------------------------

	manager.lock.Lock()
	defer manager.lock.Unlock()

	inf, err := new(manager.clientset, ns)
	if err != nil {
		return
	}
	manager.infrastructures[ns.Name] = inf
}
