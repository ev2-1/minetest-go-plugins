package pos

import (
	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"

	"sync"
	"time"
)

var pos = make(map[*minetest.Client]*mt.PlayerPos)
var posMu sync.RWMutex

var posUpdate = make(map[*minetest.Client]int64)
var posUpdateMu sync.RWMutex

var posUpdatersMu sync.RWMutex
var posUpdaters []func(c *minetest.Client, pos *mt.PlayerPos, lu int64)

func RegisterPosUpdater(pu func(c *minetest.Client, pos *mt.PlayerPos, lu int64)) {
	posUpdatersMu.Lock()
	defer posUpdatersMu.Unlock()

	posUpdaters = append(posUpdaters, pu)
}

func Update(c *minetest.Client, p *mt.PlayerPos) {
	posUpdateMu.RLock()

	time := time.Now().UnixMicro()
	dtime := time - posUpdate[c]

	posUpdatersMu.RLock()
	for _, u := range posUpdaters {
		u(c, p, dtime)
	}
	posUpdatersMu.RUnlock()

	posUpdateMu.RUnlock()

	posUpdateMu.Lock()
	posUpdate[c] = time
	posUpdateMu.Unlock()

	posMu.Lock()
	pos[c] = p
	posMu.Unlock()
}

// GetPos returns pos os player / client
func GetPos(c *minetest.Client) mt.PlayerPos {
	posMu.RLock()
	defer posMu.RUnlock()

	if pos[c] == nil {
		pos[c] = &mt.PlayerPos{}
		pos[c].SetPos(mt.Pos{0, 100, 0})
	}

	return *pos[c]
}

// SetPos sets position
func SetPos(c *minetest.Client, p mt.PlayerPos) {
	posMu.Lock()
	defer posMu.Unlock()

	pos[c] = &p
}

// deleteClt
func LeaveHook(l *minetest.Leave) {
	posMu.Lock()
	delete(pos, l.Client)
	posMu.Unlock()

	posUpdateMu.Lock()
	delete(posUpdate, l.Client)
	posUpdateMu.Unlock()
}
