package rl_first_person

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Speed             Speeds
	Mouse_sensitivity Sensitivities
	Fovs              FOVs
	Position          rl.Vector3
	Rotation          rl.Vector2
	Scale             rl.Vector3
	ConstScale        Scale
	IsCrouching       bool
	YVelocity         float32
	Gravity           float32
	JumpPower         float32
	LastKeyPressed    int32
	Controls          Controls
	Camera            rl.Camera3D
}

type Speeds struct {
	Normal       float32
	Sprint       float32
	Sneak        float32
	Current      float32
	Acceleration float32
}

type Sensitivities struct {
	Normal float32
	Zoom   float32
}

type FOVs struct {
	Normal float32
	Zoom   float32
}

type Scale struct {
	Normal float32
	Crouch float32
}

type Vector2XZ struct {
	X float32
	Z float32
}

type Controls struct {
	Forward  int32
	Backward int32
	Left     int32
	Right    int32
	Jump     int32
	Crouch   int32
	Sprint   int32
	Zoom     int32
}

func (player *Player) initPlayer() {
	player.Speed.Normal = .1
	player.Speed.Sprint = .15
	player.Speed.Sneak = .05
	player.Speed.Current = 0.
	player.Speed.Acceleration = .01
	player.Mouse_sensitivity.Normal = .0025
	player.Mouse_sensitivity.Zoom = .0005
	player.Fovs.Normal = 70.
	player.Fovs.Zoom = 20.
	player.Rotation = rl.NewVector2(0., 0.)
	player.Position = rl.NewVector3(4., .9, 4.)
	player.Scale = rl.NewVector3(.8, 1.8, .8)
	player.ConstScale.Normal = 1.8
	player.ConstScale.Crouch = .9
	player.IsCrouching = false
	player.YVelocity = 0.
	player.Gravity = .0065
	player.JumpPower = .15
	player.LastKeyPressed = -1
	player.Controls.Forward = int32(rl.KeyW)
	player.Controls.Backward = int32(rl.KeyS)
	player.Controls.Left = int32(rl.KeyA)
	player.Controls.Right = int32(rl.KeyD)
	player.Controls.Jump = int32(rl.KeySpace)
	player.Controls.Crouch = int32(rl.KeyLeftControl)
	player.Controls.Sprint = int32(rl.KeyLeftShift)
	player.Controls.Zoom = int32(rl.KeyC)
	player.initCamera()
}

func (player *Player) updatePlayer(bounding_boxes []rl.BoundingBox) {
	player.lastKeyPressedPlayer()
	player.accelerationPlayer()
	player.movePlayer(bounding_boxes)
	player.rotatePlayer()
	player.applyGravityToPlayer(bounding_boxes)
	player.updateCameraFirstPerson()
}

func (player *Player) lastKeyPressedPlayer() {
	if rl.IsKeyDown(player.Controls.Forward) {
		player.LastKeyPressed = int32(player.Controls.Forward)
	}
	if rl.IsKeyDown(player.Controls.Backward) {
		player.LastKeyPressed = int32(player.Controls.Backward)
	}
	if rl.IsKeyDown(player.Controls.Left) {
		player.LastKeyPressed = int32(player.Controls.Left)
	}
	if rl.IsKeyDown(player.Controls.Right) {
		player.LastKeyPressed = int32(player.Controls.Right)
	}
}

func (player *Player) accelerationPlayer() {
	final_speed := player.Speed.Acceleration * rl.GetFrameTime() * 60

	keys_down := map[string]bool{"shift": rl.IsKeyDown(player.Controls.Sprint), "ctrl": rl.IsKeyDown(player.Controls.Crouch), "w": rl.IsKeyDown(player.Controls.Forward), "s": rl.IsKeyDown(player.Controls.Backward), "a": rl.IsKeyDown(player.Controls.Left), "d": rl.IsKeyDown(player.Controls.Right)}
	if !keys_down["w"] && !keys_down["s"] && !keys_down["a"] && !keys_down["d"] {
		if player.Speed.Current > 0. {
			player.Speed.Current -= final_speed
		} else {
			player.Speed.Current = 0.
		}
	} else if !keys_down["shift"] && player.Speed.Current > player.Speed.Normal {
		player.Speed.Current -= final_speed
	}
	if player.IsCrouching && player.Speed.Current > player.Speed.Sneak {
		player.Speed.Current -= final_speed
	}

	if player.Speed.Current <= player.Speed.Normal && (keys_down["w"] || keys_down["s"] || keys_down["a"] || keys_down["d"]) && !keys_down["shift"] && !keys_down["ctrl"] {
		player.Speed.Current += final_speed
	}
	if keys_down["shift"] && player.Speed.Current <= player.Speed.Sprint && (keys_down["w"] || keys_down["s"] || keys_down["a"] || keys_down["d"]) {
		player.Speed.Current += final_speed
	}
	if keys_down["ctrl"] && player.Speed.Current <= player.Speed.Sneak && (keys_down["w"] || keys_down["s"] || keys_down["a"] || keys_down["d"]) {
		player.Speed.Current += final_speed
	}
}

func (player *Player) movePlayer(bounding_boxes []rl.BoundingBox) {
	half_crouch_scale := player.ConstScale.Crouch / 2

	if rl.IsKeyDown(player.Controls.Crouch) {
		player.Scale.Y = player.ConstScale.Crouch
		if !player.IsCrouching {
			player.Position.Y -= half_crouch_scale
		}
		player.IsCrouching = true
	} else if player.checkPlayerUncrouch(bounding_boxes) {
		player.Scale.Y = player.ConstScale.Normal
		if player.IsCrouching {
			player.Position.Y += half_crouch_scale
		}
		player.IsCrouching = false
	}
	if rl.IsKeyDown(player.Controls.Jump) && player.YVelocity == 0. && player.checkIfPlayerOnSurface(bounding_boxes) && !player.IsCrouching {
		player.YVelocity = player.JumpPower
	}

	player.Position.Y += player.YVelocity * (rl.GetFrameTime() * 60)

	player_position_after_moving := player.getPlayerPositionAfterMoving()

	collisions_x, collisions_z := player.checkCollisionsXZForPlayer(bounding_boxes)
	if collisions_x && collisions_z {
		return
	} else if collisions_x {
		player.Position.Z = player_position_after_moving.Z
		return
	} else if collisions_z {
		player.Position.X = player_position_after_moving.X
		return
	}

	player.Position = player_position_after_moving
}

func (player *Player) rotatePlayer() {
	mouse_delta := rl.GetMouseDelta()
	if rl.IsKeyDown(player.Controls.Zoom) {
		player.Rotation.X += mouse_delta.X * player.Mouse_sensitivity.Zoom
		player.Rotation.Y -= mouse_delta.Y * player.Mouse_sensitivity.Zoom
	} else {
		player.Rotation.X += mouse_delta.X * player.Mouse_sensitivity.Normal
		player.Rotation.Y -= mouse_delta.Y * player.Mouse_sensitivity.Normal
	}

	if player.Rotation.Y > 1.5 {
		player.Rotation.Y = 1.5
	}
	if player.Rotation.Y < -1.5 {
		player.Rotation.Y = -1.5
	}
}

func (player *Player) applyGravityToPlayer(bounding_boxes []rl.BoundingBox) {
	player.YVelocity -= player.Gravity * (rl.GetFrameTime() * 60)

	if player.checkCollisionsYForPlayer(bounding_boxes) || player.getPlayerPositionAfterMoving().Y-(player.Scale.Y/2) < 0. {
		player.YVelocity = 0.
		return
	}
}

func (player Player) checkCollisionsForPlayer(bounding_boxes []rl.BoundingBox) bool {
	player.Position = player.getPlayerPositionAfterMoving()

	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.Position.X-player.Scale.X/2, player.Position.Y-player.Scale.Y/2, player.Position.Z-player.Scale.Z/2),
			rl.NewVector3(player.Position.X+player.Scale.X/2, player.Position.Y+player.Scale.Y/2, player.Position.Z+player.Scale.Z/2)), box) {
			return true
		}
	}

	return false
}

func (player Player) checkCollisionsYForPlayer(bounding_boxes []rl.BoundingBox) bool {
	player.Speed.Normal = 0
	player.Speed.Sprint = 0
	player.Speed.Sneak = 0
	player.Speed.Current = 0

	return player.checkCollisionsForPlayer(bounding_boxes)
}

func (player Player) checkCollisionsXZForPlayer(bounding_boxes []rl.BoundingBox) (bool, bool) {
	player_position_after_moving := player.getPlayerPositionAfterMoving()

	collision_x, collision_z := false, false

	player_position_x := player_position_after_moving.X
	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player_position_x-player.Scale.X/2, player.Position.Y-player.Scale.Y/2, player.Position.Z-player.Scale.Z/2),
			rl.NewVector3(player_position_x+player.Scale.X/2, player.Position.Y+player.Scale.Y/2, player.Position.Z+player.Scale.Z/2)), box) {
			collision_x = true
		}
	}

	player_position_z := player_position_after_moving.Z
	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.Position.X-player.Scale.X/2, player.Position.Y-player.Scale.Y/2, player_position_z-player.Scale.Z/2),
			rl.NewVector3(player.Position.X+player.Scale.X/2, player.Position.Y+player.Scale.Y/2, player_position_z+player.Scale.Z/2)), box) {
			collision_z = true
		}
	}

	return collision_x, collision_z
}

func (player Player) getPlayerPositionAfterMoving() rl.Vector3 {
	frame_time := rl.GetFrameTime() * 60

	current_speed := player.Speed.Current

	if player.Speed.Normal == 0. {
		player.Position.Y += player.YVelocity * frame_time
	}

	keys_pressed := 0
	if rl.IsKeyDown(player.Controls.Forward) {
		keys_pressed++
	}
	if rl.IsKeyDown(player.Controls.Backward) {
		keys_pressed++
	}
	if rl.IsKeyDown(player.Controls.Left) {
		keys_pressed++
	}
	if rl.IsKeyDown(player.Controls.Right) {
		keys_pressed++
	}
	if keys_pressed == 2 {
		current_speed = current_speed * .707
	}

	final_speed := current_speed * frame_time

	speeds := Vector2XZ{
		float32(math.Cos(float64(player.Rotation.X))) * final_speed,
		float32(math.Sin(float64(player.Rotation.X))) * final_speed,
	}

	if rl.IsKeyDown(player.Controls.Forward) || player.LastKeyPressed == int32(player.Controls.Forward) {
		player.Position.X -= speeds.X
		player.Position.Z -= speeds.Z
	}
	if rl.IsKeyDown(player.Controls.Backward) || player.LastKeyPressed == int32(player.Controls.Backward) {
		player.Position.X += speeds.X
		player.Position.Z += speeds.Z
	}
	if rl.IsKeyDown(player.Controls.Left) || player.LastKeyPressed == int32(player.Controls.Left) {
		player.Position.Z += speeds.X
		player.Position.X -= speeds.Z
	}
	if rl.IsKeyDown(player.Controls.Right) || player.LastKeyPressed == int32(player.Controls.Right) {
		player.Position.Z -= speeds.X
		player.Position.X += speeds.Z
	}

	return player.Position
}

func (player Player) checkPlayerUncrouch(bounding_boxes []rl.BoundingBox) bool {
	player.Scale.Y = player.ConstScale.Normal
	player.Position.Y += player.ConstScale.Normal / 2

	return !player.checkCollisionsForPlayer(bounding_boxes)
}

func (player Player) checkIfPlayerOnSurface(bounding_boxes []rl.BoundingBox) bool {
	player.Position.Y -= player.Gravity * (rl.GetFrameTime() * 60)
	if player.checkCollisionsYForPlayer(bounding_boxes) || player.Position.Y-(player.Scale.Y/2) < 0. {
		return true
	}
	return false
}

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
