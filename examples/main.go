package main

import (
	rl_fp "github.com/antosmichael07/Raylib-3D-Custom-First-Person/rl_first_person"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.SetTraceLogLevel(rl.LogError)
	current_monitor := rl.GetCurrentMonitor()

	rl.InitWindow(int32(rl.GetMonitorWidth(current_monitor)), int32(rl.GetMonitorHeight(current_monitor)), "Raylib Testing")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(current_monitor)))
	rl.ToggleFullscreen()
	rl.DisableCursor()
	rl.SetExitKey(-1)

	player := rl_fp.Player{}
	player.initPlayer()

	bounding_boxes := []rl.BoundingBox{
		rl.NewBoundingBox(rl.NewVector3(-1., 1.5, -1.), rl.NewVector3(1., 3.5, 1.)),
		rl.NewBoundingBox(rl.NewVector3(2., 0., -1.), rl.NewVector3(4., 2, 1.)),
		rl.NewBoundingBox(rl.NewVector3(5., -1.5, -1.), rl.NewVector3(7., .5, 1.)),
	}

	for !rl.WindowShouldClose() {
		manageFPS(current_monitor)

		if rl.IsKeyPressed(rl.KeyEscape) {
			rl.MinimizeWindow()
		}

		player.updatePlayer(bounding_boxes)

		rl.BeginDrawing()
		{
			rl.ClearBackground(rl.Black)
			rl.DrawFPS(10, 10)

			rl.DrawRectangle(int32(rl.GetScreenWidth())/2-4, int32(rl.GetScreenHeight())/2-4, 8, 8, rl.Fade(rl.White, .5))
		}

		rl.BeginMode3D(player.Camera)
		{
			showBoundingBoxes(bounding_boxes, rl.Red)

			rl.DrawGrid(20, 1.0)
		}
		rl.EndMode3D()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func manageFPS(monitor int) {
	if rl.IsWindowFocused() {
		rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(monitor)))
	} else {
		rl.SetTargetFPS(10)
	}
}

func showBoundingBoxes(bounding_boxes []rl.BoundingBox, color rl.Color) {
	for _, box := range bounding_boxes {
		rl.DrawBoundingBox(box, color)
	}
}