package informers

import (
	log "github.com/sirupsen/logrus"
	apps_v1 "k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type DeploymentInformer struct {
	informer    cache.SharedIndexInformer
	namespace   string
	stopChannel chan struct{}
}

func newDeploymentsInformer(namespace string) Informer {
	depInformer := &DeploymentInformer{
		namespace:   namespace,
		stopChannel: make(chan struct{}),
	}

	depInformer.initInformer()

	return depInformer
}

func (depInformer *DeploymentInformer) initInformer() {
	//	Get the informer
	informer := cache.NewSharedIndexInformer(&cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return clientset.AppsV1().Deployments(depInformer.namespace).List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			return clientset.AppsV1().Deployments(depInformer.namespace).Watch(options)
		},
	},
		&apps_v1.Deployment{},
		0, //Skip resync
		cache.Indexers{},
	)

	depInformer.informer = informer
}

func (depInformer *DeploymentInformer) Start() {
	go depInformer.informer.Run(depInformer.stopChannel)
}

func (depInformer *DeploymentInformer) Stop() {
	close(depInformer.stopChannel)
}

func (depInformer *DeploymentInformer) AddEventHandler(add func(interface{}), update func(interface{}, interface{}), delete func(interface{})) {
	depInformer.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deployment := depInformer.parseObject(obj)
			if deployment != nil {
				add(deployment)
			}
		},
		UpdateFunc: func(old, new interface{}) {
		},
		DeleteFunc: func(obj interface{}) {
		},
	})
}

func (depInformer *DeploymentInformer) parseObject(obj interface{}) *apps_v1.Deployment {
	//------------------------------------
	//	Try to get it
	//------------------------------------

	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Errorln("Error while trying to parse obj:", err)
		return nil
	}

	//	try to get the object
	parsedObject, _, err := depInformer.informer.GetIndexer().GetByKey(key)
	if err != nil {
		log.Errorf("An error occurred: cannot find cache element with key %s from store %v", key, err)
		return nil
	}

	var deployment *apps_v1.Deployment
	deployment, ok := parsedObject.(*apps_v1.Deployment)
	if !ok {
		deployment, ok = obj.(*apps_v1.Deployment)
		if !ok {
			tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
			if !ok {
				log.Errorln("error decoding object, invalid type")
				return nil
			}
			deployment, ok = tombstone.Obj.(*apps_v1.Deployment)
			if !ok {
				log.Errorln("error decoding object tombstone, invalid type")
				return nil
			}
			log.Infof("Recovered deleted object '%s' from tombstone", deployment.Name)
		}
	}

	//------------------------------------
	//	Add it
	//------------------------------------
	return deployment
}
