package rl_first_person

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (player *Player) initCamera() {
	player.Camera.Position = rl.NewVector3(player.Position.X, player.Position.Y+(player.Scale.Y/2), player.Position.Z)
	player.Camera.Target = rl.NewVector3(
		player.Camera.Position.X-float32(math.Cos(float64(player.Rotation.X)))*float32(math.Cos(float64(player.Rotation.Y))),
		player.Camera.Position.Y+float32(math.Sin(float64(player.Rotation.Y)))+(player.Scale.Y/2),
		player.Camera.Position.Z-float32(math.Sin(float64(player.Rotation.X)))*float32(math.Cos(float64(player.Rotation.Y))),
	)
	player.Camera.Up = rl.NewVector3(0., 1., 0.)
	player.Camera.Fovy = player.Fovs.Normal
	player.Camera.Projection = rl.CameraPerspective
}

func (player *Player) updateCameraFirstPerson() {
	player.moveCamera()
	player.rotateCamera()
	player.zoomCamera()
}

func (player *Player) moveCamera() {
	player.Camera.Position = rl.NewVector3(player.Position.X, player.Position.Y+(player.Scale.Y/2), player.Position.Z)
}

func (player *Player) rotateCamera() {
	cos_rotation_y := float32(math.Cos(float64(player.Rotation.Y)))

	player.Camera.Target.X = player.Camera.Position.X - float32(math.Cos(float64(player.Rotation.X)))*cos_rotation_y
	player.Camera.Target.Y = player.Camera.Position.Y + float32(math.Sin(float64(player.Rotation.Y)))
	player.Camera.Target.Z = player.Camera.Position.Z - float32(math.Sin(float64(player.Rotation.X)))*cos_rotation_y
}

func (player *Player) zoomCamera() {
	if rl.IsKeyDown(player.Controls.Zoom) {
		player.Camera.Fovy = player.Fovs.Zoom
	} else {
		player.Camera.Fovy = player.Fovs.Normal
	}
}
