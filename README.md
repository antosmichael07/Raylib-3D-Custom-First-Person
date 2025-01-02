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

	rlfp "github.com/antosmichael07/Raylib-3D-Custom-First-Person"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	current_monitor := rl.GetCurrentMonitor()
	rl.InitWindow(int32(rl.GetMonitorWidth(current_monitor)), int32(rl.GetMonitorHeight(current_monitor)), "Raylib 3D Custom First Person - Example")
	rl.ToggleFullscreen()
	rl.DisableCursor()
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(current_monitor)))
	defer rl.CloseWindow()

	world := rlfp.World{}
	world.Init(0.)
	world.New(rl.NewVector3(0., 0., 0.), rl.NewVector2(0., 0.), false)

	world.AddBoundingBox(rl.NewBoundingBox(rl.NewVector3(2., 0., 0.), rl.NewVector3(3., .4, 1.)))
	world.AddBoundingBox(rl.NewBoundingBox(rl.NewVector3(4., .5, 0.), rl.NewVector3(5., 1., 1.)))
	world.AddBoundingBox(rl.NewBoundingBox(rl.NewVector3(6., 1., 0.), rl.NewVector3(7., 1.5, 1.)))

	world.AddTriggerBox(rl.NewBoundingBox(rl.NewVector3(-1., 0., 5.), rl.NewVector3(0., 1., 6.)))
	world.AddTriggerBox(rl.NewBoundingBox(rl.NewVector3(-1., 2., 7.), rl.NewVector3(0., 3., 8.)))

	world.AddInteractableBox(rl.NewBoundingBox(rl.NewVector3(-1., 0., -6.), rl.NewVector3(0., 1., -5.)))
	world.AddInteractableBox(rl.NewBoundingBox(rl.NewVector3(-1., 2., -8.), rl.NewVector3(0., 3., -7.)))

	for !rl.WindowShouldClose() {
		world.Update(int32(rl.GetMonitorWidth(current_monitor)), int32(rl.GetMonitorHeight(current_monitor)))

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginMode3D(world.Player.Camera)

		rl.DrawGrid(100, 1.)

		for i := range world.BoundingBoxes {
			rl.DrawBoundingBox(world.BoundingBoxes[i], rl.Red)
		}
		for i := range world.TriggerBoxes {
			rl.DrawBoundingBox(world.TriggerBoxes[i].BoundingBox, rl.Green)
		}
		for i := range world.InteractableBoxes {
			rl.DrawBoundingBox(world.InteractableBoxes[i].BoundingBox, rl.Blue)
		}

		world.DrawBoundingBoxOver()

		rl.EndMode3D()

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
