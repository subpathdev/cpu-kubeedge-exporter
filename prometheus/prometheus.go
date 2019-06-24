package prometheus

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/subpathdev/cpu-kubeedge-exporter/typ"
)

type Dev struct {
	Name     string
	Actual   typ.TwinValue
	Expected typ.TwinValue
	Node     [][]string
	Operator []string
}

var devices map[string][]Dev
var devMutex, nodeMutex sync.RWMutex
var nodes map[string]int64

func handleChannel(events chan watch.Event, eve chan watch.Event) {
	for {
		select {
		case ev := <-events:
			dev, ok := ev.Object.(*typ.Device)
			if !ok {
				log.Printf("in events: can not confort ev.Object to *typ.Device; err:")
				continue
			}

			switch ev.Type {
			case watch.Deleted:
				devMutex.Lock()
				delete(devices, dev.ObjectMeta.Name)
				devMutex.Unlock()
			case watch.Added:
				devMutex.Lock()
				var nodes [][]string
				var operator []string
				for _, terms := range dev.Spec.NodeSelector.NodeSelectorTerms {
					for _, expression := range terms.MatchExpressions {
						var node []string
						node = append(node, expression.Values...)
						operator = append(operator, string(expression.Operator))
						nodes = append(nodes, node)
					}
				}
				var devs []Dev
				for _, twin := range dev.Status.Twins {
					var dev Dev
					var actual, expected typ.TwinValue
					actual = twin.Actual
					expected = twin.Desired
					dev.Actual = actual
					dev.Expected = expected
					dev.Name = twin.Name
					dev.Node = nodes
					dev.Operator = operator
					devs = append(devs, dev)
				}
				devices[dev.Name] = devs
				devMutex.Unlock()
			case watch.Modified:
				devMutex.Lock()
				var devs []Dev
				var nodes [][]string
				var operator []string
				for _, terms := range dev.Spec.NodeSelector.NodeSelectorTerms {
					for _, expression := range terms.MatchExpressions {
						var node []string
						node = append(node, expression.Values...)
						operator = append(operator, string(expression.Operator))
						nodes = append(nodes, node)
					}
				}
				for _, twin := range dev.Status.Twins {
					var dev Dev
					var actual, expected typ.TwinValue
					actual = twin.Actual
					expected = twin.Desired
					dev.Actual = actual
					dev.Expected = expected
					dev.Name = twin.Name
					dev.Node = nodes
					dev.Operator = operator
					devs = append(devs, dev)
				}
				devMutex.Unlock()
				devices[dev.Name] = devs
			default:
				log.Printf("unexpected type")
			}
		case ev := <-eve:
			dev, ok := ev.Object.(*v1.Node)
			if !ok {
				log.Printf("in eve: can not convert ev.Object to *typ.Device; err:")
				continue
			}

			switch ev.Type {
			case watch.Added:
				nodeMutex.Lock()
				nodes[dev.Name] = time.Now().Unix()
				nodeMutex.Unlock()
			case watch.Modified:
				nodeMutex.Lock()
				nodes[dev.Name] = time.Now().Unix()
				nodeMutex.Unlock()
			case watch.Deleted:
				nodeMutex.Lock()
				delete(nodes, dev.Name)
				nodeMutex.Unlock()
			default:
				log.Printf("unexpected type")
			}
		}
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	message := "Displays the matched nodes, the device, the sensor name and the value:\n"
	devMutex.RLock()
	log.Printf("request over %v devices", len(devices))
	for key, value := range devices {
		for _, v := range value {
			var node string
			for _, op := range v.Operator {
				for _, name := range v.Node {
					for _, noName := range name {
						for no := range nodes {
							switch op {
							case "In":
								if strings.Contains(no, noName) {
									node += fmt.Sprintf("%s, ", no)
								}
							case "notIn":
								if !strings.Contains(no, noName) {
									node += fmt.Sprintf("%s, ", no)
								}
							case "Exists":
								log.Printf("operator is exits; acutal not implemented\n")
								//TODO
							case "DoesNotExists":
								log.Printf("operator is exits; acutal not implemented\n")
								//TODO
							case "Gt":
								log.Printf("operator is exits; acutal not implemented\n")
								//TODO
							case "Lt":
								log.Printf("operator is exits; acutal not implemented\n")
								//TODO
							default:
								log.Printf("operator is %v and not in expected scope\n", op)
							}
						}
					}
				}
			}
			message += fmt.Sprintf("Node: %v -> %v::%v: actual value: %v\t expected value:%v\n", node, key, v.Name, v.Actual.Value, v.Expected.Value)
		}
	}
	devMutex.RUnlock()
	message += fmt.Sprintf("\n\n\n\n%v", nodes)
	if _, err := w.Write([]byte(message)); err != nil {
		log.Printf("could not write message; error is: %v", err)
	}
}

func handlePrometheus(w http.ResponseWriter, r *http.Request) {
	message := "# TYPE cpu_kubeedge_exporter gauge\n"
	devMutex.RLock()
	log.Printf("request over %v devices", len(devices))
	for _, v := range devices["cpu-sensor-tag01"] {
		message += fmt.Sprintf("cpu-sensor-tag01{node=\"%v\",sensor=\"%v\",type=\"actual\"} %v\n", v.Node, v.Name, v.Actual.Value)
		message += fmt.Sprintf("cpu-sensor-tag01{node=\"%v\",sensor=\"%v\",type=\"expected\"} %v\n", v.Node, v.Name, v.Expected.Value)
	}
	devMutex.RUnlock()
	if _, err := w.Write([]byte(message)); err != nil {
		log.Printf("could not write message; error is: %v", err)
	}
}

func Init(events chan watch.Event, listen string, eve chan watch.Event) {
	devices = make(map[string][]Dev)
	nodes = make(map[string]int64)
	go handleChannel(events, eve)

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/metrics", handlePrometheus)
	if err := http.ListenAndServe(listen, nil); err != nil {
		log.Printf("could not run list and serve; error is: %v", err)
	}
}
