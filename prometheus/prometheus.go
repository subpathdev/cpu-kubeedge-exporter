package prometheus

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"k8s.io/apimachinery/pkg/watch"

	"github.com/subpathdev/cpu-kubeedge-exporter/typ"
)

type Dev struct {
	Name     string
	Actual   typ.TwinValue
	Expected typ.TwinValue
}

var devices map[string][]Dev
var devMutex sync.RWMutex

func handleChannel(events chan watch.Event) {
	for {
		ev := <-events

		dev, ok := ev.Object.(*typ.Device)
		if !ok {
			log.Fatalf("can not confort ev.Object to *typ.Device; err:")
			return
		}

		switch ev.Type {
		case watch.Deleted:
			devMutex.Lock()
			delete(devices, dev.ObjectMeta.Name)
			devMutex.Unlock()
		case watch.Added:
			devMutex.Lock()
			var devs []Dev
			for _, twin := range dev.Status.Twins {
				var dev Dev
				var actual, expected typ.TwinValue
				actual = twin.Actual
				expected = twin.Desired
				dev.Actual = actual
				dev.Expected = expected
				dev.Name = twin.Name
				devs = append(devs, dev)
			}
			devices[dev.Name] = devs
			devMutex.Unlock()
		case watch.Modified:
			devMutex.Lock()
			fmt.Printf("todo\n")
			devMutex.Unlock()
		default:
			log.Printf("unexpected type")
		}
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	message := "Displays the device, the sensor name and the value:\n"
	devMutex.RLock()
	log.Printf("request over %v devices", len(devices))
	for key, value := range devices {
		for _, v := range value {
			message += fmt.Sprintf("%v::%v: actual value: %v\t expected value:%v\n", key, v.Name, v.Actual.Value, v.Expected.Value)
		}
	}
	devMutex.RUnlock()
	if _, err := w.Write([]byte(message)); err != nil {
		log.Printf("could not write message; error is: %v", err)
	}
}

func Init(events chan watch.Event, listen string) {
	devices = make(map[string][]Dev)
	go handleChannel(events)

	http.HandleFunc("/", handleRequest)
	if err := http.ListenAndServe(listen, nil); err != nil {
		log.Printf("could not run list and serve; error is: %v", err)
	}
}
