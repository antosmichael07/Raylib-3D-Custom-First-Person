# Custom First Person For Raylib

Install with `go get github.com/antosmichael07/Raylib-3D-Custom-First-Person`

## Example

The function `player.UpdatePlayer()` has to be at the last lines, because it draws on screen `Press E to interact`, if it was as the first lines, then it would be under every other thing that you draw.

```go
package main

import (
	"fmt"

	rl_fp "github.com/antosmichael07/Raylib-3D-Custom-First-Person"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	current_monitor := rl.GetCurrentMonitor()
	rl.InitWindow(int32(rl.GetMonitorWidth(current_monitor)), int32(rl.GetMonitorHeight(current_monitor)), "Raylib 3D Custom First Person - Example")
	rl.ToggleFullscreen()
	rl.DisableCursor()
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(current_monitor)))
	defer rl.CloseWindow()

	player := rl_fp.Player{}
	player.InitPlayer()

	bounding_boxes := []rl.BoundingBox{
		rl.NewBoundingBox(rl.NewVector3(1.5, -.5, -.5), rl.NewVector3(2.5, .5, .5)),
		rl.NewBoundingBox(rl.NewVector3(-.5, 0., -.5), rl.NewVector3(.5, 1., .5)),
		rl.NewBoundingBox(rl.NewVector3(-2.5, .5, -.5), rl.NewVector3(-1.5, 1.5, .5)),
		rl.NewBoundingBox(rl.NewVector3(-4.5, 1., -.5), rl.NewVector3(-3.5, 2., .5)),
	}
	trigger_boxes := []rl_fp.TriggerBox{
		rl_fp.NewTriggerBox(rl.NewBoundingBox(rl.NewVector3(3.5, 1., -.5), rl.NewVector3(4.5, 2., .5))),
		rl_fp.NewTriggerBox(rl.NewBoundingBox(rl.NewVector3(5.5, 2.5, -.5), rl.NewVector3(6.5, 3.5, .5))),
	}
	interractable_boxes := []rl_fp.InteractableBox{
		rl_fp.NewInteractableBox(rl.NewBoundingBox(rl.NewVector3(7.5, 0., -.5), rl.NewVector3(8.5, 1., .5))),
		rl_fp.NewInteractableBox(rl.NewBoundingBox(rl.NewVector3(7.5, .5, -.5), rl.NewVector3(8.5, 1.5, .5))),
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		for i := range trigger_boxes {
			if trigger_boxes[i].Triggered {
				fmt.Printf("Triggered %d\n", i)
			}
			if trigger_boxes[i].Triggering {
				rl.DrawText(fmt.Sprintf("Triggering %d", i), 10, 30, 20, rl.White)
			}
		}
		for i := range interractable_boxes {
			if interractable_boxes[i].Interacted {
				fmt.Printf("Interacted %d\n", i)
			}
			if interractable_boxes[i].Interacting {
				rl.DrawText(fmt.Sprintf("Interacting %d", i), 10, 50, 20, rl.White)
			}
		}

		rl.BeginMode3D(player.Camera)

		rl.DrawGrid(20, 1.)
		for i := range bounding_boxes {
			rl.DrawBoundingBox(bounding_boxes[i], rl.Red)
		}
		for i := range trigger_boxes {
			rl.DrawBoundingBox(trigger_boxes[i].BoundingBox, rl.Green)
		}
		for i := range interractable_boxes {
			rl.DrawBoundingBox(interractable_boxes[i].BoundingBox, rl.Blue)
		}

		rl.EndMode3D()

		player.UpdatePlayer(bounding_boxes, trigger_boxes, interractable_boxes)
		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}
}
```
