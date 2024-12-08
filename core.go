package rl_fp

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Vector2XZ struct {
	X float32
	Z float32
}

type World struct {
	Player            Player
	Ground            float32
	BoundingBoxes     []rl.BoundingBox
	TriggerBoxes      []TriggerBox
	InteractableBoxes []InteractableBox
}

func (world *World) InitWorld(ground float32) {
	world.Player.InitPlayer()
	world.Ground = ground
	world.BoundingBoxes = []rl.BoundingBox{}
	world.TriggerBoxes = []TriggerBox{}
	world.InteractableBoxes = []InteractableBox{}
}

func (world *World) Update() {
	world.UpdateVariables()
	world.Player.UpdateRotation()
	world.UpdatePlayerPositionByStepping()
	world.UpdatePlayerPosition()
	world.UpdateTriggerBoxes()
	world.UpdateInteractableBoxes()
	world.Player.UpdateCamera()
}

func (world *World) UpdateVariables() {
	world.Player.UpdateCurrentInputs()
	world.Player.UpdateFrameTime()
	world.Player.UpdateLastDirectionalKeyPressed()
	world.UpdatePlayerCurrentSpeed()
	world.UpdatePlayerYVelocity()
}

func (player *Player) UpdateCurrentInputs() {
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

func (world *World) UpdatePlayerCurrentSpeed() {
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

func (world *World) UpdatePlayerPositionByStepping() {
	world.Player.Stepped = false
	player_position_y := world.Player.Position.Y
	world.Player.Position.Y += world.Player.StepHeight + 0.0001

	tmp_not_collisions_y_next_frame := !world.CheckPlayerCollisionsYNextFrame()
	tmp_not_collisions_xyz_next_frame := !world.CheckPlayerCollisionsXYZNextFrame()
	tmp_collisions_xyz_next_frame_after_falling_x, tmp_collisions_xyz_next_frame_after_falling_y := world.CheckPlayerCollisionsXZNextFrameAfterFalling()

	world.Player.Position.Y = player_position_y

	player_position_xyz_next_frame := world.Player.GetPositionXYZNextFrame()

	if tmp_not_collisions_y_next_frame && tmp_not_collisions_xyz_next_frame && world.CheckPlayerCollisionsXYZNextFrame() && world.Player.YVelocity == 0. {
		world.Player.Position.Y = (world.GetPlayerCollisionsXZHighestPoint() + world.Player.Scale.Y/2) + 0.0001
		world.Player.Position.X = player_position_xyz_next_frame.X
		world.Player.Position.Z = player_position_xyz_next_frame.Z
		world.Player.Stepped = true
		return
	}

	collision_x, collision_z := world.CheckPlayerCollisionsXZNextFrameAfterFalling()

	if tmp_not_collisions_y_next_frame && !tmp_collisions_xyz_next_frame_after_falling_x && collision_x && world.Player.YVelocity == 0. {
		world.Player.Position.Y = (world.GetPlayerCollisionsXHighestPoint() + world.Player.Scale.Y/2) + 0.0001
		world.Player.Position.X = player_position_xyz_next_frame.X
		world.Player.Stepped = true
		return
	}
	if tmp_not_collisions_y_next_frame && !tmp_collisions_xyz_next_frame_after_falling_y && collision_z && world.Player.YVelocity == 0. {
		world.Player.Position.Y = (world.GetPlayerCollisionsZHighestPoint() + world.Player.Scale.Y/2) + 0.0001
		world.Player.Position.Z = player_position_xyz_next_frame.Z
		world.Player.Stepped = true
		return
	}
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
