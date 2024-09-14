package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v3"
)

func readConfigFromYAML(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	err = yaml.Unmarshal(data, &Cfg)
	if err != nil {
		return fmt.Errorf("error unmarshaling YAML: %v", err)
	}
	Mu = &sync.RWMutex{}
	return nil
}

func WriteConfigToYAML(config *config) error {
	//Mu.Lock()
	//defer Mu.Unlock()
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling YAML: %v", err)
	}

	err = ioutil.WriteFile("mb.yml", data, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

func Load() {
	err := readConfigFromYAML("mb.yml")
	if err != nil {
		log.Fatalf("Error reading YAML: %v", err)
	}
	for i, node := range Cfg.Nodes {
		if node.Id == 0 {
			log.Fatal("Node id must not be zero")
		}
		if node.Scheme != "https" {
			Cfg.Nodes[i].Scheme = "http"
		}
	}
	WriteConfigToYAML(Cfg)
	if Cfg.WrongTopic {
		Bal = &balancer{Mu: &sync.Mutex{}}
		Bal.setupBal()
	}
}
