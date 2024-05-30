package main

import (
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

	player := initPlayer()

	bounding_box := rl.NewBoundingBox(rl.NewVector3(-1., 1.5, -1.), rl.NewVector3(1., 3.5, 1.))

	for !rl.WindowShouldClose() {
		manageFPS(current_monitor)

		updatePlayer(&player, bounding_box)

		rl.BeginDrawing()
		{
			rl.ClearBackground(rl.Black)
			rl.DrawFPS(10, 10)
		}

		rl.BeginMode3D(player.camera)
		{
			rl.DrawBoundingBox(bounding_box, rl.Red)

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
