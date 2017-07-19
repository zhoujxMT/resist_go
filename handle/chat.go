package handle

import (
	"sync"
)

type Chat struct {
	sync.Mutex
	Rooms map[string]*Room
}
