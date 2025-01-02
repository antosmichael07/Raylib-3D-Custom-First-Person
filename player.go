package rlfp

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chewxy/math32"
)

// Indexes of the player.CurrentInputs array
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

// Used for managing the player
type Player struct {
	// Player's constant speed variables
	Speed PlayerSpeeds
	// Player's sensitivities when zooming or not
	MouseSensitivity PlayerSensitivities
	// Player's normal and zoom FOV
	Fovs        PlayerFOVs
	BoundingBox rl.BoundingBox
	Position    rl.Vector3
	// Player's position is updated by this value
	OffsetNextFrame rl.Vector3
	Rotation        rl.Vector2
	Scale           rl.Vector3
	// Player's constant scale variables for crouching and normal scale
	ConstScale  PlayerScale
	IsCrouching bool
	// Player's position Y is updated by this value
	YVelocity float32
	// How much the player jumps
	JumpPower float32
	// Used for moving when no keys are pressed
	LastDirectionalKeyPressed int32
	// Range where the player can interact with an interactable box
	InteractRange float32
	// If the player interacted last frame
	AlreadyInteracted bool
	// How high can the player step up
	StepHeight float32
	// Constant controls (setting)
	Controls [ControlCount]int32
	// Current keys that are down
	CurrentInputs [ControlCount]bool
	Camera        rl.Camera3D
}

// Player's constant speed variables
type PlayerSpeeds struct {
	// Max values
	Normal float32
	Sprint float32
	Sneak  float32
	// Current speed
	Current float32
	// Acceleration
	Acceleration float32
}

// Player's sensitivities when zooming or not
type PlayerSensitivities struct {
	Normal float32
	Zoom   float32
}

// Player's normal and zoom FOV
type PlayerFOVs struct {
	Normal float32
	Zoom   float32
}

// Player's constant scale variables for crouching and normal scale
type PlayerScale struct {
	Normal float32
	Crouch float32
}

// Initializtion of default values, should be called at the start of the program
func (player *Player) Init() {
	player.Speed.Normal = 3.5
	player.Speed.Sprint = 6.
	player.Speed.Sneak = 1.5
	player.Speed.Acceleration = 25.
	player.MouseSensitivity.Normal = .0025
	player.MouseSensitivity.Zoom = .0005
	player.Fovs.Normal = 70.
	player.Fovs.Zoom = 20.
	player.ConstScale.Normal = 1.8
	player.ConstScale.Crouch = .9
	player.Scale = rl.Vector3{X: .6, Y: player.ConstScale.Normal, Z: .6}
	player.JumpPower = 5.
	player.InteractRange = 3.
	player.StepHeight = .4
	player.Controls[ControlForward] = rl.KeyW
	player.Controls[ControlBackward] = rl.KeyS
	player.Controls[ControlLeft] = rl.KeyA
	player.Controls[ControlRight] = rl.KeyD
	player.Controls[ControlJump] = rl.KeySpace
	player.Controls[ControlCrouch] = rl.KeyLeftControl
	player.Controls[ControlSprint] = rl.KeyLeftShift
	player.Controls[ControlZoom] = rl.KeyC
	player.Controls[ControlInteract] = rl.KeyE
}

// Initializes player's values, should be called when loading a save or starting a new game
func (player *Player) New(position rl.Vector3, rotation rl.Vector2, is_crouching bool) {
	player.Speed.Current = 0.
	player.Position = position
	player.OffsetNextFrame = rl.Vector3{X: 0., Y: 0., Z: 0.}
	player.Rotation = rotation
	if is_crouching {
		player.Scale.Y = player.ConstScale.Crouch
	} else {
		player.Scale.Y = player.ConstScale.Normal
	}
	player.BoundingBox = rl.BoundingBox{
		Min: rl.Vector3{
			X: player.Position.X - player.Scale.X/2.,
			Y: player.Position.Y - player.Scale.Y/2.,
			Z: player.Position.Z - player.Scale.Z/2.,
		},
		Max: rl.Vector3{
			X: player.Position.X + player.Scale.X/2.,
			Y: player.Position.Y + player.Scale.Y/2.,
			Z: player.Position.Z + player.Scale.Z/2.,
		},
	}
	player.IsCrouching = is_crouching
	player.YVelocity = 0.
	player.LastDirectionalKeyPressed = -1
	player.AlreadyInteracted = false
	player.CurrentInputs = [ControlCount]bool{false, false, false, false, false, false, false, false, false}
	player.InitCamera()
}

// Updating player, should be called every frame
func (world *World) UpdatePlayer() {
	// Update variables that don't affect player's current position
	world.UpdatePlayerVariables()
	// Updates player's position and states
	world.Player.UpdateRotation()
	world.UpdatePlayerCrouch()
	world.UpdatePlayerPosition()
	// Move camera to player's position and rotate it
	world.Player.UpdateCamera()
}

// Updates variables, that don't affect player's current position
func (world *World) UpdatePlayerVariables() {
	world.Player.UpdateCurrentInputs()
	world.Player.UpdateLastDirectionalKeyPressed()
	world.UpdatePlayerCurrentSpeed()
}

// Gets current keys down
func (player *Player) UpdateCurrentInputs() {
	player.CurrentInputs[ControlForward] = rl.IsKeyDown(player.Controls[ControlForward])
	player.CurrentInputs[ControlBackward] = rl.IsKeyDown(player.Controls[ControlBackward])
	player.CurrentInputs[ControlLeft] = rl.IsKeyDown(player.Controls[ControlLeft])
	player.CurrentInputs[ControlRight] = rl.IsKeyDown(player.Controls[ControlRight])
	player.CurrentInputs[ControlJump] = rl.IsKeyDown(player.Controls[ControlJump])
	player.CurrentInputs[ControlCrouch] = rl.IsKeyDown(player.Controls[ControlCrouch])
	player.CurrentInputs[ControlSprint] = rl.IsKeyDown(player.Controls[ControlSprint])
	player.CurrentInputs[ControlZoom] = rl.IsKeyDown(player.Controls[ControlZoom])
	player.CurrentInputs[ControlInteract] = rl.IsKeyDown(player.Controls[ControlInteract])
}

// Updates last key pressed, for moving without holding anything
func (player *Player) UpdateLastDirectionalKeyPressed() {
	if player.CurrentInputs[ControlForward] {
		player.LastDirectionalKeyPressed = player.Controls[ControlForward]
	}
	if player.CurrentInputs[ControlBackward] {
		player.LastDirectionalKeyPressed = player.Controls[ControlBackward]
	}
	if player.CurrentInputs[ControlLeft] {
		player.LastDirectionalKeyPressed = player.Controls[ControlLeft]
	}
	if player.CurrentInputs[ControlRight] {
		player.LastDirectionalKeyPressed = player.Controls[ControlRight]
	}
}

// Updates player's current speed
func (world *World) UpdatePlayerCurrentSpeed() {
	is_player_on_ground_next_frame := world.isPlayerOnGroundNextFrame()

	// When the player isn't holding anything, slow him to zero
	if !world.Player.CurrentInputs[ControlForward] && !world.Player.CurrentInputs[ControlBackward] &&
		!world.Player.CurrentInputs[ControlLeft] && !world.Player.CurrentInputs[ControlRight] {

		if world.Player.Speed.Current > 0. {
			world.Player.Speed.Current -= world.Player.Speed.Acceleration * world.FrameTime
			return
		} else {
			world.Player.Speed.Current = 0.
			return
		}
	}
	// When the player is faster while crouching then he should be
	if world.Player.IsCrouching && world.Player.Speed.Current > world.Player.Speed.Sneak {
		world.Player.Speed.Current -= world.Player.Speed.Acceleration * world.FrameTime
		return
	}
	// When the player is faster while not sprinting then he should be
	if (!world.Player.CurrentInputs[ControlSprint] || !is_player_on_ground_next_frame) &&
		world.Player.Speed.Current > world.Player.Speed.Normal {

		world.Player.Speed.Current -= world.Player.Speed.Acceleration * world.FrameTime
		return
	}

	// Add speed, when the player is slower than he should be
	if world.Player.Speed.Current <= world.Player.Speed.Normal && !world.Player.IsCrouching {
		world.Player.Speed.Current += world.Player.Speed.Acceleration * world.FrameTime
		return
	}
	// Add speed, when the player is sprinting and is slower than he should be
	if world.Player.CurrentInputs[ControlSprint] && world.Player.Speed.Current <= world.Player.Speed.Sprint {
		world.Player.Speed.Current += world.Player.Speed.Acceleration * world.FrameTime
		return
	}
	// Add speed, when the player is crouching and is slower than he should be
	if world.Player.IsCrouching && world.Player.Speed.Current <= world.Player.Speed.Sneak {
		world.Player.Speed.Current += world.Player.Speed.Acceleration * world.FrameTime
		return
	}
}

// Updates player's rotation
func (player *Player) UpdateRotation() {
	// Current mouse movement
	mouse_delta := rl.GetMouseDelta()

	// Rotate player with according sensitivity
	if player.CurrentInputs[ControlZoom] {
		player.Rotation.X += mouse_delta.X * player.MouseSensitivity.Zoom
		player.Rotation.Y -= mouse_delta.Y * player.MouseSensitivity.Zoom
	} else {
		player.Rotation.X += mouse_delta.X * player.MouseSensitivity.Normal
		player.Rotation.Y -= mouse_delta.Y * player.MouseSensitivity.Normal
	}

	// When the player is looking up or down, limit his view
	if player.Rotation.Y > 1.57 {
		player.Rotation.Y = 1.57
	}
	if player.Rotation.Y < -1.57 {
		player.Rotation.Y = -1.57
	}
}

// Updates player's crouch state
func (world *World) UpdatePlayerCrouch() {
	// Set player to crouching state
	if world.Player.CurrentInputs[ControlCrouch] {
		if !world.Player.IsCrouching {
			world.Player.Scale.Y = world.Player.ConstScale.Crouch
			world.Player.Position.Y -= world.Player.ConstScale.Crouch / 2
			world.Player.BoundingBox.Max.Y = world.Player.BoundingBox.Min.Y + world.Player.ConstScale.Crouch
			world.Player.IsCrouching = true

			return
		}

		return
	}

	// Set player to normal state
	if world.Player.IsCrouching && world.CanPlayerUncrouch() {
		world.Player.Scale.Y = world.Player.ConstScale.Normal
		world.Player.Position.Y += world.Player.ConstScale.Normal / 2
		world.Player.BoundingBox.Max.Y = world.Player.BoundingBox.Min.Y + world.Player.ConstScale.Normal
		world.Player.IsCrouching = false

		return
	}
}

// Check collisions in the Y axis, if the player can uncrouch
//
// #1 return: bool - if the player can uncrouch
func (world *World) CanPlayerUncrouch() bool {
	bounding_box_next_frame := world.Player.BoundingBox
	bounding_box_next_frame.Max.Y = bounding_box_next_frame.Min.Y + world.Player.ConstScale.Normal

	for i := range world.BoundingBoxes {
		if rl.CheckCollisionBoxes(bounding_box_next_frame, world.BoundingBoxes[i]) && i != 0 {
			return false
		}
	}

	return true
}

// Updates player's position and bounding box
func (world *World) UpdatePlayerPosition() {
	// Jump when the player is on the ground and the jump key is pressed
	if world.Player.CurrentInputs[ControlJump] && world.Player.YVelocity == 0. &&
		world.isPlayerOnGroundNextFrame() && !world.Player.IsCrouching {
		world.Player.YVelocity = world.Player.JumpPower
	}

	// Get player's offsets for the next frame
	world.UpdatePlayerOffsetNextFrame()

	// Update player's position Y and Y velocity
	if world.Player.OffsetNextFrame.Y != 0 {
		world.UpdatePlayerPositionY()
	} else {
		world.Player.YVelocity -= world.Gravity * world.FrameTime
	}

	// Update player's position X
	if world.Player.OffsetNextFrame.X != 0 {
		world.UpdatePlayerPositionX()
	}
	// Update player's position Z
	if world.Player.OffsetNextFrame.Z != 0 {
		world.UpdatePlayerPositionZ()
	}
}

// Gets player's offsets for the next frame
func (world *World) UpdatePlayerOffsetNextFrame() {
	offset := rl.Vector2{X: 0., Y: 0.}
	// Get player's current speed
	current_speed := world.Player.Speed.Current

	// Get the number of keys pressed
	keys_pressed := 0
	if world.Player.CurrentInputs[ControlForward] {
		keys_pressed++
	}
	if world.Player.CurrentInputs[ControlBackward] {
		keys_pressed++
	}
	if world.Player.CurrentInputs[ControlLeft] {
		keys_pressed++
	}
	if world.Player.CurrentInputs[ControlRight] {
		keys_pressed++
	}
	// If you're moving diagonally (2 keys) then multiply by math32.Sqrt(2)
	if keys_pressed == 2 {
		current_speed *= .707
	}

	current_speed *= world.FrameTime

	// Calculate the offsets according to player's inputs
	if world.Player.CurrentInputs[ControlForward] || world.Player.LastDirectionalKeyPressed == world.Player.Controls[ControlForward] {
		offset.X -= math32.Cos(world.Player.Rotation.X) * current_speed
		offset.Y -= math32.Sin(world.Player.Rotation.X) * current_speed
	}
	if world.Player.CurrentInputs[ControlBackward] || world.Player.LastDirectionalKeyPressed == world.Player.Controls[ControlBackward] {
		offset.X += math32.Cos(world.Player.Rotation.X) * current_speed
		offset.Y += math32.Sin(world.Player.Rotation.X) * current_speed
	}
	if world.Player.CurrentInputs[ControlLeft] || world.Player.LastDirectionalKeyPressed == world.Player.Controls[ControlLeft] {
		offset.Y += math32.Cos(world.Player.Rotation.X) * current_speed
		offset.X -= math32.Sin(world.Player.Rotation.X) * current_speed
	}
	if world.Player.CurrentInputs[ControlRight] || world.Player.LastDirectionalKeyPressed == world.Player.Controls[ControlRight] {
		offset.Y -= math32.Cos(world.Player.Rotation.X) * current_speed
		offset.X += math32.Sin(world.Player.Rotation.X) * current_speed
	}

	// Update player's offsets
	world.Player.OffsetNextFrame.X = offset.X
	world.Player.OffsetNextFrame.Y = world.Player.YVelocity * world.FrameTime
	world.Player.OffsetNextFrame.Z = offset.Y
}

// Updates player's position X
func (world *World) UpdatePlayerPositionX() {
	// Check collisions in the X axis
	if i, t := world.checkPlayerCollisionsXNextFrame(); i != -1 {
		// Check if the player will be colliding with an object when stepping up
		if world.BoundingBoxes[i].Max.Y-world.Player.BoundingBox.Min.Y <= world.Player.StepHeight && world.isPlayerOnGroundNextFrame() {
			if !world.checkPlayerCollisionsXYNextFrame(world.BoundingBoxes[i].Max.Y + world.FloatPrecision) {
				// Move player in the X axis
				world.Player.BoundingBox.Min.X += world.Player.OffsetNextFrame.X
				world.Player.BoundingBox.Max.X += world.Player.OffsetNextFrame.X
				world.Player.Position.X += world.Player.OffsetNextFrame.X

				// Move the player in the Y axis
				world.Player.BoundingBox.Min.Y = world.BoundingBoxes[i].Max.Y + world.FloatPrecision
				world.Player.BoundingBox.Max.Y = world.Player.BoundingBox.Min.Y + world.Player.Scale.Y
				world.Player.Position.Y = world.Player.BoundingBox.Min.Y + world.Player.Scale.Y/2
				return
			}
		}

		if t {
			// Align to an object when moving in positive X axis
			world.Player.BoundingBox.Max.X = world.BoundingBoxes[i].Min.X - world.FloatPrecision
			world.Player.BoundingBox.Min.X = world.Player.BoundingBox.Max.X - world.Player.Scale.X
			world.Player.Position.X = world.Player.BoundingBox.Min.X + world.Player.Scale.X/2

			return
		} else {
			// Align to an object when moving in negative X axis
			world.Player.BoundingBox.Min.X = world.BoundingBoxes[i].Max.X + world.FloatPrecision
			world.Player.BoundingBox.Max.X = world.Player.BoundingBox.Min.X + world.Player.Scale.X
			world.Player.Position.X = world.Player.BoundingBox.Min.X + world.Player.Scale.X/2

			return
		}
	}

	// Move player in the X axis
	world.Player.BoundingBox.Min.X += world.Player.OffsetNextFrame.X
	world.Player.BoundingBox.Max.X += world.Player.OffsetNextFrame.X
	world.Player.Position.X += world.Player.OffsetNextFrame.X
}

// Updates player's position Y
func (world *World) UpdatePlayerPositionY() {
	// Check if the player is on the ground
	if world.Player.BoundingBox.Min.Y+world.Player.OffsetNextFrame.Y < world.Ground &&
		(world.Player.BoundingBox.Max.Y > world.Ground ||
			world.Player.BoundingBox.Min.Y+world.Player.OffsetNextFrame.Y+world.Player.StepHeight > world.Ground) {

		// Reset player's Y velocity when colliding with the ground
		world.Player.YVelocity = 0.

		// Check if the player will be colliding with an object when moving in the Y axis
		if i := world.checkPlayerCollisionsYOnGround(); i != -1 {
			// Align to an object when colliding
			world.Player.BoundingBox.Max.Y = world.BoundingBoxes[i].Min.Y - world.FloatPrecision
			world.Player.BoundingBox.Min.Y = world.Player.BoundingBox.Max.Y - world.Player.Scale.Y
			world.Player.Position.Y = world.Player.BoundingBox.Min.Y + world.Player.Scale.Y/2

			return
		} else {
			// Align to the ground
			world.Player.BoundingBox.Min.Y = world.Ground + world.FloatPrecision
			world.Player.BoundingBox.Max.Y = world.Player.BoundingBox.Min.Y + world.Player.Scale.Y
			world.Player.Position.Y = world.Player.BoundingBox.Min.Y + world.Player.Scale.Y/2

			return
		}
	}

	// Check collisions in the Y axis
	if i, t := world.checkPlayerCollisionsYNextFrame(); i != -1 {
		// Reset player's Y velocity when colliding with an object
		world.Player.YVelocity = 0.

		if t {
			// Align to an object when moving in positive Y axis
			world.Player.BoundingBox.Max.Y = world.BoundingBoxes[i].Min.Y - world.FloatPrecision
			world.Player.BoundingBox.Min.Y = world.Player.BoundingBox.Max.Y - world.Player.Scale.Y
			world.Player.Position.Y = world.Player.BoundingBox.Min.Y + world.Player.Scale.Y/2

			world.Player.OffsetNextFrame.Y = 0.

			return
		} else {
			// Align to an object when moving in negative Y axis
			world.Player.BoundingBox.Min.Y = world.BoundingBoxes[i].Max.Y + world.FloatPrecision
			world.Player.BoundingBox.Max.Y = world.Player.BoundingBox.Min.Y + world.Player.Scale.Y
			world.Player.Position.Y = world.Player.BoundingBox.Min.Y + world.Player.Scale.Y/2

			return
		}
	}

	// Move player in the Y axis
	world.Player.BoundingBox.Min.Y += world.Player.OffsetNextFrame.Y
	world.Player.BoundingBox.Max.Y += world.Player.OffsetNextFrame.Y
	world.Player.Position.Y += world.Player.OffsetNextFrame.Y

	// Update player's Y velocity
	world.Player.YVelocity -= world.Gravity * world.FrameTime
}

// Updates player's position Z
func (world *World) UpdatePlayerPositionZ() {
	// Check collisions in the Z axis
	if i, t := world.checkPlayerCollisionsZNextFrame(); i != -1 {
		// Check if the player will be colliding with an object when stepping up
		if world.BoundingBoxes[i].Max.Y-world.Player.BoundingBox.Min.Y <= world.Player.StepHeight && world.isPlayerOnGroundNextFrame() {
			if !world.checkPlayerCollisionsZYNextFrame(world.BoundingBoxes[i].Max.Y + world.FloatPrecision) {
				// Move player in the X axis
				world.Player.BoundingBox.Min.Z += world.Player.OffsetNextFrame.Z
				world.Player.BoundingBox.Max.Z += world.Player.OffsetNextFrame.Z
				world.Player.Position.Z += world.Player.OffsetNextFrame.Z

				// Move the player in the Y axis
				world.Player.BoundingBox.Min.Y = world.BoundingBoxes[i].Max.Y + world.FloatPrecision
				world.Player.BoundingBox.Max.Y = world.Player.BoundingBox.Min.Y + world.Player.Scale.Y
				world.Player.Position.Y = world.Player.BoundingBox.Min.Y + world.Player.Scale.Y/2
				return
			}
		}

		if t {
			// Align to an object when moving in positive Z axis
			world.Player.BoundingBox.Max.Z = world.BoundingBoxes[i].Min.Z - world.FloatPrecision
			world.Player.BoundingBox.Min.Z = world.Player.BoundingBox.Max.Z - world.Player.Scale.Z
			world.Player.Position.Z = world.Player.BoundingBox.Min.Z + world.Player.Scale.Z/2

			return
		} else {
			// Align to an object when moving in negative Z axis
			world.Player.BoundingBox.Min.Z = world.BoundingBoxes[i].Max.Z + world.FloatPrecision
			world.Player.BoundingBox.Max.Z = world.Player.BoundingBox.Min.Z + world.Player.Scale.Z
			world.Player.Position.Z = world.Player.BoundingBox.Min.Z + world.Player.Scale.Z/2

			return
		}
	}

	// Move player in the Z axis
	world.Player.BoundingBox.Min.Z += world.Player.OffsetNextFrame.Z
	world.Player.BoundingBox.Max.Z += world.Player.OffsetNextFrame.Z
	world.Player.Position.Z += world.Player.OffsetNextFrame.Z
}
