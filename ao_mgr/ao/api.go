package ao

import (
	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"

	"fmt"
)

// the first one is reserved for the playerAOID
const lowestAOID = mt.AOID(1)

func GetAOID() mt.AOID {
	aosMu.Lock()
	defer aosMu.Unlock()

	for id := GlobalAOIDmin; id < GlobalAOIDmax; id++ {
		if _, ok := aos[id]; !ok {
			aos[id] = Global

			return id
		}
	}

	return 0
}

func FreeAOID(id mt.AOID) {
	aosMu.Lock()
	defer aosMu.Unlock()

	delete(aos, id)
}

func RmAO(ids ...mt.AOID) {
	for _, id := range ids {
		if id == 0 {
			continue
		}

		FreeAOID(id)
		rmQueue = append(rmQueue, id)
	}
}

func AOMsg(msgs ...mt.IDAOMsg) {
	for _, msg := range msgs {
		if msg.ID == 0 {
			continue
		}

		globalMsgsMu.RLock()
		globalMsgs = append(globalMsgs, msg)
		globalMsgsMu.RUnlock()
	}
}

// - abstr -

// RegisterAO registers a initialized ActiveObject
func RegisterAO(ao ActiveObject) mt.AOID {
	if ao.GetID() == 0 {
		ao.SetID(GetAOID())
	}

	activeObjectsMu.Lock()
	activeObjects[ao.GetID()] = ao
	activeObjectsMu.Unlock()

	return ao.GetID()
}

var ao0maker func(clt *minetest.Client) mt.AOInitData

// Register player AO0 / self
// RegisterSelfAOMaker is used to register the AO maker for each client
func RegisterAO0Maker(f func(clt *minetest.Client) mt.AOInitData) {
	if ao0maker == nil {
		ao0maker = f
	} else {
		panic(fmt.Errorf("[ao_mgr] Repeated AO0Maker registration attempt."))
	}
}
