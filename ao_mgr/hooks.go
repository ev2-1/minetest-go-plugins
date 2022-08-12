package main

import (
//	"github.com/ev2-1/minetest-go"
	"github.com/ev2-1/minetest-go-plugins/ao_mgr/ao"
)

/*func Tick() {
	api.Tick()
}*/

func PktTick() {
	ao.SendPkts()
}
