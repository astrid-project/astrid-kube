package main

import (
	core_v1 "k8s.io/api/core/v1"

	"os"
	"os/signal"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// TODO: absolutely change this
	kubeconfig = "/home/elis/.kube/config"
)

var (
	logger        *log.Logger
	signalChan    chan os.Signal
	stopInformers chan struct{}
	cleanupDone   chan struct{}
)

func main() {
	logger = log.New()
	l := log.NewEntry(logger)
	l.Infoln("Starting...")

	//----------------------------------------
	//	Start
	//----------------------------------------
	clientset := getClientSet()
	informer := getInformer(clientset)
	signalChan = make(chan os.Signal, 1)
	stopInformers = make(chan struct{})

	go informer.Run(stopInformers)

	cleanupDone = make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go cleanUp()
	<-cleanupDone
}

func getClientSet() kubernetes.Interface {
	//	Use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	//	Get the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func getInformer(clientset kubernetes.Interface) cache.SharedIndexInformer {
	l := log.NewEntry(logger)

	//	Get the informer
	informer := cache.NewSharedIndexInformer(&cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return clientset.CoreV1().Namespaces().List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			return clientset.CoreV1().Namespaces().Watch(options)
		},
	},
		&core_v1.Pod{},
		0, //Skip resync
		cache.Indexers{},
	)

	//	Set the events
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			l.Infoln("new namespace!")
		},
		UpdateFunc: func(old, new interface{}) {
		},
		DeleteFunc: func(obj interface{}) {

		},
	})

	return informer
}

func cleanUp() {
	l := log.NewEntry(logger)
	<-signalChan
	close(stopInformers)
	l.Infoln("Received an interrupt, stopping everything")
	//cleanup(services, c)
	close(cleanupDone)
}
