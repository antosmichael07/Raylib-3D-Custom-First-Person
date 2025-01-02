package rlfp

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Used for managing world around the player and the player itself
type World struct {
	Player Player
	Ground float32
	// Value which is subtracted from player's and object's Y velocity
	Gravity float32
	// Used for moving when players have different target FPS
	FrameTime     float32
	LastFrameTime float32
	// Boxes with collisions
	BoundingBoxes []rl.BoundingBox
	// Boxes that activate when a player walks into them
	TriggerBoxes []TriggerBox
	// Boxes that activate when a player presses a key when looking at them
	InteractableBoxes []InteractableBox
	// Minimum value for working with floats
	FloatPrecision float32
	// Distance where the player has to be in to update certain elements in the world
	CalculationDistance float32
	// For interactable boxes, so they don't update their state every frame
	AlreadySetInteractStates bool
}

// Initializes default values for the world
//
// #1 argument ground: float32 - height of the ground
func (world *World) Init(ground float32) {
	world.Player.Init()
	world.Ground = ground
	world.Gravity = 15.
	world.FrameTime = 0.
	world.FloatPrecision = .0001
	// The actual distance is math32.Sqrt(world.CalculationDistance)
	// This is because the program is faster without square rooting and it has the same effect
	world.CalculationDistance = 40000.
	world.AlreadySetInteractStates = false
}

// Creates a new world with the player at the specified position, should be called when loading a save
//
// #1 argument position: rl.Vector3 - position of the player
//
// #2 argument rotation: rl.Vector2 - rotation of the player
//
// #3 argument is_crouching: bool - if the player is crouching
func (world *World) New(position rl.Vector3, rotation rl.Vector2, is_crouching bool) {
	world.Player.New(position, rotation, is_crouching)
	world.BoundingBoxes = []rl.BoundingBox{}
	world.TriggerBoxes = []TriggerBox{}
	world.InteractableBoxes = []InteractableBox{}
}

// Adds a new bounding box to the world
//
// #1 argument box: rl.BoundingBox - bounding box to add
func (world *World) AddBoundingBox(box rl.BoundingBox) {
	world.BoundingBoxes = append(world.BoundingBoxes, box)
}

// Updates every value in the world struct, should be called every frame
//
// #1 argument windowWidth: int32 - width of the window
//
// #2 argument windowHeight: int32 - height of the window
func (world *World) Update(windowWidth, windowHeight int32) {
	world.LastFrameTime = world.FrameTime
	world.FrameTime = rl.GetFrameTime()
	if world.FrameTime > 1. {
		world.FrameTime = 1.
	}
	world.UpdatePlayer()
	world.UpdateTriggerBoxes()
	world.UpdateInteractableBoxes(windowWidth, windowHeight)
}
