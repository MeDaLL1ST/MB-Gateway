package config

import "log"

func (b *balancer) setupBal() {
	if len(Cfg.Nodes) == 0 {
		log.Fatal("At least one node must be entered during startup")
	}
	for i := range Cfg.Nodes {
		b.Nodes = append(b.Nodes, Cfg.Nodes[i].Id)
	}
	b.Current = 0
}

func (b *balancer) Next() int {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	defer b.incr()
	if b.Current+1 > len(b.Nodes) {
		b.Current = 0
		return b.Nodes[b.Current]
	} else {
		return b.Nodes[b.Current]
	}
}

func (b *balancer) incr() {
	b.Current++
}

func (b *balancer) AddNode(id int) {
	//b.Mu.Lock()
	b.Nodes = append(b.Nodes, id)
	//b.Mu.Unlock()
}

func (b *balancer) RmNode(id int) {
	//b.Mu.Lock()
	for i, node := range b.Nodes {
		if node == id {
			b.Nodes = append(b.Nodes[:i], b.Nodes[i+1:]...)
		}
	}
	//b.Mu.Unlock()
}
