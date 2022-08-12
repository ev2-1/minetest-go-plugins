package main

import (
	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"
	"github.com/ev2-1/minetest-go-plugins/tools/pos"

)

var Name string = "tools"

func ProcessPkt(c *minetest.Client, pkt *mt.Pkt) {
	switch cmd := pkt.Cmd.(type) {
	case *mt.ToSrvPlayerPos:
		pos.Update(c, &cmd.Pos)
	}
}

func LeaveHook(l *minetest.Leave) {
	pos.LeaveHook(l)
}
