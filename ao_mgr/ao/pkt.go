package ao

import (
	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"

	"sync"
)

var (
	globalMsgsMu sync.RWMutex
	globalMsgs []mt.IDAOMsg
	cltMsgsMu sync.RWMutex
	cltMsgs map[*minetest.Client][]mt.IDAOMsg


	globalAddMu sync.RWMutex
	globalAdd []mt.AOAdd
	cltAddMu sync.RWMutex
	cltAdd map[*minetest.Client][]mt.AOAdd

	globalRmMu sync.RWMutex
	globalRm []mt.AOID
	cltRmMu sync.RWMutex
	cltRm map[*minetest.Client][]mt.AOID
)

func SendPkts() {
	// adds / rm
	globalAddMu.RLock()
	cltAddMu.RLock()

	globalRmMu.RLock()
	cltRmMu.RLock()
	for clt := range minetest.Clts() {
		clt.SendCmd(&mt.ToCltAORmAdd{
			Add: append(globalAdd, cltAdd[clt]...),
			Remove: append(globalRm, cltRm[clt]...),
		})
	}
	globalAddMu.RUnlock()
	cltAddMu.RUnlock()

	globalRmMu.RUnlock()
	cltRmMu.RUnlock()
	
	globalAddMu.Lock()
	if len(globalAdd) != 0 {
		globalAdd = make([]mt.AOAdd, 0)
	}
	globalAddMu.Unlock()

	globalRmMu.Lock()
	if len(globalRm) != 0 {
		globalRm = make([]mt.AOID, 0)
	}
	globalRmMu.Unlock()

	cltAddMu.Lock()
	if len(cltAdd) != 0 {
		cltAdd = make(map[*minetest.Client][]mt.AOAdd)
	}
	cltAddMu.Unlock()

	cltRmMu.Lock()
	if len(cltRm) != 0 {
		cltRm = make(map[*minetest.Client][]mt.AOID)
	}
	cltRmMu.Unlock()

	// msgs
	globalMsgsMu.RLock()
	cltMsgsMu.RLock()
	for clt := range minetest.Clts() {
		clt.SendCmd(&mt.ToCltAOMsgs{
			Msgs: append(globalMsgs, cltMsgs[clt]...),
		})
	}
	globalMsgsMu.RUnlock()
	cltMsgsMu.RUnlock()

	globalMsgsMu.Lock()
	if len(globalMsgs) != 0 {
		globalMsgs = make([]mt.IDAOMsg, 0)
	}
	globalMsgsMu.Unlock()

	cltMsgsMu.Lock()
	if len(cltMsgs) != 0 {
		cltMsgs = make(map[*minetest.Client][]mt.IDAOMsg)
	}
	cltMsgsMu.Unlock()
}


