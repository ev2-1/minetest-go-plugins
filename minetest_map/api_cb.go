package main

import (
	"github.com/EliasFleckenstein03/mtmap"
	//	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"

	"sync"
)

var sendHooks []func(*minetest.Client, [3]int16, *mtmap.MapBlk)
var sendHooksMu sync.RWMutex

func doSBM(c *minetest.Client, p [3]int16, blkdata *mtmap.MapBlk) {
	sendHooksMu.RLock()
	defer sendHooksMu.RUnlock()

	for _, h := range sendHooks {
		h(c, p, blkdata)
	}
}

func RegisterSBM(h func(*minetest.Client, [3]int16, *mtmap.MapBlk)) {
	sendHooksMu.Lock()
	defer sendHooksMu.Unlock()

	sendHooks = append(sendHooks, h)
}
