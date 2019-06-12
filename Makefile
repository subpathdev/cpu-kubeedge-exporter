.PHONY = all lint clear
.DEFAULT = all

all: main.go typ/type.go prometheus/prometheus.go kubernetes/kubernetes.go
	go build ./

clear:
	$(RM) cpu-kubeedge-exporter

lint:
	golangci-lint run ./...
	go vet ./...
