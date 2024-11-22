# Custom First Person For Raylib

Advanced library for advanced managing first person games.<hr>
Install with `go get github.com/antosmichael07/Raylib-3D-Custom-First-Person`

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

	for !rl.WindowShouldClose() {
		world.UpdatePlayer()

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

		world.DrawInteractIndicator()

		rl.EndMode3D()
		rl.DrawFPS(10, 10)

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

		rl.EndDrawing()
	}
}
```
