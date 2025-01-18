package rlfp

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Checks if the player is colliding with a bounding box after moving by world.Player.OffsetNextFrame.X on the X axis
//
// #1 return: int - the index of the bounding box that is colliding with the player (if there is no collision, returns -1)
//
// #2 return: bool - true if world.Player.OffsetNextFrame.X is positive
func (world *World) checkPlayerCollisionsXNextFrame() (int, bool) {
	bounding_box := world.Player.BoundingBox

	// Move bounding_box by world.Player.OffsetNextFrame.X on the X axis
	if world.Player.OffsetNextFrame.X > 0 {
		bounding_box.Max.X += world.Player.OffsetNextFrame.X
	} else {
		bounding_box.Min.X += world.Player.OffsetNextFrame.X
	}

	// Check if bounding_box is colliding with another bounding box
	for i := range world.BoundingBoxes {
		if getDistance(world.Player.Position.X, world.Player.Position.Z,
			world.BoundingBoxes[i].Min.X, world.BoundingBoxes[i].Min.Z) <= world.CalculationDistance &&
			rl.CheckCollisionBoxes(bounding_box, world.BoundingBoxes[i]) {

			return i, world.Player.OffsetNextFrame.X > 0
		}
	}

	return -1, false
}

// Checks if the player is colliding with a bounding box after moving by world.Player.OffsetNextFrame.X on the X axis
//
// #1 argument y: float32 - the value to set to the Y axis of the  player's bounding box
//
// #1 return: bool - if there is a collision with a bounding box
func (world *World) checkPlayerCollisionsXYNextFrame(y float32) bool {
	bounding_box := world.Player.BoundingBox

	// Move bounding_box by world.Player.OffsetNextFrame.X on the X axis
	if world.Player.OffsetNextFrame.X > 0 {
		bounding_box.Max.X += world.Player.OffsetNextFrame.X
	} else {
		bounding_box.Min.X += world.Player.OffsetNextFrame.X - world.FrameTime
	}

	bounding_box.Min.Y = y
	bounding_box.Max.Y = y + world.Player.Scale.Y

	// Check if bounding_box is colliding with another bounding box
	for i := range world.BoundingBoxes {
		if getDistance(world.Player.Position.X, world.Player.Position.Z,
			world.BoundingBoxes[i].Min.X, world.BoundingBoxes[i].Min.Z) <= world.CalculationDistance &&
			rl.CheckCollisionBoxes(bounding_box, world.BoundingBoxes[i]) {

			return true
		}
	}

	return false
}

// Checks if the player is colliding with a bounding box after moving by world.Player.OffsetNextFrame.Y on the Y axis
//
// #1 return: int - the index of the bounding box that is colliding with the player (if there is no collision, returns -1)
//
// #2 return: bool - true if world.Player.OffsetNextFrame.Y is positive
func (world *World) checkPlayerCollisionsYNextFrame() (int, bool) {
	bounding_box := world.Player.BoundingBox

	// Move bounding_box by world.Player.OffsetNextFrame.Y on the Y axis
	if world.Player.OffsetNextFrame.Y > 0 {
		bounding_box.Max.Y += world.Player.OffsetNextFrame.Y
	} else {
		bounding_box.Min.Y += world.Player.OffsetNextFrame.Y
	}

	// Check if bounding_box is colliding with another bounding box
	for i := range world.BoundingBoxes {
		if getDistance(world.Player.Position.X, world.Player.Position.Z,
			world.BoundingBoxes[i].Min.X, world.BoundingBoxes[i].Min.Z) <= world.CalculationDistance &&
			rl.CheckCollisionBoxes(bounding_box, world.BoundingBoxes[i]) {

			return i, world.Player.OffsetNextFrame.Y > 0
		}
	}

	return -1, false
}

// Checks if the player is colliding with a bounding box after moving to world.Ground on the Y axis
//
// #1 return: int - the index of the bounding box that is colliding with the player (if there is no collision, returns -1)
func (world *World) checkPlayerCollisionsYOnGround() int {
	bounding_box := world.Player.BoundingBox

	// Move bounding_box to world.Ground on the Y axis
	bounding_box.Max.Y += world.Ground + world.FloatPrecision - world.Player.BoundingBox.Min.Y

	// Check if bounding_box is colliding with another bounding box
	for i := range world.BoundingBoxes {
		if getDistance(world.Player.Position.X, world.Player.Position.Z,
			world.BoundingBoxes[i].Min.X, world.BoundingBoxes[i].Min.Z) <= world.CalculationDistance &&
			rl.CheckCollisionBoxes(bounding_box, world.BoundingBoxes[i]) {

			return i
		}
	}

	return -1
}

// Checks if the player is colliding with a bounding box after moving by world.Player.OffsetNextFrame.Z on the Z axis
//
// #1 return: int - the index of the bounding box that is colliding with the player (if there is no collision, returns -1)
//
// #2 return: bool - true if world.Player.OffsetNextFrame.Z is positive
func (world *World) checkPlayerCollisionsZNextFrame() (int, bool) {
	bounding_box := world.Player.BoundingBox

	// Move bounding_box by world.Player.OffsetNextFrame.Z on the Z axis
	if world.Player.OffsetNextFrame.Z > 0 {
		bounding_box.Max.Z += world.Player.OffsetNextFrame.Z
	} else {
		bounding_box.Min.Z += world.Player.OffsetNextFrame.Z
	}

	// Check if bounding_box is colliding with another bounding box
	for i := range world.BoundingBoxes {
		if getDistance(world.Player.Position.X, world.Player.Position.Z,
			world.BoundingBoxes[i].Min.X, world.BoundingBoxes[i].Min.Z) <= world.CalculationDistance &&
			rl.CheckCollisionBoxes(bounding_box, world.BoundingBoxes[i]) {

			return i, world.Player.OffsetNextFrame.Z > 0
		}
	}

	return -1, false
}

// Checks if the player is colliding with a bounding box after moving by world.Player.OffsetNextFrame.Z on the Z axis
//
// #1 argument y: float32 - the value to set to the Y axis of the player's bounding box
//
// #1 return: bool - if there is a collision with a bounding box
func (world *World) checkPlayerCollisionsZYNextFrame(y float32) bool {
	bounding_box := world.Player.BoundingBox

	// Move bounding_box by world.Player.OffsetNextFrame.Z on the Z axis
	if world.Player.OffsetNextFrame.Z > 0 {
		bounding_box.Max.Z += world.Player.OffsetNextFrame.Z
	} else {
		bounding_box.Min.Z += world.Player.OffsetNextFrame.Z
	}

	bounding_box.Min.Y = y
	bounding_box.Max.Y = y + world.Player.Scale.Y

	// Check if bounding_box is colliding with another bounding box
	for i := range world.BoundingBoxes {
		if getDistance(world.Player.Position.X, world.Player.Position.Z,
			world.BoundingBoxes[i].Min.X, world.BoundingBoxes[i].Min.Z) <= world.CalculationDistance &&
			rl.CheckCollisionBoxes(bounding_box, world.BoundingBoxes[i]) {

			return true
		}
	}

	return false
}

// Checks if the player is on the ground after moving by world.FloatPrecision on the Y axis
//
// #1 return: bool - true if the player is on the ground
func (world *World) isPlayerOnGroundNextFrame() bool {
	bounding_box := world.Player.BoundingBox

	// Move bounding_box by world.FloatPrecision on the Y axis
	bounding_box.Min.Y -= world.FloatPrecision

	// Check if bounding_box is on the ground
	if bounding_box.Min.Y <= world.Ground && bounding_box.Max.Y > world.Ground {
		return true
	}
	// Check if bounding_box is colliding with a bounding box
	if i, _ := world.checkPlayerCollisionsYNextFrame(); i != -1 {
		return true
	}

	return false
}
