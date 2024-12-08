package rl_fp

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TriggerBox struct {
	BoundingBox rl.BoundingBox
	Triggered   bool
	Triggering  bool
}

func NewTriggerBox(box rl.BoundingBox) TriggerBox {
	return TriggerBox{box, false, false}
}

func (world *World) UpdateTriggerBoxes() {
	player_bounding_box := rl.NewBoundingBox(rl.NewVector3(world.Player.Position.X-world.Player.Scale.X/2, world.Player.Position.Y-world.Player.Scale.Y/2, world.Player.Position.Z-world.Player.Scale.Z/2), rl.NewVector3(world.Player.Position.X+world.Player.Scale.X/2, world.Player.Position.Y+world.Player.Scale.Y/2, world.Player.Position.Z+world.Player.Scale.Z/2))
	is_colliding := false

	for i := range world.TriggerBoxes {
		is_colliding = rl.CheckCollisionBoxes(player_bounding_box, world.TriggerBoxes[i].BoundingBox)

		if !world.TriggerBoxes[i].Triggering {
			world.TriggerBoxes[i].Triggered = is_colliding
		} else {
			world.TriggerBoxes[i].Triggered = false
		}

		world.TriggerBoxes[i].Triggering = is_colliding
	}
}
