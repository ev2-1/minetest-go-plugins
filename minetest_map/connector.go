package main

import (
	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"
	"github.com/ev2-1/minetest-go-plugins/minetest_map/api"

	"plugin"
)

func PosUpdate(c *minetest.Client, pos *mt.PlayerPos, LastUpdate int64) {
	mmap.PosUpdate(c, pos, LastUpdate)
}

func PluginsLoaded(map[string]*plugin.Plugin) {
	mmap.PluginsLoaded()
}

func ProcessPkt(c *minetest.Client, pkt *mt.Pkt) {
	mmap.ProcessPkt(c, pkt)
}
