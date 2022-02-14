package main

import (
	"fmt"
	v1 "k8s-watch/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

var CfgPath string

func main() {
	if home := homedir.HomeDir(); home != "" {
		CfgPath = filepath.Join(home, ".kube", "config")
	} else {
		CfgPath = "/root/.kube/config"
	}

	watchLoop()
}

func toUser(obj interface{}) (*v1.User, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	u := &v1.User{}

	err = json.Unmarshal(data, u)
	return u, err
}

func startWatching(stopCh <-chan struct{}, s cache.SharedIndexInformer) {
	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			u, _ := toUser(obj)
			fmt.Printf("ADD: %s %s %s\n", u.Namespace, u.Name, u.Spec)
		},
		UpdateFunc: func(oldObj, obj interface{}) {
			oldUser, _ := toUser(oldObj)
			u, _ := toUser(obj)
			fmt.Printf("MOD: %s %s %s  =>  %s %s %s\n",
				oldUser.Namespace, oldUser.Name, oldUser.Spec, u.Namespace, u.Name, u.Spec)
		},
		DeleteFunc: func(obj interface{}) {
			u, _ := toUser(obj)
			fmt.Printf("DEL: %s %s %s\n", u.Namespace, u.Name, u.Spec)
		},
	}
	s.AddEventHandler(handlers)
	s.Run(stopCh)
}

func watchLoop() {
	cfg, err := clientcmd.BuildConfigFromFlags("", CfgPath)
	if err != nil {
		panic(err)
	}
	client, err := dynamic.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	f := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client,
		0, "", nil)
	i := f.ForResource(schema.GroupVersionResource{
		Group:    "ir0.cn",
		Version:  "v1",
		Resource: "users"})

	startWatching(make(chan struct{}), i.Informer())
}
