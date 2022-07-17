package main

import (
	"github.com/EliasFleckenstein03/mtmap"
	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"

	"plugin"
	"sync"
)

// a list of all clients and their loaded chunks
var loadedChunks map[*minetest.Client]map[pos]bool
var loadedChunksMu sync.RWMutex

// configuration
var load bool = true

var (
	MapBlkUpdateRate   int64 = 2         // in seconds
	MapBlkUpdateRange        = int16(10) // in mapblks
	MapBlkUpdateHeight       = int16(5)  // in mapblks

	heigthOff = -MapBlkUpdateHeight / 2
)

var stone mt.Content
var grass mt.Content
var exampleBlk mtmap.MapBlk

func init() {
	loadedChunks = make(map[*minetest.Client]map[pos]bool)
	OpenDB(minetest.Path("/map.sqlite"))
}

func PluginsLoaded(map[string]*plugin.Plugin) {
	minetest.FillNameIdMap()

	s := minetest.GetNodeDef("mcl_core:stone")
	if s != nil {
		stone = s.Param0
	}

	s = minetest.GetNodeDef("mcl_core:dirt_with_grass")
	if s != nil {
		grass = s.Param0
	}

	exampleBlk = mtmap.MapBlk{}

	for i := 0; i < 4096; i++ {
		exampleBlk.Param0[i] = 126
	}

	for i := 0; i < 16*16; i++ {
		exampleBlk.Param0[i] = stone
	}

	// center block is stone:
	exampleBlk.Param0[4096/2+16/2] = grass // some wool
}

func PosUpdate(c *minetest.Client, pos *mt.PlayerPos, LastUpdate int64) {
	if load {
		p := Pos2int(pos.Pos())
		blkpos, _ := mt.Pos2Blkpos(p)

		go func() {
			for _, sp := range spiral(MapBlkUpdateRange) {
				for i := int16(0); i < MapBlkUpdateRange; i++ {
					// generate absolute position
					ap := sp.add(blkpos).add([3]int16{0, heigthOff + i})

					// load block
					LoadBlk(c, ap, false)
				}
			}
		}()
	}
}
