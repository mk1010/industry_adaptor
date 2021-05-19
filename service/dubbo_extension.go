package service

import (
	"hash/crc32"
	"sync"

	json "github.com/json-iterator/go"

	"github.com/apache/dubbo-go/cluster"
	"github.com/apache/dubbo-go/cluster/loadbalance"
	"github.com/apache/dubbo-go/common/extension"
	"github.com/apache/dubbo-go/protocol"
)

const (
	// ConsistentHash consistent hash
	NcLinkConsistentHash = "nclinkConsistentHash"
)

var (
	selectors = make(map[string]*NcLinkConsistentHashSelector)
	writeMap  sync.Map
)

func init() {
	extension.SetLoadbalance(NcLinkConsistentHash, NewNcLinkConsistentHashLoadBalance)
}

// ConsistentHashLoadBalance implementation of load balancing: using consistent hashing
type NcLinkConsistentHashLoadBalance struct{}

// NewConsistentHashLoadBalance creates NewConsistentHashLoadBalance
//
// The same parameters of the request is always sent to the same provider.
func NewNcLinkConsistentHashLoadBalance() cluster.LoadBalance {
	return &NcLinkConsistentHashLoadBalance{}
}

// Select gets invoker based on load balancing strategy
func (lb *NcLinkConsistentHashLoadBalance) Select(invokers []protocol.Invoker, invocation protocol.Invocation) protocol.Invoker {
	key := invokers[0].GetUrl().ServiceKey()

	// hash the invokers
	bs := make([]byte, 0)
	for _, invoker := range invokers {
		b, err := json.Marshal(invoker)
		if err != nil {
			return nil
		}
		bs = append(bs, b...)
	}
	hashCode := crc32.ChecksumIEEE(bs)
	selector, ok := selectors[key]
	if ok {
		if selector.hashCode == hashCode {
			return selector.selectInvoker
		} else {
			historyInvoker := selector.selectInvoker
			if historyInvoker.IsAvailable() {
				targetIp := historyInvoker.GetUrl().Ip
				targetPort := historyInvoker.GetUrl().Port
				for _, invoker := range invokers {
					if invoker.GetUrl().Ip == targetIp && invoker.GetUrl().Port == targetPort {
						selector.mutex.Lock()
						selector.hashCode = hashCode
						selector.selectInvoker = invoker
						selector.mutex.Unlock()
						return invoker
					}
				}
			}
		}
	}
	m, _ := writeMap.LoadOrStore(key, &sync.Mutex{})
	mu := m.(*sync.Mutex)
	mu.Lock()
	selector, ok = selectors[key]
	if !ok || selector.hashCode != hashCode {
		selectors[key] = newNcLinkConsistentHashSelector(invokers, invocation, hashCode)
		selector = selectors[key]
	}
	mu.Unlock()
	return selector.selectInvoker
}

type NcLinkConsistentHashSelector struct {
	hashCode      uint32
	selectInvoker protocol.Invoker
	mutex         sync.Mutex
}

func newNcLinkConsistentHashSelector(invokers []protocol.Invoker, invocation protocol.Invocation,
	hashCode uint32) *NcLinkConsistentHashSelector {
	invoker := loadbalance.NewRandomLoadBalance().Select(invokers, invocation)
	return &NcLinkConsistentHashSelector{
		hashCode:      hashCode,
		selectInvoker: invoker,
	}
}
