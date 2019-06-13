package kubernetes

import (
	"log"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	awatch "k8s.io/apimachinery/pkg/watch"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/tools/watch"

	"github.com/subpathdev/cpu-kubeedge-exporter/typ"
)

type ResourceEventHandler struct {
	events chan awatch.Event
}

func (r ResourceEventHandler) OnAdd(obj interface{}) {
	log.Printf("adding new device")
	r.obj2Event(awatch.Added, obj)
}

func (r ResourceEventHandler) OnUpdate(odlObj, obj interface{}) {
	log.Printf("modify device")
	r.obj2Event(awatch.Modified, obj)
}

func (r ResourceEventHandler) OnDelete(obj interface{}) {
	log.Printf("delete device")
	r.obj2Event(awatch.Deleted, obj)
}

func (r ResourceEventHandler) obj2Event(typ awatch.EventType, obj interface{}) {
	eventObj, ok := obj.(runtime.Object)
	if !ok {
		log.Printf("unknow type: %T, ignore", obj)
		return
	}
	r.events <- awatch.Event{Type: typ, Object: eventObj}
}

var kubernetesRestClient *rest.RESTClient

// Init will initialise the connection to kubernetes api server
// kubeMaster is the url of the master
// kubeConfig is the path to the kubeconfig
func Init(kubeMaster string, kubeConfig string, events chan awatch.Event) error {
	conf, err := clientcmd.BuildConfigFromFlags(kubeMaster, kubeConfig)
	if err != nil {
		log.Fatalf("can not connect to kubernetes api server: %v", err)
		return err
	}

	conf.ContentType = runtime.ContentTypeJSON
	conf.APIPath = "/apis"
	conf.GroupVersion = &schema.GroupVersion{Group: "devices.kubeedge.io", Version: "v1alpha1"}

	kubernetesRestClient, err = rest.RESTClientFor(conf)
	if err != nil {
		log.Fatalf("can not create REST client, error is: %v", err)
		return err
	}

	lw := cache.NewListWatchFromClient(kubernetesRestClient, "devices", "default", fields.Everything())
	si := cache.NewSharedInformer(lw, &typ.Device{}, 0)
	reh := ResourceEventHandler{events: events}
	si.AddEventHandler(reh)
	stopNever := make(chan struct{})
	go si.Run(stopNever)
	return nil
}
