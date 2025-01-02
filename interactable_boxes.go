package rlfp

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Used for interacting with the world with a key press
type InteractableBox struct {
	// The bounding box of the interactable object
	BoundingBox rl.BoundingBox
	// The bounding box of the interactable object with a little extra space, used for drawing, when the player can interact with it
	BoundingBoxOver rl.BoundingBox
	// If the player has interacted with the object (one frame)
	Interacted bool
	// If the player is interacting with the object (holding the key)
	Interacting bool
	// The ray collision of the mouse ray and the interactable object
	RayCollision rl.RayCollision
}

// Creates a new interactable box and puts it in world.InteractableBoxes array
//
// #1 argument box: rl.BoundingBox - the bounding box of the interactable object (position)
func (world *World) AddInteractableBox(box rl.BoundingBox) {
	world.InteractableBoxes = append(world.InteractableBoxes, InteractableBox{
		box,
		rl.BoundingBox{
			Min: rl.Vector3{
				X: box.Min.X - .02,
				Y: box.Min.Y - .02,
				Z: box.Min.Z - .02,
			},
			Max: rl.Vector3{
				X: box.Max.X + .02,
				Y: box.Max.Y + .02,
				Z: box.Max.Z + .02,
			},
		},
		false,
		false,
		rl.RayCollision{
			Hit:      false,
			Distance: 0.,
			Point:    rl.Vector3{X: 0., Y: 0., Z: 0.},
			Normal:   rl.Vector3{X: 0., Y: 0., Z: 0.},
		},
	})
}

// Updates the interactable boxes
//
// #1 argument window_width: int32 - the width of the window
//
// #2 argument window_height: int32 - the height of the window
func (world *World) UpdateInteractableBoxes(window_width int32, window_height int32) {
	// If the player has already interacted with an object, reset the interacted state of all interactable boxes
	if !world.AlreadySetInteractStates && world.Player.AlreadyInteracted {
		for i := range world.InteractableBoxes {
			world.InteractableBoxes[i].Interacted = false
		}

		world.AlreadySetInteractStates = true
	}

	mouse_ray := rl.GetScreenToWorldRay(
		rl.Vector2{
			X: float32(window_width) / 2,
			Y: float32(window_height) / 2,
		},
		world.Player.Camera,
	)

	// Update the individual interactable boxes
	for i := range world.InteractableBoxes {
		if getDistance(world.Player.Position.X, world.Player.Position.Z,
			world.InteractableBoxes[i].BoundingBox.Min.X, world.InteractableBoxes[i].BoundingBox.Min.Z) <= world.CalculationDistance {

			world.UpdateInteractableBox(i, &mouse_ray)
		}
	}

	// If the player is not interacting with any object, reset the interacted state of all interactable boxes
	if world.Player.CurrentInputs[ControlInteract] {
		world.Player.AlreadyInteracted = true
	} else {
		world.Player.AlreadyInteracted = false
		world.AlreadySetInteractStates = false
	}
}

// Updates an interactable box
//
// #1 argument i: int - the index of the interactable box to update
//
// #2 argument mouse_ray: *rl.Ray - the current mouse ray
func (world *World) UpdateInteractableBox(i int, mouse_ray *rl.Ray) {
	world.InteractableBoxes[i].RayCollision = rl.GetRayCollisionBox(*mouse_ray, world.InteractableBoxes[i].BoundingBox)

	// Setting the states of the interactable box
	if world.Player.CurrentInputs[ControlInteract] && (!world.Player.AlreadyInteracted ||
		world.InteractableBoxes[i].RayCollision.Distance > world.Player.InteractRange) {

		if (world.InteractableBoxes[i].RayCollision.Hit && world.InteractableBoxes[i].RayCollision.Distance <= world.Player.InteractRange) ||
			world.InteractableBoxes[i].Interacting {

			world.Player.AlreadyInteracted = true
			if !world.InteractableBoxes[i].Interacting {
				world.InteractableBoxes[i].Interacted = true
			} else {
				world.InteractableBoxes[i].Interacted = false
			}
			world.InteractableBoxes[i].Interacting = true
		} else {
			world.InteractableBoxes[i].Interacting = false
			world.InteractableBoxes[i].Interacted = false
		}
	} else if !world.Player.CurrentInputs[ControlInteract] {
		world.InteractableBoxes[i].Interacting = false
		world.InteractableBoxes[i].Interacted = false
		world.Player.AlreadyInteracted = false
	}
}

// Draws boundingBoxOver of an interactable object, if the player can interact with it
func (world *World) DrawBoundingBoxOver() {
	// If the player is interacting with an object, don't draw the bounding box
	for i := range world.InteractableBoxes {
		if world.InteractableBoxes[i].Interacting {
			return
		}
	}

	// If the player can interact with an object, draw the bounding box
	for i := range world.InteractableBoxes {
		if world.InteractableBoxes[i].RayCollision.Hit && world.InteractableBoxes[i].RayCollision.Distance <= world.Player.InteractRange {
			rl.DrawBoundingBox(world.InteractableBoxes[i].BoundingBoxOver, rl.White)
			return
		}
	}
}
