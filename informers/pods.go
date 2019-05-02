package informers

import (
	log "github.com/sirupsen/logrus"
	core_v1 "k8s.io/api/core/v1"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type PodsInformer struct {
	informer    cache.SharedIndexInformer
	namespace   string
	stopChannel chan struct{}
}

func newPodsInformer(namespace string) Informer {
	podInformer := &PodsInformer{
		namespace:   namespace,
		stopChannel: make(chan struct{}),
	}

	podInformer.initInformer()

	return podInformer
}

func (podInformer *PodsInformer) initInformer() {
	//	Get the informer
	informer := cache.NewSharedIndexInformer(&cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return clientset.CoreV1().Pods(podInformer.namespace).List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			return clientset.CoreV1().Pods(podInformer.namespace).Watch(options)
		},
	},
		&core_v1.Pod{},
		0, //Skip resync
		cache.Indexers{},
	)

	podInformer.informer = informer
}

func (podInformer *PodsInformer) Start() {
	go podInformer.informer.Run(podInformer.stopChannel)
}

func (podInformer *PodsInformer) Stop() {
	close(podInformer.stopChannel)
}

func (podInformer *PodsInformer) AddEventHandler(add func(interface{}), update func(interface{}, interface{}), delete func(interface{})) {
	podInformer.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := podInformer.parseObject(obj)
			if pod != nil && add != nil {
				add(pod)
			}
		},
		UpdateFunc: func(old, obj interface{}) {
			pod := podInformer.parseObject(obj)
			if pod != nil && update != nil {
				update(old, obj)
			}
		},
		DeleteFunc: func(obj interface{}) {
		},
	})
}

func (podInformer *PodsInformer) parseObject(obj interface{}) *core_v1.Pod {
	//------------------------------------
	//	Try to get it
	//------------------------------------

	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Errorln("Error while trying to parse obj:", err)
		return nil
	}

	//	try to get the object
	parsedObject, _, err := podInformer.informer.GetIndexer().GetByKey(key)
	if err != nil {
		log.Errorf("An error occurred: cannot find cache element with key %s from store %v", key, err)
		return nil
	}

	var pod *core_v1.Pod
	pod, ok := parsedObject.(*core_v1.Pod)
	if !ok {
		pod, ok = obj.(*core_v1.Pod)
		if !ok {
			tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
			if !ok {
				log.Errorln("error decoding object, invalid type")
				return nil
			}
			pod, ok = tombstone.Obj.(*core_v1.Pod)
			if !ok {
				log.Errorln("error decoding object tombstone, invalid type")
				return nil
			}
			log.Infof("Recovered deleted object '%s' from tombstone", pod.Name)
		}
	}

	//------------------------------------
	//	Add it
	//------------------------------------
	return pod
}
