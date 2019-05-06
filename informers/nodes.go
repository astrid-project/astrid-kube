package informers

import (
	log "github.com/sirupsen/logrus"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type NodeInformer struct {
	informer    cache.SharedIndexInformer
	stopChannel chan struct{}
}

func newNodesInformer() Informer {
	nodeInformer := &NodeInformer{
		stopChannel: make(chan struct{}),
	}

	nodeInformer.initInformer()

	return nodeInformer
}

func (nodeInformer *NodeInformer) initInformer() {
	//	Get the informer
	informer := cache.NewSharedIndexInformer(&cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return clientset.CoreV1().Nodes().List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			return clientset.CoreV1().Nodes().Watch(options)
		},
	},
		&core_v1.Node{},
		0, //Skip resync
		cache.Indexers{},
	)

	nodeInformer.informer = informer
}

func (nodeInformer *NodeInformer) Start() {
	go nodeInformer.informer.Run(nodeInformer.stopChannel)
}

func (nodeInformer *NodeInformer) Stop() {
	close(nodeInformer.stopChannel)
}

func (nodeInformer *NodeInformer) AddEventHandler(add func(interface{}), update func(interface{}, interface{}), delete func(interface{})) {
	nodeInformer.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			node := nodeInformer.parseObject(obj)
			if node != nil {
				add(node)
			}
		},
		UpdateFunc: func(old, new interface{}) {
		},
		DeleteFunc: func(obj interface{}) {
		},
	})
}

func (nodeInformer *NodeInformer) parseObject(obj interface{}) *core_v1.Node {
	//------------------------------------
	//	Try to get it
	//------------------------------------

	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Errorln("Error while trying to parse obj:", err)
		return nil
	}

	//	try to get the object
	parsedObject, _, err := nodeInformer.informer.GetIndexer().GetByKey(key)
	if err != nil {
		log.Errorf("An error occurred: cannot find cache element with key %s from store %v", key, err)
		return nil
	}

	var node *core_v1.Node
	node, ok := parsedObject.(*core_v1.Node)
	if !ok {
		node, ok = obj.(*core_v1.Node)
		if !ok {
			tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
			if !ok {
				log.Errorln("error decoding object, invalid type")
				return nil
			}
			node, ok = tombstone.Obj.(*core_v1.Node)
			if !ok {
				log.Errorln("error decoding object tombstone, invalid type")
				return nil
			}
			log.Infof("Recovered deleted object '%s' from tombstone", node.Name)
		}
	}

	//------------------------------------
	//	Add it
	//------------------------------------
	return node
}
