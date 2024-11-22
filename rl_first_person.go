package rl_fp

import (
	"fmt"
	"math"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type World struct {
	Player            Player
	Ground            float32
	BoundingBoxes     []rl.BoundingBox
	TriggerBoxes      []TriggerBox
	InteractableBoxes []InteractableBox
}

type Player struct {
	Speed             Speeds
	MouseSensitivity  Sensitivities
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
	FrameTime         float32
	InteractRange     float32
	AlreadyInteracted bool
	StepHeight        float32
	Stepped           bool
	Controls          Controls
	Camera            rl.Camera3D
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
	player.Gravity = .04
	player.JumpPower = .8
	player.LastKeyPressed = -1
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
	player.InitCamera()
}

func (world *World) UpdatePlayer() {
	world.Player.FrameTime = rl.GetFrameTime()
	world.Player.LastKeyPressedPlayer()
	world.AccelerationPlayer()
	world.ApplyGravityToPlayer()
	if world.Player.Speed.Acceleration != 0. {
		world.StepPlayer()
		world.MovePlayer()
		world.CheckTriggerBoxes()
	}
	world.UpdateInteractableBoxes()
	world.Player.RotatePlayer()
	world.Player.UpdateCameraFirstPerson()
}

func (player *Player) LastKeyPressedPlayer() {
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

func (world *World) AccelerationPlayer() {
	final_speed := world.Player.Speed.Acceleration * world.Player.FrameTime * 60

	keys_down := []bool{rl.IsKeyDown(world.Player.Controls.Sprint), rl.IsKeyDown(world.Player.Controls.Crouch), rl.IsKeyDown(world.Player.Controls.Forward), rl.IsKeyDown(world.Player.Controls.Backward), rl.IsKeyDown(world.Player.Controls.Left), rl.IsKeyDown(world.Player.Controls.Right)}
	if !keys_down[2] && !keys_down[3] && !keys_down[4] && !keys_down[5] {
		if world.Player.Speed.Current > 0. {
			world.Player.Speed.Current -= final_speed
		} else {
			world.Player.Speed.Current = 0.
		}
	} else if (!keys_down[0] || !world.CheckIfPlayerOnSurface(&world.BoundingBoxes)) && world.Player.Speed.Current > world.Player.Speed.Normal {
		world.Player.Speed.Current -= final_speed
	}
	if world.Player.IsCrouching && world.Player.Speed.Current > world.Player.Speed.Sneak {
		world.Player.Speed.Current -= final_speed
	}

	if world.Player.Speed.Current <= world.Player.Speed.Normal && (keys_down[2] || keys_down[3] || keys_down[4] || keys_down[5]) && (!keys_down[0] || !world.CheckIfPlayerOnSurface(&world.BoundingBoxes)) && !keys_down[1] {
		world.Player.Speed.Current += final_speed
	}
	if keys_down[0] && !world.Player.IsCrouching && world.CheckIfPlayerOnSurface(&world.BoundingBoxes) && world.Player.Speed.Current <= world.Player.Speed.Sprint && (keys_down[2] || keys_down[3] || keys_down[4] || keys_down[5]) {
		world.Player.Speed.Current += final_speed
	}
	if keys_down[1] && world.Player.Speed.Current <= world.Player.Speed.Sneak && (keys_down[2] || keys_down[3] || keys_down[4] || keys_down[5]) {
		world.Player.Speed.Current += final_speed
	}
}

func (world *World) StepPlayer() {
	world.Player.Stepped = false
	player_tmp := world.Player
	player_tmp.Position.Y += world.Player.StepHeight + 0.0001
	if !player_tmp.CheckCollisionsYForPlayer(&world.BoundingBoxes) && !player_tmp.CheckCollisionsForPlayer(&world.BoundingBoxes) && world.Player.CheckCollisionsForPlayer(&world.BoundingBoxes) && world.Player.YVelocity == 0. {
		player_position_after_moving := world.Player.GetPlayerPositionAfterMoving()
		world.Player.Position.Y = (world.CheckCollisionsForPlayerAsHighestPoint() + world.Player.Scale.Y/2) + 0.0001
		world.Player.Position.X = player_position_after_moving.X
		world.Player.Position.Z = player_position_after_moving.Z
		world.Player.Stepped = true
		return
	}
	collision_x, collision_z := world.Player.CheckCollisionsXZForPlayerWithY(&world.BoundingBoxes)
	tmp_collision_x, tmp_collision_z := player_tmp.CheckCollisionsXZForPlayerWithY(&world.BoundingBoxes)
	if !player_tmp.CheckCollisionsYForPlayer(&world.BoundingBoxes) && !tmp_collision_x && collision_x && world.Player.YVelocity == 0. {
		world.Player.Position.Y = (world.CheckCollisionsXForPlayerAsHighestPoint() + world.Player.Scale.Y/2) + 0.0001
		world.Player.Position.X = world.Player.GetPlayerPositionAfterMoving().X
		world.Player.Stepped = true
		return
	}
	if !player_tmp.CheckCollisionsYForPlayer(&world.BoundingBoxes) && !tmp_collision_z && collision_z && world.Player.YVelocity == 0. {
		world.Player.Position.Y = (world.CheckCollisionsZForPlayerAsHighestPoint() + world.Player.Scale.Y/2) + 0.0001
		world.Player.Position.Z = world.Player.GetPlayerPositionAfterMoving().Z
		world.Player.Stepped = true
		return
	}
}

func (world *World) CheckCollisionsForPlayerAsHighestPoint() float32 {
	player_position := world.Player.GetPlayerPositionAfterMoving()

	highest_y := float32(0.)
	for i := 0; i < len(world.BoundingBoxes); i++ {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player_position.X-world.Player.Scale.X/2, player_position.Y-world.Player.Scale.Y/2, player_position.Z-world.Player.Scale.Z/2),
			rl.NewVector3(player_position.X+world.Player.Scale.X/2, player_position.Y+world.Player.Scale.Y/2, player_position.Z+world.Player.Scale.Z/2)), world.BoundingBoxes[i]) {
			if world.BoundingBoxes[i].Max.Y > highest_y {
				if world.BoundingBoxes[i].Min.Y <= world.BoundingBoxes[i].Max.Y {
					highest_y = world.BoundingBoxes[i].Max.Y
				} else {
					highest_y = world.BoundingBoxes[i].Min.Y
				}
			}
		}
	}

	return highest_y
}

func (world *World) CheckCollisionsXForPlayerAsHighestPoint() float32 {
	player_position_after_moving := world.Player.GetPlayerPositionAfterMoving()

	var highest_y float32 = 0.
	for i := 0; i < len(world.BoundingBoxes); i++ {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player_position_after_moving.X-world.Player.Scale.X/2, player_position_after_moving.Y-world.Player.Scale.Y/2, world.Player.Position.Z-world.Player.Scale.Z/2),
			rl.NewVector3(player_position_after_moving.X+world.Player.Scale.X/2, player_position_after_moving.Y+world.Player.Scale.Y/2, world.Player.Position.Z+world.Player.Scale.Z/2)), world.BoundingBoxes[i]) {
			if world.BoundingBoxes[i].Max.Y > highest_y {
				if world.BoundingBoxes[i].Min.Y <= world.BoundingBoxes[i].Max.Y {
					highest_y = world.BoundingBoxes[i].Max.Y
				} else {
					highest_y = world.BoundingBoxes[i].Min.Y
				}
			}
		}
	}

	return highest_y
}

func (world *World) CheckCollisionsZForPlayerAsHighestPoint() float32 {
	player_position_after_moving := world.Player.GetPlayerPositionAfterMoving()

	highest_y := float32(0.)
	for i := 0; i < len(world.BoundingBoxes); i++ {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(world.Player.Position.X-world.Player.Scale.X/2, player_position_after_moving.Y-world.Player.Scale.Y/2, player_position_after_moving.Z-world.Player.Scale.Z/2),
			rl.NewVector3(world.Player.Position.X+world.Player.Scale.X/2, player_position_after_moving.Y+world.Player.Scale.Y/2, player_position_after_moving.Z+world.Player.Scale.Z/2)), world.BoundingBoxes[i]) {
			if world.BoundingBoxes[i].Max.Y > highest_y {
				if world.BoundingBoxes[i].Min.Y <= world.BoundingBoxes[i].Max.Y {
					highest_y = world.BoundingBoxes[i].Max.Y
				} else {
					highest_y = world.BoundingBoxes[i].Min.Y
				}
			}
		}
	}

	return highest_y
}

func (world *World) MovePlayer() {
	half_crouch_scale := world.Player.ConstScale.Crouch / 2

	if rl.IsKeyDown(world.Player.Controls.Crouch) {
		world.Player.Scale.Y = world.Player.ConstScale.Crouch
		if !world.Player.IsCrouching {
			world.Player.Position.Y -= half_crouch_scale
		}
		world.Player.IsCrouching = true
	} else if world.Player.CheckPlayerUncrouch(&world.BoundingBoxes) {
		world.Player.Scale.Y = world.Player.ConstScale.Normal
		if world.Player.IsCrouching {
			world.Player.Position.Y += half_crouch_scale
		}
		world.Player.IsCrouching = false
	}
	if rl.IsKeyDown(world.Player.Controls.Jump) && world.Player.YVelocity == 0. && world.CheckIfPlayerOnSurface(&world.BoundingBoxes) && !world.Player.IsCrouching {
		world.Player.YVelocity = world.Player.JumpPower
	}

	world.Player.Position.Y += world.Player.YVelocity * (world.Player.FrameTime * 60)

	if !world.Player.Stepped {
		player_position_after_moving := world.Player.GetPlayerPositionAfterMoving()

		collisions_x, collisions_z := world.CheckCollisionsXZForPlayer()
		if collisions_x && collisions_z {
			return
		} else if collisions_x {
			world.Player.Position.Z = player_position_after_moving.Z
			return
		} else if collisions_z {
			world.Player.Position.X = player_position_after_moving.X
			return
		}

		if world.Player.CheckCollisionsForPlayer(&world.BoundingBoxes) {
			world.Player.Position.X = player_position_after_moving.X
			return
		}

		world.Player.Position.X = player_position_after_moving.X
		world.Player.Position.Z = player_position_after_moving.Z
	}
}

func (player *Player) RotatePlayer() {
	mouse_delta := rl.GetMouseDelta()
	if rl.IsKeyDown(player.Controls.Zoom) {
		player.Rotation.X += mouse_delta.X * player.MouseSensitivity.Zoom
		player.Rotation.Y -= mouse_delta.Y * player.MouseSensitivity.Zoom
	} else {
		player.Rotation.X += mouse_delta.X * player.MouseSensitivity.Normal
		player.Rotation.Y -= mouse_delta.Y * player.MouseSensitivity.Normal
	}

	if player.Rotation.Y > 1.5 {
		player.Rotation.Y = 1.5
	}
	if player.Rotation.Y < -1.5 {
		player.Rotation.Y = -1.5
	}
}

func (world *World) ApplyGravityToPlayer() {
	frame_time := world.Player.FrameTime * 60

	world.Player.YVelocity -= world.Player.Gravity * frame_time

	player_y_after_falling := world.Player.Position.Y + world.Player.YVelocity*frame_time
	if world.Player.CheckCollisionsYForPlayer(&world.BoundingBoxes) {
		world.Player.YVelocity = 0.
		return
	}
	if player_y_after_falling-(world.Player.Scale.Y/2) < world.Ground {
		world.Player.YVelocity = 0.
		world.Player.Position.Y = world.Ground + world.Player.Scale.Y/2 - .1
	}
}

func (player *Player) CheckCollisionsForPlayer(bounding_boxes *[]rl.BoundingBox) bool {
	player_position_after_moving := player.GetPlayerPositionAfterMoving()

	for i := 0; i < len(*bounding_boxes); i++ {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player_position_after_moving.X-player.Scale.X/2, player_position_after_moving.Y-player.Scale.Y/2, player_position_after_moving.Z-player.Scale.Z/2),
			rl.NewVector3(player_position_after_moving.X+player.Scale.X/2, player_position_after_moving.Y+player.Scale.Y/2, player_position_after_moving.Z+player.Scale.Z/2)), (*bounding_boxes)[i]) {
			return true
		}
	}

	return false
}

func (player *Player) CheckCollisionsYForPlayer(bounding_boxes *[]rl.BoundingBox) bool {
	player_position_y_after_moving := player.Position.Y + (player.YVelocity * player.FrameTime * 60)

	for i := 0; i < len(*bounding_boxes); i++ {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.Position.X-player.Scale.X/2, player_position_y_after_moving-player.Scale.Y/2, player.Position.Z-player.Scale.Z/2),
			rl.NewVector3(player.Position.X+player.Scale.X/2, player_position_y_after_moving+player.Scale.Y/2, player.Position.Z+player.Scale.Z/2)), (*bounding_boxes)[i]) {
			return true
		}
	}

	return false
}

func (world *World) CheckCollisionsXZForPlayer() (bool, bool) {
	player_position_after_moving := world.Player.GetPlayerPositionAfterMoving()

	collision_x, collision_z := false, false

	player_position_x := player_position_after_moving.X
	for i := 0; i < len(world.BoundingBoxes); i++ {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player_position_x-world.Player.Scale.X/2, world.Player.Position.Y-world.Player.Scale.Y/2, world.Player.Position.Z-world.Player.Scale.Z/2),
			rl.NewVector3(player_position_x+world.Player.Scale.X/2, world.Player.Position.Y+world.Player.Scale.Y/2, world.Player.Position.Z+world.Player.Scale.Z/2)), world.BoundingBoxes[i]) {
			collision_x = true
		}
	}

	player_position_z := player_position_after_moving.Z
	for i := 0; i < len(world.BoundingBoxes); i++ {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(world.Player.Position.X-world.Player.Scale.X/2, world.Player.Position.Y-world.Player.Scale.Y/2, player_position_z-world.Player.Scale.Z/2),
			rl.NewVector3(world.Player.Position.X+world.Player.Scale.X/2, world.Player.Position.Y+world.Player.Scale.Y/2, player_position_z+world.Player.Scale.Z/2)), world.BoundingBoxes[i]) {
			collision_z = true
		}
	}

	return collision_x, collision_z
}

func (player *Player) CheckCollisionsXZForPlayerWithY(bounding_boxes *[]rl.BoundingBox) (bool, bool) {
	player_position_after_moving := player.GetPlayerPositionAfterMoving()

	collision_x, collision_z := false, false

	for i := 0; i < len(*bounding_boxes); i++ {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player_position_after_moving.X-player.Scale.X/2, player_position_after_moving.Y-player.Scale.Y/2, player.Position.Z-player.Scale.Z/2),
			rl.NewVector3(player_position_after_moving.X+player.Scale.X/2, player_position_after_moving.Y+player.Scale.Y/2, player.Position.Z+player.Scale.Z/2)), (*bounding_boxes)[i]) {
			collision_x = true
		}
	}

	for i := 0; i < len(*bounding_boxes); i++ {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.Position.X-player.Scale.X/2, player_position_after_moving.Y-player.Scale.Y/2, player_position_after_moving.Z-player.Scale.Z/2),
			rl.NewVector3(player.Position.X+player.Scale.X/2, player_position_after_moving.Y+player.Scale.Y/2, player_position_after_moving.Z+player.Scale.Z/2)), (*bounding_boxes)[i]) {
			collision_z = true
		}
	}

	return collision_x, collision_z
}

func (player *Player) GetPlayerPositionAfterMoving() rl.Vector3 {
	player_position := player.Position

	frame_time := player.FrameTime * 60

	current_speed := player.Speed.Current

	if player.Speed.Normal == 0. {
		player_position.Y += player.YVelocity * frame_time
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
		player_position.X -= speeds.X
		player_position.Z -= speeds.Z
	}
	if rl.IsKeyDown(player.Controls.Backward) || player.LastKeyPressed == int32(player.Controls.Backward) {
		player_position.X += speeds.X
		player_position.Z += speeds.Z
	}
	if rl.IsKeyDown(player.Controls.Left) || player.LastKeyPressed == int32(player.Controls.Left) {
		player_position.Z += speeds.X
		player_position.X -= speeds.Z
	}
	if rl.IsKeyDown(player.Controls.Right) || player.LastKeyPressed == int32(player.Controls.Right) {
		player_position.Z -= speeds.X
		player_position.X += speeds.Z
	}

	return player_position
}

func (player *Player) CheckPlayerUncrouch(bounding_boxes *[]rl.BoundingBox) bool {
	player_position_after_moving := player.GetPlayerPositionAfterMoving()

	for i := 0; i < len(*bounding_boxes); i++ {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player_position_after_moving.X-player.Scale.X/2, player.ConstScale.Normal/2-player.ConstScale.Normal/2, player_position_after_moving.Z-player.Scale.Z/2),
			rl.NewVector3(player_position_after_moving.X+player.Scale.X/2, player.ConstScale.Normal/2+player.ConstScale.Normal/2, player_position_after_moving.Z+player.Scale.Z/2)), (*bounding_boxes)[i]) {
			return false
		}
	}

	return true
}

func (world *World) CheckIfPlayerOnSurface(bounding_boxes *[]rl.BoundingBox) bool {
	player_position_y := world.Player.Position.Y - (world.Player.Gravity * world.Player.FrameTime * 60)
	if player_position_y-(world.Player.Scale.Y/2) < world.Ground {
		return true
	}

	player_position_y_after_moving := player_position_y + (world.Player.YVelocity * world.Player.FrameTime * 60)

	for i := 0; i < len(*bounding_boxes); i++ {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(world.Player.Position.X-world.Player.Scale.X/2, player_position_y_after_moving-world.Player.Scale.Y/2, world.Player.Position.Z-world.Player.Scale.Z/2),
			rl.NewVector3(world.Player.Position.X+world.Player.Scale.X/2, player_position_y_after_moving+world.Player.Scale.Y/2, world.Player.Position.Z+world.Player.Scale.Z/2)), (*bounding_boxes)[i]) {
			return true
		}
	}

	return false
}

func (world *World) CheckTriggerBoxes() {
	for i := range world.TriggerBoxes {
		if !(world.TriggerBoxes)[i].Triggering {
			(world.TriggerBoxes)[i].Triggered = rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(world.Player.Position.X-world.Player.Scale.X/2, world.Player.Position.Y-world.Player.Scale.Y/2, world.Player.Position.Z-world.Player.Scale.Z/2), rl.NewVector3(world.Player.Position.X+world.Player.Scale.X/2, world.Player.Position.Y+world.Player.Scale.Y/2, world.Player.Position.Z+world.Player.Scale.Z/2)), (world.TriggerBoxes)[i].BoundingBox)
		} else {
			(world.TriggerBoxes)[i].Triggered = false
		}
		(world.TriggerBoxes)[i].Triggering = rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(world.Player.Position.X-world.Player.Scale.X/2, world.Player.Position.Y-world.Player.Scale.Y/2, world.Player.Position.Z-world.Player.Scale.Z/2), rl.NewVector3(world.Player.Position.X+world.Player.Scale.X/2, world.Player.Position.Y+world.Player.Scale.Y/2, world.Player.Position.Z+world.Player.Scale.Z/2)), (world.TriggerBoxes)[i].BoundingBox)
	}
}

func NewTriggerBox(box rl.BoundingBox) TriggerBox {
	return TriggerBox{box, false, false}
}

func (world *World) UpdateInteractableBoxes() {
	for i := range world.InteractableBoxes {
		world.InteractableBoxes[i].RayCollision = rl.GetRayCollisionBox(rl.GetMouseRay(rl.NewVector2(float32(rl.GetMonitorWidth(rl.GetCurrentMonitor()))/2, float32(rl.GetMonitorHeight(rl.GetCurrentMonitor()))/2), world.Player.Camera), (world.InteractableBoxes)[i].BoundingBox)
	}
	world.CheckInteractableBoxes()
}

func (world *World) DrawInteractIndicator() {
	for i := range world.InteractableBoxes {
		if world.InteractableBoxes[i].Interacting {
			return
		}
	}
	text := fmt.Sprintf("Press %s to interact", strings.ToUpper(string(world.Player.Controls.Interact)))
	text_size := rl.MeasureText(text, 30)
	for i := range world.InteractableBoxes {
		if world.InteractableBoxes[i].RayCollision.Hit && world.InteractableBoxes[i].RayCollision.Distance <= world.Player.InteractRange {
			rl.DrawText(text, int32(rl.GetScreenWidth()/2)-text_size/2, int32(rl.GetScreenHeight()/2)-30, 30, rl.White)
		}
	}
}

func (world *World) CheckInteractableBoxes() {
	for i := range world.InteractableBoxes {
		if world.Player.AlreadyInteracted {
			world.InteractableBoxes[i].Interacted = false
		}
		if rl.IsKeyDown(world.Player.Controls.Interact) && (!world.Player.AlreadyInteracted || world.InteractableBoxes[i].RayCollision.Distance > world.Player.InteractRange) {
			if world.InteractableBoxes[i].RayCollision.Hit && world.InteractableBoxes[i].RayCollision.Distance <= world.Player.InteractRange {
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
		} else if !rl.IsKeyDown(world.Player.Controls.Interact) {
			world.InteractableBoxes[i].Interacting = false
			world.InteractableBoxes[i].Interacted = false
			world.Player.AlreadyInteracted = false
		}
	}
	if rl.IsKeyDown(world.Player.Controls.Interact) {
		world.Player.AlreadyInteracted = true
	} else {
		world.Player.AlreadyInteracted = false
	}
}

func NewInteractableBox(box rl.BoundingBox) InteractableBox {
	return InteractableBox{box, false, false, rl.NewRayCollision(false, 0., rl.NewVector3(0., 0., 0.), rl.NewVector3(0., 0., 0.))}
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

func (player *Player) UpdateCameraFirstPerson() {
	player.MoveCamera()
	player.RotateCamera()
	player.ZoomCamera()
}

func (player *Player) MoveCamera() {
	player.Camera.Position = rl.NewVector3(player.Position.X, player.Position.Y+(player.Scale.Y/2), player.Position.Z)
}

func (player *Player) RotateCamera() {
	cos_rotation_y := float32(math.Cos(float64(player.Rotation.Y)))

	player.Camera.Target.X = player.Camera.Position.X - float32(math.Cos(float64(player.Rotation.X)))*cos_rotation_y
	player.Camera.Target.Y = player.Camera.Position.Y + float32(math.Sin(float64(player.Rotation.Y)))
	player.Camera.Target.Z = player.Camera.Position.Z - float32(math.Sin(float64(player.Rotation.X)))*cos_rotation_y
}

func (player *Player) ZoomCamera() {
	if rl.IsKeyDown(player.Controls.Zoom) {
		player.Camera.Fovy = player.Fovs.Zoom
	} else {
		player.Camera.Fovy = player.Fovs.Normal
	}
}
