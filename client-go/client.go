package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	context := context.Background()

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("failed to find home directory for config: %v\n", err)
		return
	}

	// The ~ does not seem to work here, but hardcoding the path is not appropriate either
	kubeconfig := flag.String("kubeconfig", home+"/.kube/config", "kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("failed to read kube config: %v\n", err)
		return
	}

	// Optional config customization
	config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	config.UserAgent = fmt.Sprintf("sample/v1.0 (%s/%s) kubernetes/v1.0", runtime.GOOS, runtime.GOARCH)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("failed to create client set: %v\n", err)
		return
	}

	pods, err := clientset.CoreV1().Pods("book").List(context, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("failed to get pods: %v\n", err)
		return
	}

	fmt.Printf("pod %v\n", pods)

	informerFactory := informers.NewSharedInformerFactory(clientset, time.Minute*30)
	podInformer := informerFactory.Core().V1().Pods()

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(new interface{}) { fmt.Printf("add pod %v\n", new) },
		DeleteFunc: func(obj interface{}) { fmt.Printf("delete pod %v\n", obj) },
		UpdateFunc: func(old, new interface{}) { fmt.Printf("update pod %v\n", new) },
	})

	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
}
