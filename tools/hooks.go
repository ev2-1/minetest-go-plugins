package main

import (
	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"
	"github.com/ev2-1/minetest-go-plugins/tools/pos"

	"plugin"
)

var posUpdaters []func(*minetest.Client, *mt.PlayerPos, int64)

func broadcastPosUpdate(clt *minetest.Client, pos *mt.PlayerPos, lu int64) {
	for _, p := range posUpdaters {
		p(clt, pos, lu)
	}
}

func PluginsLoaded(pl map[string]*plugin.Plugin) {
	//func PosUpdate(c *minetest.Client, pos *mt.PlayerPos, LastUpdate int64)
	for _, p := range pl {
		s, err := p.Lookup("PosUpdate")
		if err == nil {
			f, ok := s.(func(c *minetest.Client, pos *mt.PlayerPos, LastUpdate int64))
			if ok {
				posUpdaters = append(posUpdaters, f)
			}
		}
	}

	pos.RegisterPosUpdater(broadcastPosUpdate)
}
