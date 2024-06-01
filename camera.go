package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (player *Player) initCamera() {
	player.camera.Position = rl.NewVector3(player.position.X, player.position.Y+(player.scale.Y/2), player.position.Z)
	player.camera.Target = rl.NewVector3(
		player.camera.Position.X-float32(math.Cos(float64(player.rotation.X)))*float32(math.Cos(float64(player.rotation.Y))),
		player.camera.Position.Y+float32(math.Sin(float64(player.rotation.Y)))+(player.scale.Y/2),
		player.camera.Position.Z-float32(math.Sin(float64(player.rotation.X)))*float32(math.Cos(float64(player.rotation.Y))),
	)
	player.camera.Up = rl.NewVector3(0., 1., 0.)
	player.camera.Fovy = player.fovs.normal
	player.camera.Projection = rl.CameraPerspective
}

func (player *Player) updateCameraFirstPerson() {
	player.moveCamera()
	player.rotateCamera()
	player.zoomCamera()
}

func (player *Player) moveCamera() {
	player.camera.Position = rl.NewVector3(player.position.X, player.position.Y+(player.scale.Y/2), player.position.Z)
}

func (player *Player) rotateCamera() {
	cos_rotation_y := float32(math.Cos(float64(player.rotation.Y)))

	player.camera.Target.X = player.camera.Position.X - float32(math.Cos(float64(player.rotation.X)))*cos_rotation_y
	player.camera.Target.Y = player.camera.Position.Y + float32(math.Sin(float64(player.rotation.Y)))
	player.camera.Target.Z = player.camera.Position.Z - float32(math.Sin(float64(player.rotation.X)))*cos_rotation_y
}

func (player *Player) zoomCamera() {
	if rl.IsKeyDown(rl.KeyC) {
		player.camera.Fovy = player.fovs.zoom
	} else {
		player.camera.Fovy = player.fovs.normal
	}
}
