package rl_fp

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

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
