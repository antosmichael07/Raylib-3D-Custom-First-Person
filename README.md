# Custom First Person For Raylib

Install with `go get github.com/antosmichael07/Raylib-3D-Custom-First-Person`

## Example

```go
package main

import (
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

	collision_boxes_with_player := []rl.BoundingBox{
		rl.NewBoundingBox(rl.NewVector3(-.5, 0., -.5), rl.NewVector3(.5, 1., .5)),
		rl.NewBoundingBox(rl.NewVector3(-2.5, .5, -.5), rl.NewVector3(-1.5, 1.5, .5)),
		rl.NewBoundingBox(rl.NewVector3(-4.5, 1., -.5), rl.NewVector3(-3.5, 2., .5)),
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		rl.BeginMode3D(player.Camera)

		player.UpdatePlayer(collision_boxes_with_player)

		rl.DrawGrid(20, 1.)
		for i := range collision_boxes_with_player {
			rl.DrawBoundingBox(collision_boxes_with_player[i], rl.Red)
		}

		rl.EndMode3D()

		rl.EndDrawing()
	}
}
```
