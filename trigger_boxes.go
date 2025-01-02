package rlfp

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Used when a player is inside a trigger box
type TriggerBox struct {
	// The box that triggers the event
	BoundingBox rl.BoundingBox
	// If the player is inside the box (one frame)
	Triggered bool
	// If the player is inside the box (staying inside)
	Triggering bool
}

// Creates a new trigger box and puts it in world.TriggerBoxes array
//
// #1 argument box: rl.BoundingBox - the box that triggers the event
func (world *World) AddTriggerBox(box rl.BoundingBox) {
	world.TriggerBoxes = append(world.TriggerBoxes, TriggerBox{box, false, false})
}

// Updates all trigger boxes
func (world *World) UpdateTriggerBoxes() {
	for i := range world.TriggerBoxes {
		if getDistance(world.Player.Position.X, world.Player.Position.Z,
			world.TriggerBoxes[i].BoundingBox.Min.X, world.TriggerBoxes[i].BoundingBox.Min.Z) <= world.CalculationDistance {

			world.UpdateTriggerBox(i)
		}
	}
}

// Updates a trigger box
//
// #1 argument i: int - index of the trigger box
func (world *World) UpdateTriggerBox(i int) {
	is_colliding := rl.CheckCollisionBoxes(world.Player.BoundingBox, world.TriggerBoxes[i].BoundingBox)

	if !world.TriggerBoxes[i].Triggering {
		world.TriggerBoxes[i].Triggered = is_colliding
	} else {
		world.TriggerBoxes[i].Triggered = false
	}

	world.TriggerBoxes[i].Triggering = is_colliding
}
