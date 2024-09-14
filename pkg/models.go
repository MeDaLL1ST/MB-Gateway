package pkg

type AddItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Topic string `json:"topic"`
}

type AddNodeItem struct {
	Id     string `json:"id"`
	Addr   string `json:"addr"`
	APIKey string `json:"api_key"`
	Scheme string `json:"scheme"`
}

type RmNodeItem struct {
	Id string `json:"id"`
}

type AddTopicItem struct {
	Topic  string `json:"topic"`
	NodeId string `json:"node_id"`
}

type RmTopicItem struct {
	Topic string `json:"topic"`
}
