package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func initCamera(player Player) rl.Camera3D {
	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(player.position.X, player.position.Y+(player.scale.Y/2), player.position.Z)
	camera.Target = rl.NewVector3(
		player.camera.Position.X-float32(math.Cos(float64(player.rotation.X)))*float32(math.Cos(float64(player.rotation.Y))),
		player.camera.Position.Y+float32(math.Sin(float64(player.rotation.Y)))+(player.scale.Y/2),
		player.camera.Position.Z-float32(math.Sin(float64(player.rotation.X)))*float32(math.Cos(float64(player.rotation.Y))),
	)
	camera.Up = rl.NewVector3(0., 1., 0.)
	camera.Fovy = player.fov
	camera.Projection = rl.CameraPerspective

	return camera
}

func updateCameraFirstPerson(player *Player) {
	moveCamera(player)
	rotateCamera(player)
}

func moveCamera(player *Player) {
	player.camera.Position = rl.NewVector3(player.position.X, player.position.Y+(player.scale.Y/2), player.position.Z)
}

func rotateCamera(player *Player) {
	player.camera.Target.X = player.camera.Position.X - float32(math.Cos(float64(player.rotation.X)))*float32(math.Cos(float64(player.rotation.Y)))
	player.camera.Target.Y = player.camera.Position.Y + float32(math.Sin(float64(player.rotation.Y)))
	player.camera.Target.Z = player.camera.Position.Z - float32(math.Sin(float64(player.rotation.X)))*float32(math.Cos(float64(player.rotation.Y)))
}
