package ao

import (
	"github.com/ev2-1/minetest-go"
	"github.com/anon55555/mt"
	"sync"
)

var activeObjects map[mt.AOID]*ActiveObject
var activeObjectsMu sync.RWMutex

// ActiveObject describes a ActiveObject fully
type ActiveObject struct {
	ID mt.AOID
	
	Anim mt.AOAnim
	AnimSpeed float32
	Attach mt.AOAttach
	Bones map[string]mt.AOBonePos
	Props mt.AOProps
	TextureMod mt.Texture
	
	ArmorGroups []mt.Group

	Controller AOController
}

// AOController is used to controll a activeobject
// it is meant in a way, so you describe active objects as a type which implements this interface
type AOController interface {
	Interact(ActiveObject, AOInteract)
}

// AOInteract describes a interaction with a active object
type AOInteract struct {
	Player minetest.Client

	Action mt.Interaction
	ItemSlot uint16
	Pos mt.PlayerPos
}
