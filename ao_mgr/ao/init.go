package ao

import (
	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"

	"fmt"
	"sync"
)

// Data kept per client
type ClientData struct {
	clt *minetest.Client

	_AO0 bool // player AO initialized; wont send any non AOID 0 packets

	// which AOs do you have?
	aos map[mt.AOID]struct{}

	// queues
	queueAdd []mt.AOID
	queueRm  []mt.AOID
}

func (cd *ClientData) QueueAdd(adds ...mt.AOID) {
	cd.queueAdd = append(cd.queueAdd, adds...)

	for _, id := range adds {
		cd.aos[id] = struct{}{}
	}
}

func makeClientData(c *minetest.Client) *ClientData {
	return &ClientData{
		clt: c,

		aos: make(map[mt.AOID]struct{}),
	}
}

var clients = make(map[*minetest.Client]*ClientData)
var clientsMu sync.RWMutex

func JoinHook(clt *minetest.Client) {
	cd := makeClientData(clt)

	// give client data
	clientsMu.Lock()
	clients[clt] = cd
	clientsMu.Unlock()

	if ao0maker == nil {
		panic(fmt.Errorf("no AO0Maker registerd, please ensure you have a player managing plugin installed."))
	}

	// send client AO0
	cd.clt.SendCmd(&mt.ToCltAORmAdd{
		Add: []mt.AOAdd{
			mt.AOAdd{
				ID:       0,
				InitData: ao0maker(clt),
			},
		},
	})
}
