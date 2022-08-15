package main

import (
	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"

	"github.com/ev2-1/minetest-go-plugins/ao_mgr/ao"
)

func Tick() {
	ao.Tick()
}

func PktTick() {
	ao.SendPkts()
}

func JoinHook(clt *minetest.Client) {
	ao.JoinHook(clt)
}

func ProcessPkt(clt *minetest.Client, pkt *mt.Pkt) {
	ao.ProcessPkt(clt, pkt)
}
