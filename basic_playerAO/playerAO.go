package main

import (
	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"

	"errors"
	"image/color"
	"log"
	"plugin"
	"sync"
)

var errAOID = errors.New("cant assign AOID")

var playerInitialized = make(map[*minetest.Client]mt.AOID)
var playerInitializedMu sync.RWMutex

var GetPos func(c *minetest.Client) mt.PlayerPos
var _GetAOID func() *mt.AOID

func GetAOID() mt.AOID {
	id := _GetAOID()
	if id == nil {
		panic(errAOID)
	}

	return *id
}

func PluginsLoaded(m map[string]*plugin.Plugin) {
	tools, ok := m["tools"]
	if !ok {
		log.Fatal("No tools installed")
	}

	//func GetPos(c *minetest.Client) mt.PlayerPos {
	f, err := tools.Lookup("GetPos")
	if err != nil {
		log.Fatal("Tool plugin does not expose 'GetPos' function")
	}

	gp, ok := f.(func(*minetest.Client) mt.PlayerPos)
	if !ok {
		log.Fatal("tools.GetPos has incompatible type")
	}

	GetPos = gp

	ao, ok := m["ao"]
	if !ok {
		log.Fatal("AO manager not installed")
	}

	//func GetPos(c *minetest.Client) mt.PlayerPos {
	f, err = ao.Lookup("GetAOID")
	if err != nil {
		log.Fatal("Tool plugin does not expose 'GetPos' function")
	}

	ga, ok := f.(func() *mt.AOID)
	if !ok {
		log.Fatal("ao.GetAOID has incompatible type")
	}

	_GetAOID = ga
}

func ProcessPkt(c *minetest.Client, pkt *mt.Pkt) {
	switch pkt.Cmd.(type) {
	//	case *mt.ToSrvCltReady:
	case *mt.ToSrvPlayerPos:
		playerInitializedMu.RLock()
		defer playerInitializedMu.RUnlock()
		if _, ok := playerInitialized[c]; !ok {
			go initPlayer(c)
		}
	}
}

func LeaveHook(l *minetest.Leave) {
	playerInitializedMu.Lock()
	defer playerInitializedMu.Unlock()

	delete(playerInitialized, l.Client)

	// tell all clients player is gone now
}

func PosUpdate(clt *minetest.Client, p *mt.PlayerPos, _ int64) {
	playerInitializedMu.RLock()
	defer playerInitializedMu.RUnlock()

	cmd := &mt.ToCltAOMsgs{
		Msgs: []mt.IDAOMsg{
			mt.IDAOMsg{
				ID: playerInitialized[clt],
				Msg: &mt.AOCmdPos{
					Pos: mt.AOPos{
						Pos: p.Pos(),
						Rot: mt.Vec{0, p.Yaw()},

						Interpolate: true,
					},
				},
			},
			mt.IDAOMsg{
				ID: playerInitialized[clt],
				Msg: &mt.AOCmdBonePos{
					Bone: "Head_Control",
					Pos: mt.AOBonePos{
						Pos: mt.Vec{0, 6.3, 0},
						Rot: mt.Vec{-p.Pitch(), 0, 0},
					},
				},
			},
		},
	}

	for c := range playerInitialized {
		if c != clt {
			c.SendCmd(cmd)
		}
	}
}

func initPlayer(clt *minetest.Client) {
	playerInitializedMu.Lock()
	playerInitialized[clt] = GetAOID()
	playerInitializedMu.Unlock()

	playerInitializedMu.RLock()
	defer playerInitializedMu.RUnlock()

	var add []mt.AOAdd

	// self:
	newClt := playerAO(clt, false)
	newCmd := &mt.ToCltAORmAdd{
		Add: []mt.AOAdd{newClt},
	}

	for c, _ := range playerInitialized {
		if c != clt {
			add = append(add, playerAO(c, false))

			// send new client to players:
			c.SendCmd(newCmd)
		}
	}

	clt.Log("sending", len(add), "adds")

	// TODO fix this: send new client own thing
	ack, _ := clt.SendCmd(&mt.ToCltAORmAdd{
		Add: []mt.AOAdd{playerAO(clt, true)},
	})
	<-ack

	// send new client all data:
	clt.SendCmd(&mt.ToCltAORmAdd{
		Add: add,
	})
}

func playerAO(c *minetest.Client, self bool) mt.AOAdd {
	var pos mt.PlayerPos
	var aoid mt.AOID
	if !self {
		pos = GetPos(c)
		aoid = playerInitialized[c]
	} else {
		aoid = 0
	}

	name := c.Name

	return mt.AOAdd{
		ID: aoid,
		InitData: mt.AOInitData{
			Name:     name,
			IsPlayer: true,

			ID: aoid,

			Pos: pos.Pos(),
			Rot: mt.Vec{0, pos.Yaw()},

			HP: 20,

			Msgs: []mt.AOMsg{
				&mt.AOCmdProps{
					Props: mt.AOProps{
						MaxHP:      20,
						ColBox:     mt.Box{mt.Vec{-0.312, 0, -0.312}, mt.Vec{0.312, 1.8, 0.312}},
						SelBox:     mt.Box{mt.Vec{-0.312, 0, -0.312}, mt.Vec{0.312, 1.8, 0.312}},
						Pointable:  true,
						Visual:     "mesh",
						VisualSize: mt.Vec{1, 1, 1},

						Visible:  true,
						Textures: []mt.Texture{"mcl_skins_character_1.png", "blank.png", "blank.png"},

						SpriteSheetSize:  [2]int16{1, 1},
						SpritePos:        [2]int16{0, 0},
						MakeFootstepSnds: true,
						RotateSpeed:      0,
						Mesh:             "mcl_armor_character_female.b3d",

						Colors: []color.NRGBA{color.NRGBA{R: 255, G: 255, B: 255, A: 255}},

						CollideWithAOs: true,
						StepHeight:     6,
						NametagColor:   color.NRGBA{R: 255, G: 255, B: 255, A: 255},

						FaceRotateSpeed: -1,
						MaxBreath:       10,
						EyeHeight:       1.5,
						Shaded:          true,
						ShowOnMinimap:   true,

						//Textures:   []mt.Texture{"mcl_flowerpots_flowerpot.png", "mcl_flowerpots_flowerpot.png", "mcl_flowerpots_flowerpot.png", "mcl_flowerpots_flowerpot.png", "mcl_flowerpots_flowerpot.png", "mcl_flowerpots_flowerpot.png"},
						//Mesh:       "flowerpot.obj",
					},
				},
				&mt.AOCmdArmorGroups{
					Armor: []mt.Group{},
				},
				&mt.AOCmdAnim{
					Anim: mt.AOAnim{
						Frames: [2]int32{0, 1117650944},
						Speed:  30,
					},
				},
				&mt.AOCmdBonePos{
					Bone: "Body_Control",
					Pos: mt.AOBonePos{
						Pos: mt.Vec{0, 6.3, 0},
						Rot: mt.Vec{0, 0, 0},
					},
				},
				&mt.AOCmdBonePos{
					Bone: "Head_Control",
					Pos: mt.AOBonePos{
						Pos: mt.Vec{0, 6.3, 0},
						Rot: mt.Vec{-pos.Pitch(), 0, 0},
					},
				},
				&mt.AOCmdBonePos{
					Bone: "Arm_Right_Pitch_Control",
					Pos: mt.AOBonePos{
						Pos: mt.Vec{-3, 5.785, 0},
						Rot: mt.Vec{0, 0, 0},
					},
				},
				&mt.AOCmdBonePos{
					Bone: "Arm_Left_Pitch_Control",
					Pos: mt.AOBonePos{
						Pos: mt.Vec{3, 5.785, 0},
						Rot: mt.Vec{0, 0, 0},
					},
				},
				&mt.AOCmdBonePos{
					Bone: "Wield_Item",
					Pos: mt.AOBonePos{
						Pos: mt.Vec{-1.5, 4.9, 1.8},
						Rot: mt.Vec{135, 0, 90},
					},
				},

				&mt.AOCmdAttach{
					Attach: mt.AOAttach{ForceVisible: true},
				},

				&mt.AOCmdPhysOverride{
					Phys: mt.AOPhysOverride{
						Walk:    1,
						Jump:    1,
						Gravity: 1,
					},
				},
			},
		},
	}
}
