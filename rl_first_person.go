package rl_fp

import (
	"fmt"
	"math"
	"strings"

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
)

type World struct {
	Player            Player
	Ground            float32
	BoundingBoxes     []rl.BoundingBox
	TriggerBoxes      []TriggerBox
	InteractableBoxes []InteractableBox
}

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
	Controls                  Controls
	CurrentInputs             [9]bool
	Camera                    rl.Camera3D
}

type SuitControls struct {
	AirPipe            bool
	PressureStabilizer bool
	ThermalRegulator   bool
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
	Interact int32
}

type TriggerBox struct {
	BoundingBox rl.BoundingBox
	Triggered   bool
	Triggering  bool
}

type InteractableBox struct {
	BoundingBox  rl.BoundingBox
	Interacted   bool
	Interacting  bool
	RayCollision rl.RayCollision
}

func (world *World) InitWorld(ground float32) {
	world.Player.InitPlayer()
	world.Ground = ground
	world.BoundingBoxes = []rl.BoundingBox{}
	world.TriggerBoxes = []TriggerBox{}
	world.InteractableBoxes = []InteractableBox{}
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
	player.Controls.Forward = rl.KeyW
	player.Controls.Backward = rl.KeyS
	player.Controls.Left = rl.KeyA
	player.Controls.Right = rl.KeyD
	player.Controls.Jump = rl.KeySpace
	player.Controls.Crouch = rl.KeyLeftControl
	player.Controls.Sprint = rl.KeyLeftShift
	player.Controls.Zoom = rl.KeyC
	player.Controls.Interact = rl.KeyE
	player.CurrentInputs = [9]bool{false, false, false, false, false, false, false, false, false}
	player.InitCamera()
}

func (world *World) UpdatePlayer() {
	world.Player.GetInputs()
	world.UpdateVariables()
	world.Player.UpdateRotation()
	world.UpdatePlayerPositionByStepping()
	world.UpdatePlayerPosition()
	world.UpdateTriggerBoxes()
	world.UpdateInteractableBoxes()
	world.Player.UpdateCamera()
}

func (player *Player) GetInputs() {
	player.CurrentInputs[ControlForward] = rl.IsKeyDown(player.Controls.Forward)
	player.CurrentInputs[ControlBackward] = rl.IsKeyDown(player.Controls.Backward)
	player.CurrentInputs[ControlLeft] = rl.IsKeyDown(player.Controls.Left)
	player.CurrentInputs[ControlRight] = rl.IsKeyDown(player.Controls.Right)
	player.CurrentInputs[ControlJump] = rl.IsKeyDown(player.Controls.Jump)
	player.CurrentInputs[ControlCrouch] = rl.IsKeyDown(player.Controls.Crouch)
	player.CurrentInputs[ControlSprint] = rl.IsKeyDown(player.Controls.Sprint)
	player.CurrentInputs[ControlZoom] = rl.IsKeyDown(player.Controls.Zoom)
	player.CurrentInputs[ControlInteract] = rl.IsKeyDown(player.Controls.Interact)
}

func (world *World) UpdateVariables() {
	world.Player.UpdateFrameTime()
	world.Player.UpdateLastDirectionalKeyPressed()
	world.UpdateCurrentSpeed()
	world.UpdatePlayerYVelocity()
}

func (player *Player) UpdateFrameTime() {
	player.FrameTime = rl.GetFrameTime() * 60
}

func (player *Player) UpdateLastDirectionalKeyPressed() {
	if player.CurrentInputs[ControlForward] {
		player.LastDirectionalKeyPressed = player.Controls.Forward
	}
	if player.CurrentInputs[ControlBackward] {
		player.LastDirectionalKeyPressed = player.Controls.Backward
	}
	if player.CurrentInputs[ControlLeft] {
		player.LastDirectionalKeyPressed = player.Controls.Left
	}
	if player.CurrentInputs[ControlRight] {
		player.LastDirectionalKeyPressed = player.Controls.Right
	}
}

func (world *World) UpdateCurrentSpeed() {
	final_speed := world.Player.Speed.Acceleration * world.Player.FrameTime
	is_player_on_ground_next_frame := world.IsPlayerOnGroundNextFrame()

	if !world.Player.CurrentInputs[ControlForward] && !world.Player.CurrentInputs[ControlBackward] && !world.Player.CurrentInputs[ControlLeft] && !world.Player.CurrentInputs[ControlRight] {
		if world.Player.Speed.Current > 0. {
			world.Player.Speed.Current -= final_speed
		} else {
			world.Player.Speed.Current = 0.
		}
	} else if (!world.Player.CurrentInputs[ControlSprint] || !is_player_on_ground_next_frame) && world.Player.Speed.Current > world.Player.Speed.Normal {
		world.Player.Speed.Current -= final_speed
	}
	if world.Player.IsCrouching && world.Player.Speed.Current > world.Player.Speed.Sneak {
		world.Player.Speed.Current -= final_speed
	}

	if world.Player.Speed.Current <= world.Player.Speed.Normal && (world.Player.CurrentInputs[ControlForward] || world.Player.CurrentInputs[ControlBackward] || world.Player.CurrentInputs[ControlLeft] || world.Player.CurrentInputs[ControlRight]) && (!world.Player.CurrentInputs[ControlSprint] || !is_player_on_ground_next_frame) && !world.Player.CurrentInputs[ControlCrouch] {
		world.Player.Speed.Current += final_speed
	}
	if world.Player.CurrentInputs[ControlSprint] && !world.Player.IsCrouching && is_player_on_ground_next_frame && world.Player.Speed.Current <= world.Player.Speed.Sprint && (world.Player.CurrentInputs[ControlForward] || world.Player.CurrentInputs[ControlBackward] || world.Player.CurrentInputs[ControlLeft] || world.Player.CurrentInputs[ControlRight]) {
		world.Player.Speed.Current += final_speed
	}
	if world.Player.CurrentInputs[ControlCrouch] && world.Player.Speed.Current <= world.Player.Speed.Sneak && (world.Player.CurrentInputs[ControlForward] || world.Player.CurrentInputs[ControlBackward] || world.Player.CurrentInputs[ControlLeft] || world.Player.CurrentInputs[ControlRight]) {
		world.Player.Speed.Current += final_speed
	}
}

func (world *World) UpdatePlayerPositionByStepping() {
	world.Player.Stepped = false
	world.Player.Position.Y += world.Player.StepHeight + 0.0001

	if !world.CheckPlayerCollisionsYNextFrame() && !world.CheckPlayerCollisionsXYZNextFrame() && world.CheckPlayerCollisionsXYZNextFrame() && world.Player.YVelocity == 0. {
		player_position_next_frame := world.Player.GetPositionXYZNextFrame()
		world.Player.Position.Y = (world.GetPlayerCollisionsXZHighestPoint() + world.Player.Scale.Y/2) + 0.0001
		world.Player.Position.X = player_position_next_frame.X
		world.Player.Position.Z = player_position_next_frame.Z
		world.Player.Stepped = true
		return
	}

	collision_x, collision_z := world.CheckPlayerCollisionsXZNextFrameAfterFalling()
	tmp_collision_x, tmp_collision_z := world.CheckPlayerCollisionsXZNextFrameAfterFalling()
	if !world.CheckPlayerCollisionsYNextFrame() && !tmp_collision_x && collision_x && world.Player.YVelocity == 0. {
		world.Player.Position.Y = (world.GetPlayerCollisionsXHighestPoint() + world.Player.Scale.Y/2) + 0.0001
		world.Player.Position.X = world.Player.GetPositionXYZNextFrame().X
		world.Player.Stepped = true
		return
	}
	if !world.CheckPlayerCollisionsYNextFrame() && !tmp_collision_z && collision_z && world.Player.YVelocity == 0. {
		world.Player.Position.Y = (world.GetPlayerCollisionsZHighestPoint() + world.Player.Scale.Y/2) + 0.0001
		world.Player.Position.Z = world.Player.GetPositionXYZNextFrame().Z
		world.Player.Stepped = true
		return
	}

	world.Player.Position.Y -= world.Player.StepHeight + 0.0001
}

func (world *World) GetPlayerCollisionsXZHighestPoint() float32 {
	player_position_next_frame := world.Player.GetPositionXYZNextFrame()
	player_bounding_box_next_frame := rl.NewBoundingBox(rl.NewVector3(player_position_next_frame.X-world.Player.Scale.X/2, player_position_next_frame.Y-world.Player.Scale.Y/2, player_position_next_frame.Z-world.Player.Scale.Z/2), rl.NewVector3(player_position_next_frame.X+world.Player.Scale.X/2, player_position_next_frame.Y+world.Player.Scale.Y/2, player_position_next_frame.Z+world.Player.Scale.Z/2))
	highest_y := world.Ground

	for i := range world.BoundingBoxes {
		if rl.CheckCollisionBoxes(player_bounding_box_next_frame, world.BoundingBoxes[i]) && world.BoundingBoxes[i].Max.Y > highest_y {
			highest_y = world.BoundingBoxes[i].Max.Y
		}
	}

	return highest_y
}

func (world *World) GetPlayerCollisionsXHighestPoint() float32 {
	player_position_next_frame := world.Player.GetPositionXYZNextFrame()
	player_bounding_box_next_frame := rl.NewBoundingBox(rl.NewVector3(player_position_next_frame.X-world.Player.Scale.X/2, player_position_next_frame.Y-world.Player.Scale.Y/2, world.Player.Position.Z-world.Player.Scale.Z/2), rl.NewVector3(player_position_next_frame.X+world.Player.Scale.X/2, player_position_next_frame.Y+world.Player.Scale.Y/2, world.Player.Position.Z+world.Player.Scale.Z/2))
	highest_y := world.Ground

	for i := range world.BoundingBoxes {
		if rl.CheckCollisionBoxes(player_bounding_box_next_frame, world.BoundingBoxes[i]) && world.BoundingBoxes[i].Max.Y > highest_y {
			highest_y = world.BoundingBoxes[i].Max.Y
		}
	}

	return highest_y
}

func (world *World) GetPlayerCollisionsZHighestPoint() float32 {
	player_position_next_frame := world.Player.GetPositionXYZNextFrame()
	player_bounding_box_next_frame := rl.NewBoundingBox(rl.NewVector3(world.Player.Position.X-world.Player.Scale.X/2, player_position_next_frame.Y-world.Player.Scale.Y/2, player_position_next_frame.Z-world.Player.Scale.Z/2), rl.NewVector3(world.Player.Position.X+world.Player.Scale.X/2, player_position_next_frame.Y+world.Player.Scale.Y/2, player_position_next_frame.Z+world.Player.Scale.Z/2))
	highest_y := world.Ground

	for i := range world.BoundingBoxes {
		if rl.CheckCollisionBoxes(player_bounding_box_next_frame, world.BoundingBoxes[i]) && world.BoundingBoxes[i].Max.Y > highest_y {
			highest_y = world.BoundingBoxes[i].Max.Y
		}
	}

	return highest_y
}

func (world *World) UpdatePlayerPosition() {
	half_crouch_scale := world.Player.ConstScale.Crouch / 2

	if world.Player.CurrentInputs[ControlCrouch] {
		world.Player.Scale.Y = world.Player.ConstScale.Crouch
		if !world.Player.IsCrouching {
			world.Player.Position.Y -= half_crouch_scale
		}
		world.Player.IsCrouching = true
	} else if world.CanPlayerUncrouch() {
		world.Player.Scale.Y = world.Player.ConstScale.Normal
		if world.Player.IsCrouching {
			world.Player.Position.Y += half_crouch_scale
		}
		world.Player.IsCrouching = false
	}
	if world.Player.CurrentInputs[ControlJump] && world.Player.YVelocity == 0. && world.IsPlayerOnGroundNextFrame() && !world.Player.IsCrouching {
		world.Player.YVelocity = world.Player.JumpPower
	}

	world.Player.Position.Y += world.Player.YVelocity * world.Player.FrameTime

	if !world.Player.Stepped {
		player_position_next_frame := world.Player.GetPositionXYZNextFrame()

		collisions_x, collisions_z := world.CheckPlayerCollisionsXZNextFrame()
		if collisions_x && collisions_z {
			return
		} else if collisions_x {
			world.Player.Position.Z = player_position_next_frame.Z
			return
		} else if collisions_z {
			world.Player.Position.X = player_position_next_frame.X
			return
		}

		if world.CheckPlayerCollisionsXYZNextFrame() {
			world.Player.Position.X = player_position_next_frame.X
			return
		}

		world.Player.Position.X = player_position_next_frame.X
		world.Player.Position.Z = player_position_next_frame.Z
	}
}

func (player *Player) UpdateRotation() {
	mouse_delta := rl.GetMouseDelta()
	if player.CurrentInputs[ControlZoom] {
		player.Rotation.X += mouse_delta.X * player.MouseSensitivity.Zoom
		player.Rotation.Y -= mouse_delta.Y * player.MouseSensitivity.Zoom
	} else {
		player.Rotation.X += mouse_delta.X * player.MouseSensitivity.Normal
		player.Rotation.Y -= mouse_delta.Y * player.MouseSensitivity.Normal
	}

	if player.Rotation.Y > 1.57 {
		player.Rotation.Y = 1.57
	}
	if player.Rotation.Y < -1.57 {
		player.Rotation.Y = -1.57
	}
}

func (world *World) UpdatePlayerYVelocity() {
	world.Player.YVelocity -= world.Player.Gravity * world.Player.FrameTime

	if world.CheckPlayerCollisionsYNextFrame() {
		world.Player.YVelocity = 0.
		return
	}
	if world.Player.Position.Y+world.Player.YVelocity*world.Player.FrameTime-(world.Player.Scale.Y/2) < world.Ground {
		world.Player.YVelocity = 0.
		world.Player.Position.Y = world.Ground + world.Player.Scale.Y/2
	}
}

func (world *World) CheckPlayerCollisionsXYZNextFrame() bool {
	player_position_next_frame := world.Player.GetPositionXYZNextFrame()
	player_bounding_box_next_frame := rl.NewBoundingBox(rl.NewVector3(player_position_next_frame.X-world.Player.Scale.X/2, player_position_next_frame.Y-world.Player.Scale.Y/2, player_position_next_frame.Z-world.Player.Scale.Z/2), rl.NewVector3(player_position_next_frame.X+world.Player.Scale.X/2, player_position_next_frame.Y+world.Player.Scale.Y/2, player_position_next_frame.Z+world.Player.Scale.Z/2))

	for i := range world.BoundingBoxes {
		if rl.CheckCollisionBoxes(player_bounding_box_next_frame, world.BoundingBoxes[i]) {
			return true
		}
	}

	return false
}

func (world *World) CheckPlayerCollisionsYNextFrame() bool {
	player_position_y_next_frame := world.Player.Position.Y + (world.Player.YVelocity * world.Player.FrameTime)
	player_bounding_box_next_frame := rl.NewBoundingBox(rl.NewVector3(world.Player.Position.X-world.Player.Scale.X/2, player_position_y_next_frame-world.Player.Scale.Y/2, world.Player.Position.Z-world.Player.Scale.Z/2), rl.NewVector3(world.Player.Position.X+world.Player.Scale.X/2, player_position_y_next_frame+world.Player.Scale.Y/2, world.Player.Position.Z+world.Player.Scale.Z/2))

	for i := range world.BoundingBoxes {
		if rl.CheckCollisionBoxes(player_bounding_box_next_frame, world.BoundingBoxes[i]) {
			return true
		}
	}

	return false
}

func (world *World) CheckPlayerCollisionsXZNextFrame() (bool, bool) {
	player_position_next_frame := world.Player.GetPositionXYZNextFrame()
	player_bounding_box_next_frame_x := rl.NewBoundingBox(rl.NewVector3(player_position_next_frame.X-world.Player.Scale.X/2, world.Player.Position.Y-world.Player.Scale.Y/2, world.Player.Position.Z-world.Player.Scale.Z/2), rl.NewVector3(player_position_next_frame.X+world.Player.Scale.X/2, world.Player.Position.Y+world.Player.Scale.Y/2, world.Player.Position.Z+world.Player.Scale.Z/2))
	player_bounding_box_next_frame_z := rl.NewBoundingBox(rl.NewVector3(world.Player.Position.X-world.Player.Scale.X/2, world.Player.Position.Y-world.Player.Scale.Y/2, player_position_next_frame.Z-world.Player.Scale.Z/2), rl.NewVector3(world.Player.Position.X+world.Player.Scale.X/2, world.Player.Position.Y+world.Player.Scale.Y/2, player_position_next_frame.Z+world.Player.Scale.Z/2))

	collision_x, collision_z := false, false

	for i := range world.BoundingBoxes {
		if rl.CheckCollisionBoxes(player_bounding_box_next_frame_x, world.BoundingBoxes[i]) {
			collision_x = true
		}
	}

	for i := range world.BoundingBoxes {
		if rl.CheckCollisionBoxes(player_bounding_box_next_frame_z, world.BoundingBoxes[i]) {
			collision_z = true
		}
	}

	return collision_x, collision_z
}

func (world *World) CheckPlayerCollisionsXZNextFrameAfterFalling() (bool, bool) {
	player_position_next_frame := world.Player.GetPositionXYZNextFrame()
	player_bounding_box_next_frame_x := rl.NewBoundingBox(rl.NewVector3(player_position_next_frame.X-world.Player.Scale.X/2, player_position_next_frame.Y-world.Player.Scale.Y/2, world.Player.Position.Z-world.Player.Scale.Z/2), rl.NewVector3(player_position_next_frame.X+world.Player.Scale.X/2, player_position_next_frame.Y+world.Player.Scale.Y/2, world.Player.Position.Z+world.Player.Scale.Z/2))
	player_bounding_box_next_frame_z := rl.NewBoundingBox(rl.NewVector3(world.Player.Position.X-world.Player.Scale.X/2, player_position_next_frame.Y-world.Player.Scale.Y/2, player_position_next_frame.Z-world.Player.Scale.Z/2), rl.NewVector3(world.Player.Position.X+world.Player.Scale.X/2, player_position_next_frame.Y+world.Player.Scale.Y/2, player_position_next_frame.Z+world.Player.Scale.Z/2))

	collision_x, collision_z := false, false

	for i := range world.BoundingBoxes {
		if rl.CheckCollisionBoxes(player_bounding_box_next_frame_x, world.BoundingBoxes[i]) {
			collision_x = true
		}
	}

	for i := range world.BoundingBoxes {
		if rl.CheckCollisionBoxes(player_bounding_box_next_frame_z, world.BoundingBoxes[i]) {
			collision_z = true
		}
	}

	return collision_x, collision_z
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
		float32(math.Cos(float64(player.Rotation.X))) * final_speed,
		float32(math.Sin(float64(player.Rotation.X))) * final_speed,
	}

	if player.CurrentInputs[ControlForward] || player.LastDirectionalKeyPressed == player.Controls.Forward {
		player_position.X -= speeds.X
		player_position.Z -= speeds.Z
	}
	if player.CurrentInputs[ControlBackward] || player.LastDirectionalKeyPressed == player.Controls.Backward {
		player_position.X += speeds.X
		player_position.Z += speeds.Z
	}
	if player.CurrentInputs[ControlLeft] || player.LastDirectionalKeyPressed == player.Controls.Left {
		player_position.Z += speeds.X
		player_position.X -= speeds.Z
	}
	if player.CurrentInputs[ControlRight] || player.LastDirectionalKeyPressed == player.Controls.Right {
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

func NewTriggerBox(box rl.BoundingBox) TriggerBox {
	return TriggerBox{box, false, false}
}

func (world *World) UpdateTriggerBoxes() {
	player_bounding_box := rl.NewBoundingBox(rl.NewVector3(world.Player.Position.X-world.Player.Scale.X/2, world.Player.Position.Y-world.Player.Scale.Y/2, world.Player.Position.Z-world.Player.Scale.Z/2), rl.NewVector3(world.Player.Position.X+world.Player.Scale.X/2, world.Player.Position.Y+world.Player.Scale.Y/2, world.Player.Position.Z+world.Player.Scale.Z/2))
	is_colliding := false

	for i := range world.TriggerBoxes {
		is_colliding = rl.CheckCollisionBoxes(player_bounding_box, world.TriggerBoxes[i].BoundingBox)

		if !world.TriggerBoxes[i].Triggering {
			world.TriggerBoxes[i].Triggered = is_colliding
		} else {
			world.TriggerBoxes[i].Triggered = false
		}

		world.TriggerBoxes[i].Triggering = is_colliding
	}
}

func NewInteractableBox(box rl.BoundingBox) InteractableBox {
	return InteractableBox{box, false, false, rl.NewRayCollision(false, 0., rl.NewVector3(0., 0., 0.), rl.NewVector3(0., 0., 0.))}
}

func (world *World) UpdateInteractableBoxes() {
	mouse_ray := rl.GetMouseRay(rl.NewVector2(float32(rl.GetMonitorWidth(rl.GetCurrentMonitor()))/2, float32(rl.GetMonitorHeight(rl.GetCurrentMonitor()))/2), world.Player.Camera)
	in_distance := false

	if world.Player.AlreadyInteracted {
		for i := range world.InteractableBoxes {
			world.InteractableBoxes[i].Interacted = false
		}
	}

	for i := range world.InteractableBoxes {
		world.InteractableBoxes[i].RayCollision = rl.GetRayCollisionBox(mouse_ray, world.InteractableBoxes[i].BoundingBox)
		in_distance = world.InteractableBoxes[i].RayCollision.Distance <= world.Player.InteractRange

		if world.Player.CurrentInputs[ControlInteract] && (!world.Player.AlreadyInteracted || !in_distance) {
			if world.InteractableBoxes[i].RayCollision.Hit && in_distance {
				world.Player.AlreadyInteracted = true
				if !world.InteractableBoxes[i].Interacting {
					world.InteractableBoxes[i].Interacted = true
				} else {
					world.InteractableBoxes[i].Interacted = false
				}
				world.InteractableBoxes[i].Interacting = true
			} else {
				world.InteractableBoxes[i].Interacting = false
				world.InteractableBoxes[i].Interacted = false
			}
		} else if !world.Player.CurrentInputs[ControlInteract] {
			world.InteractableBoxes[i].Interacting = false
			world.InteractableBoxes[i].Interacted = false
			world.Player.AlreadyInteracted = false
		}
	}

	if world.Player.CurrentInputs[ControlInteract] {
		world.Player.AlreadyInteracted = true
	} else {
		world.Player.AlreadyInteracted = false
	}
}

func (world *World) DrawInteractIndicator() {
	for i := range world.InteractableBoxes {
		if world.InteractableBoxes[i].Interacting {
			return
		}
	}

	text := fmt.Sprintf("Press %s to interact", strings.ToUpper(string(world.Player.Controls.Interact)))
	for i := range world.InteractableBoxes {
		if world.InteractableBoxes[i].RayCollision.Hit && world.InteractableBoxes[i].RayCollision.Distance <= world.Player.InteractRange {
			rl.DrawText(text, int32(rl.GetScreenWidth()/2)-rl.MeasureText(text, 30)/2, int32(rl.GetScreenHeight()/2)-30, 30, rl.White)
			return
		}
	}
}

func (player *Player) InitCamera() {
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
	cos_rotation_y := float32(math.Cos(float64(player.Rotation.Y)))

	player.Camera.Target.X = player.Camera.Position.X - float32(math.Cos(float64(player.Rotation.X)))*cos_rotation_y
	player.Camera.Target.Y = player.Camera.Position.Y + float32(math.Sin(float64(player.Rotation.Y)))
	player.Camera.Target.Z = player.Camera.Position.Z - float32(math.Sin(float64(player.Rotation.X)))*cos_rotation_y
}

func (player *Player) UpdateCameraFOVY() {
	if player.CurrentInputs[ControlZoom] {
		player.Camera.Fovy = player.Fovs.Zoom
	} else {
		player.Camera.Fovy = player.Fovs.Normal
	}
}
