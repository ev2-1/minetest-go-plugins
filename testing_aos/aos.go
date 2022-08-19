package main

import (
	"github.com/ev2-1/minetest-go"
	"github.com/ev2-1/minetest-go-plugins/ao_mgr/ao"
	"github.com/ev2-1/minetest-go-plugins/tools/pos"

	"github.com/anon55555/mt"

	"image/color"
)

func init() {

}

func testAO(pos mt.Pos) ao.ActiveObject {
	return &ao.ActiveObjectS{
		AOState: ao.AOState{
			Pos: mt.AOPos{
				Pos: pos,
			},
			HP: 10,
		},

		Props: mt.AOProps{
			Mesh:      "",
			MaxHP:     10,
			Pointable: false,
			ColBox: mt.Box{
				mt.Vec{-0.5, -0.5, -0.5},
				mt.Vec{0.5, 0.5, 0.5},
			},
			SelBox: mt.Box{
				mt.Vec{-0.5, -0.5, -0.5},
				mt.Vec{0.5, 0.5, 0.5},
			},
			Visual:          "cube",
			VisualSize:      [3]float32{1.0, 1.0, 1.0},
			Textures:        []mt.Texture{"default_tnt_top.png", "default_tnt_bottom.png", "default_tnt_side.png", "default_tnt_side.png", "default_tnt_side.png", "default_tnt_side.png"},
			DmgTextureMod:   "^[brighten",
			Shaded:          true,
			SpriteSheetSize: [2]int16{1, 1},
			SpritePos:       [2]int16{0, 0},
			Visible:         true,
			Colors:          []color.NRGBA{color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}},
			BackfaceCull:    true,
			NametagColor:    color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
			NametagBG:       color.NRGBA{R: 0x01, G: 0x01, B: 0x01, A: 0x00},
			FaceRotateSpeed: -1,
			Infotext:        "",
			Itemstring:      "",
		},
	}
}

func ProcessPkt(clt *minetest.Client, pkt *mt.Pkt) {
	switch cmd := pkt.Cmd.(type) {
	case *mt.ToSrvChatMsg:
		switch cmd.Msg {
		case "spawn":
			ao.RegisterAO(testAO(pos.GetPos(clt).Pos()))

			break
		}

		break
	}
}
