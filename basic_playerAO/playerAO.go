package main

import (
	"github.com/anon55555/mt"
	"github.com/ev2-1/minetest-go"
	"github.com/ev2-1/minetest-go-plugins/ao_mgr/ao"
	"github.com/ev2-1/minetest-go-plugins/tools/pos"

	"errors"
	"fmt"
	"image/color"
	"sync"
	"time"
)

var errAOID = errors.New("cant assign AOID")

type playerData struct {
	ID     mt.AOID
	expire *bool // in unixseconds
	killCh chan struct{}
}

var playerInitialized = make(map[*minetest.Client]*playerData)
var playerInitializedMu sync.RWMutex

func GetAOID() mt.AOID {
	id := ao.GetAOID()
	if id == 0 {
		panic(errAOID)
	}

	return id
}

func ProcessPkt(c *minetest.Client, pkt *mt.Pkt) {
	switch cmd := pkt.Cmd.(type) {
	//	case *mt.ToSrvCltReady:
	case *mt.ToSrvPlayerPos:
		playerInitializedMu.RLock()
		defer playerInitializedMu.RUnlock()
		if playerInitialized[c] == nil {
			playerInitialized[c] = &playerData{}
			go initPlayer(c)
		}

	case *mt.ToSrvChatMsg:
		switch cmd.Msg { // return own pos
		case "pos":
			pp := pos.GetPos(c)
			pos := pp.Pos()

			c.SendCmd(&mt.ToCltChatMsg{
				Type: mt.RawMsg,

				Text: fmt.Sprintf("Your position: (%f, %f, %f) pitch: %f, yaw: %f",
					pos[0], pos[1], pos[2],
					pp.Pitch(), pp.Yaw(),
				),

				Timestamp: time.Now().Unix(),
			})

		case "activeobjects":
			var text string

			playerInitializedMu.RLock()
			defer playerInitializedMu.RUnlock()

			for c, pd := range playerInitialized {
				text += fmt.Sprintf("%s=%d; ", c.Name, pd.ID)
			}

			c.SendCmd(&mt.ToCltChatMsg{
				Type: mt.RawMsg,

				Text: text,
			})
		}
	}
}

func LeaveHook(l *minetest.Leave) {
	go func() {
		playerInitializedMu.Lock()
		defer playerInitializedMu.Unlock()

		// check if player actually existed
		data := playerInitialized[l.Client]
		if data == nil {
			return
		}

		// create pkt
		remove := &mt.ToCltAORmAdd{
			Remove: []mt.AOID{data.ID},
		}

		// tell all clients player is gone now
		for clt := range playerInitialized {
			if clt == l.Client { // skip self
				continue
			}

			clt.SendCmd(remove)
		}

		// actually delete the AOID from all caches
		ao.FreeAOID(data.ID)
		delete(playerInitialized, l.Client)
	}()
}

func PosUpdate(clt *minetest.Client, p *mt.PlayerPos, dt int64) {
	playerInitializedMu.RLock()
	defer playerInitializedMu.RUnlock()

	//clt.Log("last update is", dt, "old")

	data := playerInitialized[clt]
	if data == nil {
		playerInitialized[clt] = &playerData{}
		data = playerInitialized[clt]
		go initPlayer(clt)
	}

	d := true
	data.expire = &d // expires in one second

	id := data.ID

	cmd := mt.ToCltAOMsgs{
		Msgs: []mt.IDAOMsg{
			mt.IDAOMsg{
				ID: id,
				Msg: &mt.AOCmdPos{
					Pos: mt.AOPos{
						Pos: p.Pos(),
						Rot: mt.Vec{0, p.Yaw()},

						Interpolate: true,
					},
				},
			},
			mt.IDAOMsg{
				ID: id,
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
		if c != clt && c.State == minetest.CsActive {
			c.SendCmd(&cmd)
		}
	}
}

func init() {
	ao.RegisterAO0Maker(func(clt *minetest.Client) mt.AOInitData {
		return playerAO(clt, true).InitData
	})
}

func initPlayer(clt *minetest.Client) {
	playerInitializedMu.Lock()
	playerInitialized[clt] = &playerData{ID: GetAOID()}
	playerInitializedMu.Unlock()

	playerInitializedMu.RLock()
	defer playerInitializedMu.RUnlock()

	add := make([]mt.AOAdd, len(playerInitialized)-1)

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

	time.Sleep(time.Second)

	// TODO fix this: send new client own thing
	//ack, _ := clt.SendCmd(&mt.ToCltAORmAdd{
	//	Add: []mt.AOAdd{playerAO(clt, true)},
	//})
	//<-ack

	time.Sleep(time.Second)

	// send new client all data:
	if len(add) != 0 {
		clt.SendCmd(&mt.ToCltAORmAdd{
			Add: add,
		})
	}
}

func playerAO(c *minetest.Client, self bool) mt.AOAdd {
	var p mt.PlayerPos
	var aoid mt.AOID
	if !self {
		p = pos.GetPos(c)
		aoid = playerInitialized[c].ID
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

			Pos: p.Pos(),
			Rot: mt.Vec{0, p.Yaw()},

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
						Rot: mt.Vec{-p.Pitch(), 0, 0},
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
