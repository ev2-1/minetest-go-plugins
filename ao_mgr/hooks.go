package main

import (
//	"github.com/ev2-1/minetest-go"
	"github.com/ev2-1/minetest-go-plugins/basic_ao/api"
)

/*func Tick() {
	api.Tick()
}*/

func PktTick() {
	ao.SendPkts()
}
