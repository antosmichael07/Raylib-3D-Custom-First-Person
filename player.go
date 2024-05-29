package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	speed             float32
	mouse_sensitivity float32
	fov               float32
	position          rl.Vector3
	rotation          rl.Vector2
	scale             rl.Vector3
	camera            rl.Camera3D
}

func initPlayer() Player {
	player := Player{}
	player.speed = .2
	player.mouse_sensitivity = .0025
	player.fov = 70.
	player.rotation = rl.NewVector2(0., 0.)
	player.position = rl.NewVector3(4., .9, 4.)
	player.scale = rl.NewVector3(.8, 1.8, .8)
	player.camera = initCamera(player)

	return player
}

func updatePlayer(player *Player) {
	movePlayer(player)
	rotatePlayer(player)
	updateCamera(player)
}

func movePlayer(player *Player) {
	speeds := rl.Vector2{
		X: float32(math.Cos(float64(player.rotation.X))) * player.speed * (rl.GetFrameTime() * 60),
		Y: float32(math.Sin(float64(player.rotation.X))) * player.speed * (rl.GetFrameTime() * 60),
	}

	if rl.IsKeyDown(rl.KeyW) {
		player.position.X -= speeds.X
		player.position.Z -= speeds.Y
	}
	if rl.IsKeyDown(rl.KeyS) {
		player.position.X += speeds.X
		player.position.Z += speeds.Y
	}
	if rl.IsKeyDown(rl.KeyA) {
		player.position.Z += speeds.X
		player.position.X -= speeds.Y
	}
	if rl.IsKeyDown(rl.KeyD) {
		player.position.Z -= speeds.X
		player.position.X += speeds.Y
	}
}

func rotatePlayer(player *Player) {
	if player.rotation.Y > 1.5 {
		player.rotation.Y = 1.5
	}
	if player.rotation.Y < -1.5 {
		player.rotation.Y = -1.5
	}

	player.rotation.X += rl.GetMouseDelta().X * player.mouse_sensitivity
	player.rotation.Y -= rl.GetMouseDelta().Y * player.mouse_sensitivity
}
