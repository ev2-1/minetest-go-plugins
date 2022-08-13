package mmap

import (
	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"
)

// DO NOT CALL IF YOU DONT KNOW WHAT YOUR DOING
func ProcessPkt(c *minetest.Client, pkt *mt.Pkt) {
	switch cmd := pkt.Cmd.(type) {
	case *mt.ToSrvInteract:
		interact(cmd)
	}
}

func interact(m *mt.ToSrvInteract) {
	switch thing := m.Pointed.(type) {
	case *mt.PointedNode:
		pos := thing.Under

		switch m.Action {
		case mt.Dig:
		case mt.Dug:
			SetNode(pos, mt.Node{Param0: mt.Air})
		}

	default:
		return
	}
}
