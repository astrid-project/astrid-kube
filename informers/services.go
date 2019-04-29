package informers

import (
	log "github.com/sirupsen/logrus"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ServicesInformer struct {
	informer    cache.SharedIndexInformer
	namespace   string
	stopChannel chan struct{}
}

func newServicesInformer(namespace string) Informer {
	servInformer := &ServicesInformer{
		namespace:   namespace,
		stopChannel: make(chan struct{}),
	}

	servInformer.initInformer()

	return servInformer
}

func (servInformer *ServicesInformer) initInformer() {
	//	Get the informer
	informer := cache.NewSharedIndexInformer(&cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return clientset.CoreV1().Services(servInformer.namespace).List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			return clientset.CoreV1().Services(servInformer.namespace).Watch(options)
		},
	},
		&core_v1.Service{},
		0, //Skip resync
		cache.Indexers{},
	)

	servInformer.informer = informer
}

func (servInformer *ServicesInformer) Start() {
	go servInformer.informer.Run(servInformer.stopChannel)
}

func (servInformer *ServicesInformer) Stop() {
	close(servInformer.stopChannel)
}

func (servInformer *ServicesInformer) AddEventHandler(add func(interface{}), update func(interface{}, interface{}), delete func(interface{})) {
	servInformer.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			service := servInformer.parseObject(obj)
			if service != nil {
				add(service)
			}
		},
		UpdateFunc: func(old, new interface{}) {
		},
		DeleteFunc: func(obj interface{}) {
		},
	})
}

func (servInformer *ServicesInformer) parseObject(obj interface{}) *core_v1.Service {
	//------------------------------------
	//	Try to get it
	//------------------------------------

	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Errorln("Error while trying to parse obj:", err)
		return nil
	}

	//	try to get the object
	parsedObject, _, err := servInformer.informer.GetIndexer().GetByKey(key)
	if err != nil {
		log.Errorf("An error occurred: cannot find cache element with key %s from store %v", key, err)
		return nil
	}

	var service *core_v1.Service
	service, ok := parsedObject.(*core_v1.Service)
	if !ok {
		service, ok = obj.(*core_v1.Service)
		if !ok {
			tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
			if !ok {
				log.Errorln("error decoding object, invalid type")
				return nil
			}
			service, ok = tombstone.Obj.(*core_v1.Service)
			if !ok {
				log.Errorln("error decoding object tombstone, invalid type")
				return nil
			}
			log.Infof("Recovered deleted object '%s' from tombstone", service.Name)
		}
	}

	//------------------------------------
	//	Add it
	//------------------------------------
	return service
}
