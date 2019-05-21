package main

import (
	"os"
	"os/signal"

	types "github.com/SunSince90/ASTRID-kube/types"

	"github.com/SunSince90/ASTRID-kube/informers"

	graph "github.com/SunSince90/ASTRID-kube/graph"
	"github.com/SunSince90/ASTRID-kube/settings"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	signalChan  chan os.Signal
	stop        chan struct{}
	cleanupDone chan struct{}
)

func main() {
	log.Infoln("Starting...")

	settings.Load("./settings/conf.yaml")
	log.Infoln("Configuration file loaded successfully")

	//----------------------------------------
	//	Start
	//----------------------------------------
	clientset := getClientSet()
	settings.Clientset = clientset

	//	Set up the node informer
	informers.Nodes = informers.New(types.Nodes, "").(*informers.NodeInformer)
	informers.Nodes.AddEventHandler(nil, nil, nil)
	informers.Nodes.Start()

	signalChan = make(chan os.Signal, 1)
	stop = make(chan struct{})
	graphManager := graph.InitManager(clientset, stop)
	graphManager.Start()

	cleanupDone = make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go cleanUp()
	<-cleanupDone
}

func getClientSet() kubernetes.Interface {
	//	Use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", settings.Settings.Paths.Kubeconfig)
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

func cleanUp() {
	<-signalChan
	close(stop)
	log.Infoln("Received an interrupt, stopping everything")
	//cleanup(services, c)
	close(cleanupDone)
}
