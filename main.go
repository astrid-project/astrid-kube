package main

import (
	"os"
	"os/signal"

	graph "github.com/SunSince90/ASTRID-kube/graph"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// TODO: absolutely change this
	kubeconfig = "/home/elis/.kube/config"
)

var (
	signalChan    chan os.Signal
	stopInformers chan struct{}
	cleanupDone   chan struct{}
)

func main() {
	log.Infoln("Starting...")

	//----------------------------------------
	//	Start
	//----------------------------------------
	clientset := getClientSet()
	informer := graph.GetInformer(clientset)
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

func cleanUp() {
	<-signalChan
	close(stopInformers)
	log.Infoln("Received an interrupt, stopping everything")
	//cleanup(services, c)
	close(cleanupDone)
}
