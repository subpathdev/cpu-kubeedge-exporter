package main

import (
	"log"
	"os"
	"strconv"

	"k8s.io/apimachinery/pkg/watch"

	flag "github.com/jessevdk/go-flags"

	"github.com/subpathdev/cpu-kubeedge-exporter/kubernetes"
	"github.com/subpathdev/cpu-kubeedge-exporter/prometheus"
)

var opts struct {
	Server     string `short:"s" long:"server" required:"no" description:"kubernetes address of the kubernetes api server"`
	ConfigPath string `short:"c" long:"configPath" required:"no" description:"path of the kuberentes config"`
	Address    string `short:"a" long:"address" required:"yes" description:"listen address of the webserver"`
	Port       int    `short:"p" long:"port" required:"yes" description:"listen port of the webserver"`
}

func main() {
	_, err := flag.Parse(&opts)
	if err != nil {
		print := flag.WroteHelp(err)
		args := []string{
			"-h",
		}
		if !print {
			_, err := flag.ParseArgs(&opts, args)
			if err != nil {
				log.Panicf("unexpected error, err is: %v", err)
			}
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	b10 := []byte(":")
	b10 = strconv.AppendInt(b10, int64(opts.Port), 10)
	listen := opts.Address + string(b10)

	events := make(chan watch.Event)

	if err := kubernetes.Init(opts.Server, opts.ConfigPath, events); err != nil {
		log.Panicf("clould not run successfully")
	}
	prometheus.Init(events, listen)
}
