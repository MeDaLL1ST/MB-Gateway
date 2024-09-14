package config

import "sync"

var Cfg *config
var Mu *sync.RWMutex
var Bal *balancer

type config struct {
	Port       int    `yaml:"port"`
	PromPort   int    `yaml:"prom_port"`
	APIKey     string `yaml:"api_key"`
	WrongTopic bool   `yaml:"wrong_topic"`
	Nodes      []Node `yaml:"nodes"`
}

type Node struct {
	Id     int      `yaml:"id"`
	Topics []string `yaml:"topics"`
	IP     string   `yaml:"ip"`
	Scheme string   `yaml:"scheme"`
	APIKey string   `yaml:"api_key"`
}

type balancer struct {
	Nodes   []int
	Current int
	Mu      *sync.Mutex
}
