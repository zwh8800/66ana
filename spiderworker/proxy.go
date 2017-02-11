package spiderworker

import (
	"math/rand"
	"sync"
)

type ProxyPool struct {
	lock     sync.RWMutex
	proxyMap map[string]bool
}

func NewProxyPoll(proxyList []string) *ProxyPool {
	p := &ProxyPool{
		proxyMap: make(map[string]bool, len(proxyList)),
	}
	for _, paddr := range proxyList {
		p.proxyMap[paddr] = true
	}
	return p
}

func (p *ProxyPool) GetRandomProxy() string {
	p.lock.RLock()
	defer p.lock.RUnlock()
	proxyIndex := rand.Intn(p.Length() + 1)
	if proxyIndex < p.Length() {
		for k, _ := range p.proxyMap {
			if proxyIndex == 0 {
				return k
			}
			proxyIndex--
		}
	}
	return ""
}

func (p *ProxyPool) DeleteProxy(paddr string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.proxyMap, paddr)
}

func (p *ProxyPool) Length() int {
	return len(p.proxyMap)
}

var proxyList = []string{
	"10.0.0.109:1080",
	"10.0.0.109:1081",
	"10.0.0.109:1082",
	"10.0.0.109:1083",
	"10.0.0.109:1084",
	"10.0.0.109:1085",
}
