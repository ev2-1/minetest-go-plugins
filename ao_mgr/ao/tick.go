package ao

import (
	"github.com/anon55555/mt"

	"fmt"
)

var rmQueue []mt.AOID

func Tick() {
	// check if each client has all aos
	clientsMu.RLock()
	defer clientsMu.RUnlock()
	activeObjectsMu.RLock()
	defer activeObjectsMu.RUnlock()

	for _, d := range clients {
		for id, _ := range activeObjects {
			if _, ok := d.aos[id]; !ok {
				d.clt.Log(fmt.Sprintf("Adding AOID %d to client add queue", id))

				// clt dosn't have AO, adding to queue:
				d.QueueAdd(id)
			}
		}
	}
}
