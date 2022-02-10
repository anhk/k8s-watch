package main

import (
	"fmt"
	v1 "k8s-watch/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	watchList := cache.NewListWatchFromClient(clientSet.CoreV1().RESTClient(),
		"users",
		"skywing",
		fields.Everything())

	_, controller := cache.NewInformer(watchList, &v1.User{}, time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				fmt.Printf("add User: %s", obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				fmt.Printf("update User from [%s] to [%s]", oldObj, newObj)
			},
			DeleteFunc: func(obj interface{}) {
				fmt.Printf("delete User: %s", obj)
			},
		})
	controller.Run(make(chan struct{}))
}
