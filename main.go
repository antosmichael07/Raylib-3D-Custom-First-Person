package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.SetTraceLogLevel(rl.LogError)
	current_monitor := rl.GetCurrentMonitor()

	rl.InitWindow(int32(rl.GetMonitorWidth(current_monitor)), int32(rl.GetMonitorHeight(current_monitor)), "Test Raylib Go")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(current_monitor)))
	rl.ToggleFullscreen()
	rl.DisableCursor()
	rl.SetExitKey(-1)

	player := initPlayer()

	for !rl.WindowShouldClose() {
		updatePlayer(&player)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		{
			rl.DrawFPS(10, 10)
		}

		rl.BeginMode3D(player.camera)
		{
			rl.DrawGrid(20, 1.0)
		}
		rl.EndMode3D()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
