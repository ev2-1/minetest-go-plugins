package ao

import (
	"github.com/anon55555/mt"
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

func AddAO(data mt.AOInitData) mt.AOID {
	if data.ID == 0 {
		data.ID = GetAOID()
	}

	globalAdd = append(globalAdd, mt.AOAdd{
		ID: data.ID,
		InitData: data,
	})

	return data.ID
}

func RmAO(id mt.AOID) {
	if id == 0 { return }

	activeObjectsMu.Lock()
	delete(activeObjects, id)
	activeObjectsMu.Unlock()

	FreeAOID(id)
	globalRm = append(globalRm, id)
}

func AOMsg(msg... []mt.IDAOMsg) {
	
}
