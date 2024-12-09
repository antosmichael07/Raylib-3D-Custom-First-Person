package rl_fp

import (
	"github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (player *Player) InitCamera() {
	player.Camera.Position = rl.NewVector3(player.Position.X, player.Position.Y+(player.Scale.Y/2), player.Position.Z)
	player.Camera.Target = rl.NewVector3(
		player.Camera.Position.X-math32.Cos(player.Rotation.X)*math32.Cos(player.Rotation.Y),
		player.Camera.Position.Y+math32.Sin(player.Rotation.Y)+(player.Scale.Y/2),
		player.Camera.Position.Z-math32.Sin(player.Rotation.X)*math32.Cos(player.Rotation.Y),
	)
	player.Camera.Up = rl.NewVector3(0., 1., 0.)
	player.Camera.Fovy = player.Fovs.Normal
	player.Camera.Projection = rl.CameraPerspective
}

func (player *Player) UpdateCamera() {
	player.UpdateCameraPosition()
	player.UpdateCameraRotation()
	player.UpdateCameraFOVY()
}

func (player *Player) UpdateCameraPosition() {
	player.Camera.Position = player.Position
	player.Camera.Position.Y += player.Scale.Y / 2
}

func (player *Player) UpdateCameraRotation() {
	cos_rotation_y := math32.Cos(player.Rotation.Y)

	player.Camera.Target.X = player.Camera.Position.X - math32.Cos(player.Rotation.X)*cos_rotation_y
	player.Camera.Target.Y = player.Camera.Position.Y + math32.Sin(player.Rotation.Y)
	player.Camera.Target.Z = player.Camera.Position.Z - math32.Sin(player.Rotation.X)*cos_rotation_y
}

func (player *Player) UpdateCameraFOVY() {
	if player.CurrentInputs[ControlZoom] {
		player.Camera.Fovy = player.Fovs.Zoom
	} else {
		player.Camera.Fovy = player.Fovs.Normal
	}
}
