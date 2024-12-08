package rl_fp

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type InteractableBox struct {
	BoundingBox  rl.BoundingBox
	Interacted   bool
	Interacting  bool
	RayCollision rl.RayCollision
}

func NewInteractableBox(box rl.BoundingBox) InteractableBox {
	return InteractableBox{box, false, false, rl.NewRayCollision(false, 0., rl.NewVector3(0., 0., 0.), rl.NewVector3(0., 0., 0.))}
}

func (world *World) UpdateInteractableBoxes() {
	mouse_ray := rl.GetMouseRay(rl.NewVector2(float32(rl.GetMonitorWidth(rl.GetCurrentMonitor()))/2, float32(rl.GetMonitorHeight(rl.GetCurrentMonitor()))/2), world.Player.Camera)
	in_distance := false

	if world.Player.AlreadyInteracted {
		for i := range world.InteractableBoxes {
			world.InteractableBoxes[i].Interacted = false
		}
	}

	for i := range world.InteractableBoxes {
		world.InteractableBoxes[i].RayCollision = rl.GetRayCollisionBox(mouse_ray, world.InteractableBoxes[i].BoundingBox)
		in_distance = world.InteractableBoxes[i].RayCollision.Distance <= world.Player.InteractRange

		if world.Player.CurrentInputs[ControlInteract] && (!world.Player.AlreadyInteracted || !in_distance) {
			if world.InteractableBoxes[i].RayCollision.Hit && in_distance {
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

	if world.Player.CurrentInputs[ControlInteract] {
		world.Player.AlreadyInteracted = true
	} else {
		world.Player.AlreadyInteracted = false
	}
}

func (world *World) DrawInteractIndicator() {
	for i := range world.InteractableBoxes {
		if world.InteractableBoxes[i].Interacting {
			return
		}
	}

	text := fmt.Sprintf("Press %s to interact", strings.ToUpper(string(world.Player.Controls.Interact)))
	for i := range world.InteractableBoxes {
		if world.InteractableBoxes[i].RayCollision.Hit && world.InteractableBoxes[i].RayCollision.Distance <= world.Player.InteractRange {
			rl.DrawText(text, int32(rl.GetScreenWidth()/2)-rl.MeasureText(text, 30)/2, int32(rl.GetScreenHeight()/2)-30, 30, rl.White)
			return
		}
	}
}
