package main

import (
	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"
)

var Name = "ao"

// the first one is reserved for the playerAOID
const lowestAOID = mt.AOID(1)

func GetCltAOID(c *minetest.Client) *mt.AOID {
	aosMu.Lock()
	defer aosMu.Unlock()

	for id := mt.AOID(0); id < 2^16; id++ {
		if v, ok := aos[id]; !ok {
			aos[id] = Client
			if aosClt[id] == nil {
				aosClt[id] = make(map[*minetest.Client]struct{})
			}

			aosClt[id][c] = struct{}{}
		} else if _, ok := aosClt[id][c]; !ok && v == Client {

		}
	}

	return nil
}

func GetAOID() *mt.AOID {
	aosMu.Lock()
	defer aosMu.Unlock()

	for id := GlobalAOIDmin; id < GlobalAOIDmax; id++ {
		if _, ok := aos[id]; !ok {
			aos[id] = Global

			return &id
		}
	}

	return nil
}
