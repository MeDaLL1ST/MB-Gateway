package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mbgateway/config"
	"mbgateway/metrics"
	"net/http"
	"strconv"
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if !validateToken(authHeader) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	go metrics.Incr()
	var item AddItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	config.Mu.RLock()
	existskey := false
	for _, node := range config.Cfg.Nodes {
		for _, top := range node.Topics {
			if item.Topic == top {
				existskey = true
				err = sendKeyToNode(item.Key, item.Value, node.Scheme, node.IP, node.APIKey)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					defer config.Mu.RUnlock()
					return
				}
			}
		}
	}
	config.Mu.RUnlock()
	if !existskey && config.Cfg.WrongTopic {
		config.Mu.RLock()
		send_id := config.Bal.Next()
		for _, node := range config.Cfg.Nodes {

			if send_id == node.Id {
				err = sendKeyToNode(item.Key, item.Value, node.Scheme, node.IP, node.APIKey)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					defer config.Mu.RUnlock()
					return
				}
			}

		}
		config.Mu.RUnlock()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(item)
	} else if !existskey && !config.Cfg.WrongTopic {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(item)
	}
}
func AddNodeHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if !validateToken(authHeader) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	go metrics.Incr()
	var item AddNodeItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(item.Id)
	if err != nil {
		http.Error(w, "Id must be integer", http.StatusBadRequest)
		return
	}
	if id == 0 {
		http.Error(w, "Id dont must be zero", http.StatusBadRequest)
		return
	}
	config.Mu.RLock()
	for _, node := range config.Cfg.Nodes {
		if node.Id == id {
			http.Error(w, "Node already exists", http.StatusBadRequest)
			defer config.Mu.RUnlock()
			return
		}
	}
	config.Mu.RUnlock()

	if item.Scheme != "https" {
		item.Scheme = "http"
	}

	config.Mu.Lock()
	config.Cfg.Nodes = append(config.Cfg.Nodes, config.Node{Id: id, Scheme: item.Scheme, APIKey: item.APIKey, IP: item.Addr})
	config.WriteConfigToYAML(config.Cfg)
	if config.Cfg.WrongTopic {
		config.Bal.AddNode(id)
	}
	config.Mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}
func AddTopicHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if !validateToken(authHeader) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	go metrics.Incr()
	var item AddTopicItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(item.NodeId)
	if err != nil {
		http.Error(w, "Id must be integer", http.StatusBadRequest)
		return
	}
	config.Mu.Lock()
	for i, node := range config.Cfg.Nodes {
		if node.Id == id {
			for _, top := range node.Topics {
				if item.Topic == top {
					http.Error(w, "Node already have topic", http.StatusBadRequest)
					defer config.Mu.Unlock()
					return
				}
			}
			config.Cfg.Nodes[i].Topics = append(config.Cfg.Nodes[i].Topics, item.Topic)

		}
	}
	config.WriteConfigToYAML(config.Cfg)
	config.Mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}
func RmNodeHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if !validateToken(authHeader) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	go metrics.Incr()
	var item RmNodeItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(item.Id)
	if err != nil {
		http.Error(w, "Id must be integer", http.StatusBadRequest)
		return
	}
	config.Mu.Lock()
	for i, node := range config.Cfg.Nodes {
		if node.Id == id {
			config.Cfg.Nodes = append(config.Cfg.Nodes[:i], config.Cfg.Nodes[i+1:]...)
		}
	}
	config.WriteConfigToYAML(config.Cfg)
	if config.Cfg.WrongTopic {
		config.Bal.RmNode(id)
	}
	config.Mu.Unlock()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}
func RmTopicHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if !validateToken(authHeader) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	go metrics.Incr()
	var item RmTopicItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	config.Mu.Lock()
	for i, node := range config.Cfg.Nodes {
		for ti, top := range node.Topics {
			if item.Topic == top {
				config.Cfg.Nodes[i].Topics = append(config.Cfg.Nodes[i].Topics[:ti], config.Cfg.Nodes[i].Topics[ti+1:]...)
			}
		}
	}
	config.WriteConfigToYAML(config.Cfg)
	config.Mu.Unlock()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if !validateToken(authHeader) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	go metrics.Incr()
	config.Mu.RLock()
	response := struct {
		Nodes []config.Node `json:"nodes"`
	}{
		Nodes: config.Cfg.Nodes,
	}
	config.Mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func sendKeyToNode(key, value, scheme, addr, api_key string) error {

	if api_key == "" {
		api_key = config.Cfg.APIKey
	}

	data := map[string]interface{}{
		"key":   key,
		"value": value,
	}

	// Преобразуем data в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}
	// Отправляем POST-запрос
	req, err := http.NewRequest("GET", scheme+"://"+addr+"/add", nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Добавляем заголовок Content-Type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", api_key)
	// Добавляем данные в тело запроса
	req.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err //fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	return nil
}
