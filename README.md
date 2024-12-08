# Custom First Person For Raylib

Advanced library for advanced managing advanced first person in advanced games.

## Installation

```
go get -u github.com/antosmichael07/Raylib-3D-Custom-First-Person
```

## Example

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
	rl.SetExitKey(-1)
	defer rl.CloseWindow()

	world := rl_fp.World{}
	world.InitWorld(0.)

	world.BoundingBoxes = []rl.BoundingBox{
		rl.NewBoundingBox(rl.NewVector3(20, 0, 0), rl.NewVector3(30, 5, 10)),
		rl.NewBoundingBox(rl.NewVector3(40, 5, 0), rl.NewVector3(50, 15, 10)),
		rl.NewBoundingBox(rl.NewVector3(60, 10, 0), rl.NewVector3(70, 20, 10)),
	}
	world.TriggerBoxes = []rl_fp.TriggerBox{
		rl_fp.NewTriggerBox(rl.NewBoundingBox(rl.NewVector3(-10, 0, 50), rl.NewVector3(0, 10, 60))),
		rl_fp.NewTriggerBox(rl.NewBoundingBox(rl.NewVector3(-10, 20, 70), rl.NewVector3(0, 30, 80))),
	}
	world.InteractableBoxes = []rl_fp.InteractableBox{
		rl_fp.NewInteractableBox(rl.NewBoundingBox(rl.NewVector3(-10, 0, -60), rl.NewVector3(0, 10, -50))),
		rl_fp.NewInteractableBox(rl.NewBoundingBox(rl.NewVector3(-10, 20, -80), rl.NewVector3(0, 30, -70))),
	}

	for !rl.WindowShouldClose() {
		world.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginMode3D(world.Player.Camera)

		rl.DrawGrid(100, 10.)

		for i := range world.BoundingBoxes {
			rl.DrawBoundingBox(world.BoundingBoxes[i], rl.Red)
		}
		for i := range world.TriggerBoxes {
			rl.DrawBoundingBox(world.TriggerBoxes[i].BoundingBox, rl.Green)
		}
		for i := range world.InteractableBoxes {
			rl.DrawBoundingBox(world.InteractableBoxes[i].BoundingBox, rl.Blue)
		}

		rl.EndMode3D()

		world.DrawInteractIndicator()

		for i := range world.TriggerBoxes {
			if world.TriggerBoxes[i].Triggered {
				fmt.Printf("Triggered %d\n", i)
			}
			if world.TriggerBoxes[i].Triggering {
				rl.DrawText(fmt.Sprintf("Triggering %d", i), 10, 30, 20, rl.White)
			}
		}
		for i := range world.InteractableBoxes {
			if world.InteractableBoxes[i].Interacted {
				fmt.Printf("Interacted %d\n", i)
			}
			if world.InteractableBoxes[i].Interacting {
				rl.DrawText(fmt.Sprintf("Interacting %d", i), 10, 50, 20, rl.White)
			}
		}

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}
```
