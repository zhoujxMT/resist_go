package handle

// 票仓 主要是运用于每局的投票仓库
import (
	"sync"
)

type VoteSet struct {
	sync.RWMutex
	m map[string]bool
}

func NewVote() *VoteSet {
	return &VoteSet{
		m:       map[string]bool{},
		RWMutex: sync.RWMutex{}}
}

func (vote *VoteSet) Add(clientName string) {
	vote.Lock()
	defer vote.Unlock()
	vote.m[clientName] = true
}

func (vote *VoteSet) Remove(clientName string) {
	vote.Lock()
	defer vote.Unlock()
	delete(vote.m, clientName)
}

func (vote *VoteSet) Len() int {
	return len(vote.List())

}

func (vote *VoteSet) List() []string {
	vote.RLock()
	defer vote.RUnlock()
	list := []string{}
	for item := range vote.m {
		list = append(list, item)
	}
	return list
}

func (vote *VoteSet) Clear() {
	vote.Lock()
	defer vote.Unlock()
	vote.m = map[string]bool{}
}
