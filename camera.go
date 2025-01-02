package rlfp

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chewxy/math32"
)

// Initializes all of the camera's variables
func (player *Player) InitCamera() {
	player.Camera.Position = rl.Vector3{
		X: player.Position.X,
		Y: player.BoundingBox.Max.Y - .2,
		Z: player.Position.Z,
	}
	player.Camera.Target = rl.Vector3{
		X: player.Camera.Position.X - math32.Cos(player.Rotation.X)*math32.Cos(player.Rotation.Y),
		Y: player.Camera.Position.Y + math32.Sin(player.Rotation.Y) + (player.Scale.Y / 2),
		Z: player.Camera.Position.Z - math32.Sin(player.Rotation.X)*math32.Cos(player.Rotation.Y),
	}
	player.Camera.Up = rl.Vector3{X: 0., Y: 1., Z: 0.}
	player.Camera.Fovy = player.Fovs.Normal
	player.Camera.Projection = rl.CameraPerspective
}

// Updates camera's position, rotation, fovy
func (player *Player) UpdateCamera() {
	player.UpdateCameraPosition()
	player.UpdateCameraRotation()
	player.UpdateCameraFOVY()
}

// Updates player.Camera.Position to player.Position
//
// player.Camera.Position.Y is set to player.BoundingBox.Max.Y - .2
func (player *Player) UpdateCameraPosition() {
	player.Camera.Position.X = player.Position.X
	player.Camera.Position.Y = player.BoundingBox.Max.Y - .2
	player.Camera.Position.Z = player.Position.Z
}

// Calculates player.Camera.Target position with player.Rotation
func (player *Player) UpdateCameraRotation() {
	cos_rotation_y := math32.Cos(player.Rotation.Y)

	player.Camera.Target.X = player.Camera.Position.X - math32.Cos(player.Rotation.X)*cos_rotation_y
	player.Camera.Target.Y = player.Camera.Position.Y + math32.Sin(player.Rotation.Y)
	player.Camera.Target.Z = player.Camera.Position.Z - math32.Sin(player.Rotation.X)*cos_rotation_y
}

// If the player presses the zoom control, then the fovy updates to player.Fovs.Zoom
func (player *Player) UpdateCameraFOVY() {
	if player.CurrentInputs[ControlZoom] {
		player.Camera.Fovy = player.Fovs.Zoom
	} else {
		player.Camera.Fovy = player.Fovs.Normal
	}
}
