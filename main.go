package main

import (
	"log"
	"mbgateway/config"
	"mbgateway/metrics"
	"mbgateway/pkg"
	"net/http"
	"strconv"
)

func main() {
	go metrics.StartMonitor()
	http.HandleFunc("/add", pkg.AddHandler)
	http.HandleFunc("/addnode", pkg.AddNodeHandler)
	http.HandleFunc("/addtopic", pkg.AddTopicHandler)
	http.HandleFunc("/rmnode", pkg.RmNodeHandler)
	http.HandleFunc("/rmtopic", pkg.RmTopicHandler)
	http.HandleFunc("/info", pkg.InfoHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Cfg.Port), nil))
}

func init() {
	config.Load()
}
