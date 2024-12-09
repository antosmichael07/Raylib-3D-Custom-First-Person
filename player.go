package rl_fp

import (
	"github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ControlForward = iota
	ControlBackward
	ControlLeft
	ControlRight
	ControlJump
	ControlCrouch
	ControlSprint
	ControlZoom
	ControlInteract

	ControlCount
)

type Player struct {
	Speed                     Speeds
	MouseSensitivity          Sensitivities
	Fovs                      FOVs
	Position                  rl.Vector3
	Rotation                  rl.Vector2
	Scale                     rl.Vector3
	ConstScale                Scale
	IsCrouching               bool
	YVelocity                 float32
	Gravity                   float32
	JumpPower                 float32
	LastDirectionalKeyPressed int32
	FrameTime                 float32
	InteractRange             float32
	AlreadyInteracted         bool
	StepHeight                float32
	Stepped                   bool
	Controls                  [ControlCount]int32
	CurrentInputs             [ControlCount]bool
	Camera                    rl.Camera3D
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

func (player *Player) InitPlayer() {
	player.Speed.Normal = .65
	player.Speed.Sprint = 1.
	player.Speed.Sneak = .3
	player.Speed.Current = 0.
	player.Speed.Acceleration = .1
	player.MouseSensitivity.Normal = .0025
	player.MouseSensitivity.Zoom = .0005
	player.Fovs.Normal = 70.
	player.Fovs.Zoom = 20.
	player.Rotation = rl.NewVector2(0., 0.)
	player.Position = rl.NewVector3(0., 0., 0.)
	player.Scale = rl.NewVector3(8., 18., 8.)
	player.ConstScale.Normal = 18.
	player.ConstScale.Crouch = 9.
	player.IsCrouching = false
	player.YVelocity = 0.
	player.Gravity = .06
	player.JumpPower = 1.2
	player.LastDirectionalKeyPressed = -1
	player.FrameTime = 0.
	player.InteractRange = 30.
	player.AlreadyInteracted = false
	player.StepHeight = 5.
	player.Stepped = false
	player.Controls[ControlForward] = rl.KeyW
	player.Controls[ControlBackward] = rl.KeyS
	player.Controls[ControlLeft] = rl.KeyA
	player.Controls[ControlRight] = rl.KeyD
	player.Controls[ControlJump] = rl.KeySpace
	player.Controls[ControlCrouch] = rl.KeyLeftControl
	player.Controls[ControlSprint] = rl.KeyLeftShift
	player.Controls[ControlZoom] = rl.KeyC
	player.Controls[ControlInteract] = rl.KeyE
	player.CurrentInputs = [9]bool{false, false, false, false, false, false, false, false, false}
	player.InitCamera()
}

func (player *Player) GetPositionXYZNextFrame() rl.Vector3 {
	player_position := player.Position
	current_speed := player.Speed.Current

	if player.Speed.Normal == 0. {
		player_position.Y += player.YVelocity * player.FrameTime
	}

	keys_pressed := 0
	if player.CurrentInputs[ControlForward] {
		keys_pressed++
	}
	if player.CurrentInputs[ControlBackward] {
		keys_pressed++
	}
	if player.CurrentInputs[ControlLeft] {
		keys_pressed++
	}
	if player.CurrentInputs[ControlRight] {
		keys_pressed++
	}
	if keys_pressed == 2 {
		current_speed = current_speed * .707
	}

	final_speed := current_speed * player.FrameTime

	speeds := Vector2XZ{
		math32.Cos(player.Rotation.X) * final_speed,
		math32.Sin(player.Rotation.X) * final_speed,
	}

	if player.CurrentInputs[ControlForward] || player.LastDirectionalKeyPressed == player.Controls[ControlForward] {
		player_position.X -= speeds.X
		player_position.Z -= speeds.Z
	}
	if player.CurrentInputs[ControlBackward] || player.LastDirectionalKeyPressed == player.Controls[ControlBackward] {
		player_position.X += speeds.X
		player_position.Z += speeds.Z
	}
	if player.CurrentInputs[ControlLeft] || player.LastDirectionalKeyPressed == player.Controls[ControlLeft] {
		player_position.Z += speeds.X
		player_position.X -= speeds.Z
	}
	if player.CurrentInputs[ControlRight] || player.LastDirectionalKeyPressed == player.Controls[ControlRight] {
		player_position.Z -= speeds.X
		player_position.X += speeds.Z
	}

	return player_position
}

func (world *World) CanPlayerUncrouch() bool {
	player_position_next_frame := world.Player.GetPositionXYZNextFrame()
	player_bounding_box_next_frame := rl.NewBoundingBox(rl.NewVector3(player_position_next_frame.X-world.Player.Scale.X/2, world.Player.ConstScale.Normal/2-world.Player.ConstScale.Normal/2, player_position_next_frame.Z-world.Player.Scale.Z/2), rl.NewVector3(player_position_next_frame.X+world.Player.Scale.X/2, world.Player.ConstScale.Normal/2+world.Player.ConstScale.Normal/2, player_position_next_frame.Z+world.Player.Scale.Z/2))

	for i := range world.BoundingBoxes {
		if rl.CheckCollisionBoxes(player_bounding_box_next_frame, world.BoundingBoxes[i]) {
			return false
		}
	}

	return true
}

func (world *World) IsPlayerOnGroundNextFrame() bool {
	player_position_y_next_frame := world.Player.Position.Y + (world.Player.YVelocity-world.Player.Gravity*world.Player.FrameTime)*world.Player.FrameTime
	player_bounding_box_next_frame := rl.NewBoundingBox(rl.NewVector3(world.Player.Position.X-world.Player.Scale.X/2, player_position_y_next_frame-world.Player.Scale.Y/2, world.Player.Position.Z-world.Player.Scale.Z/2), rl.NewVector3(world.Player.Position.X+world.Player.Scale.X/2, player_position_y_next_frame+world.Player.Scale.Y/2, world.Player.Position.Z+world.Player.Scale.Z/2))

	if player_position_y_next_frame-(world.Player.Scale.Y/2) <= world.Ground {
		return true
	}

	for i := range world.BoundingBoxes {
		if rl.CheckCollisionBoxes(player_bounding_box_next_frame, world.BoundingBoxes[i]) {
			return true
		}
	}

	return false
}
