package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	speed             Speeds
	mouse_sensitivity Sensitivities
	fovs              FOVs
	position          rl.Vector3
	rotation          rl.Vector2
	scale             rl.Vector3
	const_scale       Scale
	is_crouching      bool
	y_velocity        float32
	gravity           float32
	jump_power        float32
	last_key_pressed  int32
	camera            rl.Camera3D
}

type Speeds struct {
	normal       float32
	sprint       float32
	sneak        float32
	current      float32
	acceleration float32
}

type Sensitivities struct {
	normal float32
	zoom   float32
}

type FOVs struct {
	normal float32
	zoom   float32
}

type Scale struct {
	normal float32
	crouch float32
}

func (player *Player) initPlayer() {
	player.speed.normal = .1
	player.speed.sprint = .15
	player.speed.sneak = .05
	player.speed.current = 0.
	player.speed.acceleration = .01
	player.mouse_sensitivity.normal = .0025
	player.mouse_sensitivity.zoom = .0005
	player.fovs.normal = 70.
	player.fovs.zoom = 20.
	player.rotation = rl.NewVector2(0., 0.)
	player.position = rl.NewVector3(4., .9, 4.)
	player.scale = rl.NewVector3(.8, 1.8, .8)
	player.const_scale.normal = 1.8
	player.const_scale.crouch = .9
	player.is_crouching = false
	player.y_velocity = 0.
	player.gravity = .0065
	player.jump_power = .15
	player.last_key_pressed = -1
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

func (player *Player) movePlayer(bounding_boxes []rl.BoundingBox) {
	half_crouch_scale := player.const_scale.crouch / 2

	if rl.IsKeyDown(rl.KeyLeftControl) {
		player.scale.Y = player.const_scale.crouch
		if !player.is_crouching {
			player.position.Y -= half_crouch_scale
		}
		player.is_crouching = true
	} else if player.checkPlayerUncrouch(bounding_boxes) {
		player.scale.Y = player.const_scale.normal
		if player.is_crouching {
			player.position.Y += half_crouch_scale
		}
		player.is_crouching = false
	}
	if rl.IsKeyDown(rl.KeySpace) && player.y_velocity == 0. && player.checkIfPlayerOnSurface(bounding_boxes) {
		player.y_velocity = player.jump_power
	}

	player.position.Y += player.y_velocity * (rl.GetFrameTime() * 60)

	player_position_after_moving := player.getPlayerPositionAfterMoving()
	collisions_x, collisions_z := player.checkCollisionsXZForPlayer(bounding_boxes)
	if collisions_x && collisions_z {
		return
	} else if collisions_x {
		player.position.Z = player_position_after_moving.Z
		return
	} else if collisions_z {
		player.position.X = player_position_after_moving.X
		return
	}

	player.position = player_position_after_moving
}

func (player *Player) rotatePlayer() {
	mouse_delta := rl.GetMouseDelta()
	if rl.IsKeyDown(rl.KeyC) {
		player.rotation.X += mouse_delta.X * player.mouse_sensitivity.zoom
		player.rotation.Y -= mouse_delta.Y * player.mouse_sensitivity.zoom
	} else {
		player.rotation.X += mouse_delta.X * player.mouse_sensitivity.normal
		player.rotation.Y -= mouse_delta.Y * player.mouse_sensitivity.normal
	}

	if player.rotation.Y > 1.5 {
		player.rotation.Y = 1.5
	}
	if player.rotation.Y < -1.5 {
		player.rotation.Y = -1.5
	}
}

func (player Player) checkCollisionsForPlayer(bounding_boxes []rl.BoundingBox) bool {
	player.position = player.getPlayerPositionAfterMoving()

	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.position.X-player.scale.X/2, player.position.Y-player.scale.Y/2, player.position.Z-player.scale.Z/2),
			rl.NewVector3(player.position.X+player.scale.X/2, player.position.Y+player.scale.Y/2, player.position.Z+player.scale.Z/2)), box) {
			return true
		}
	}

	return false
}

func (player Player) checkCollisionsXZForPlayer(bounding_boxes []rl.BoundingBox) (bool, bool) {
	collision_x, collision_z := false, false

	player_position_after_moving := player.getPlayerPositionAfterMoving()
	player_position_x := player_position_after_moving.X

	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player_position_x-player.scale.X/2, player.position.Y-player.scale.Y/2, player.position.Z-player.scale.Z/2),
			rl.NewVector3(player_position_x+player.scale.X/2, player.position.Y+player.scale.Y/2, player.position.Z+player.scale.Z/2)), box) {
			collision_x = true
		}
	}

	player_position_z := player_position_after_moving.Z

	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.position.X-player.scale.X/2, player.position.Y-player.scale.Y/2, player_position_z-player.scale.Z/2),
			rl.NewVector3(player.position.X+player.scale.X/2, player.position.Y+player.scale.Y/2, player_position_z+player.scale.Z/2)), box) {
			collision_z = true
		}
	}

	return collision_x, collision_z
}

func (player Player) getPlayerPositionAfterMoving() rl.Vector3 {
	frame_time := rl.GetFrameTime() * 60

	current_speed := player.speed.current

	if player.speed.normal == 0. {
		player.position.Y += player.y_velocity * frame_time
	}

	keys_pressed := 0
	if rl.IsKeyDown(rl.KeyW) {
		keys_pressed++
	}
	if rl.IsKeyDown(rl.KeyS) {
		keys_pressed++
	}
	if rl.IsKeyDown(rl.KeyA) {
		keys_pressed++
	}
	if rl.IsKeyDown(rl.KeyD) {
		keys_pressed++
	}
	if keys_pressed == 2 {
		current_speed = current_speed * .707
	}

	final_speed := current_speed * frame_time

	speeds := rl.NewVector2(
		float32(math.Cos(float64(player.rotation.X)))*final_speed,
		float32(math.Sin(float64(player.rotation.X)))*final_speed,
	)

	if rl.IsKeyDown(rl.KeyW) || player.last_key_pressed == int32(rl.KeyW) {
		player.position.X -= speeds.X
		player.position.Z -= speeds.Y
	}
	if rl.IsKeyDown(rl.KeyS) || player.last_key_pressed == int32(rl.KeyS) {
		player.position.X += speeds.X
		player.position.Z += speeds.Y
	}
	if rl.IsKeyDown(rl.KeyA) || player.last_key_pressed == int32(rl.KeyA) {
		player.position.Z += speeds.X
		player.position.X -= speeds.Y
	}
	if rl.IsKeyDown(rl.KeyD) || player.last_key_pressed == int32(rl.KeyD) {
		player.position.Z -= speeds.X
		player.position.X += speeds.Y
	}

	return player.position
}

func (player Player) checkPlayerUncrouch(bounding_boxes []rl.BoundingBox) bool {
	player.scale.Y = player.const_scale.normal
	player.position.Y += player.const_scale.normal / 2

	return !player.checkCollisionsForPlayer(bounding_boxes)
}

func (player *Player) applyGravityToPlayer(bounding_boxes []rl.BoundingBox) {
	player.y_velocity -= player.gravity * (rl.GetFrameTime() * 60)

	if player.checkCollisionsYForPlayer(bounding_boxes) || player.getPlayerPositionAfterMoving().Y-(player.scale.Y/2) < 0. {
		player.y_velocity = 0.
		return
	}
}

func (player Player) checkCollisionsYForPlayer(bounding_boxes []rl.BoundingBox) bool {
	player.speed.normal = 0
	player.speed.sprint = 0
	player.speed.sneak = 0
	player.speed.current = 0

	return player.checkCollisionsForPlayer(bounding_boxes)
}

func (player Player) checkIfPlayerOnSurface(bounding_boxes []rl.BoundingBox) bool {
	player.position.Y -= player.gravity * (rl.GetFrameTime() * 60)
	if player.checkCollisionsYForPlayer(bounding_boxes) || player.position.Y-(player.scale.Y/2) < 0. {
		return true
	}
	return false
}

func (player *Player) accelerationPlayer() {
	final_speed := player.speed.acceleration * rl.GetFrameTime() * 60
	keys_down := map[string]bool{"shift": rl.IsKeyDown(rl.KeyLeftShift), "ctrl": rl.IsKeyDown(rl.KeyLeftControl), "w": rl.IsKeyDown(rl.KeyW), "s": rl.IsKeyDown(rl.KeyS), "a": rl.IsKeyDown(rl.KeyA), "d": rl.IsKeyDown(rl.KeyD)}

	if !keys_down["w"] && !keys_down["s"] && !keys_down["a"] && !keys_down["d"] {
		if player.speed.current > 0. {
			player.speed.current -= final_speed
		} else {
			player.speed.current = 0.
		}
	} else if !keys_down["shift"] && player.speed.current > player.speed.normal {
		player.speed.current -= final_speed
	}
	if player.is_crouching && player.speed.current > player.speed.sneak {
		player.speed.current -= final_speed
	}

	if player.speed.current <= player.speed.normal && (keys_down["w"] || keys_down["s"] || keys_down["a"] || keys_down["d"]) && !keys_down["shift"] && !keys_down["ctrl"] {
		player.speed.current += final_speed
	}
	if keys_down["shift"] && player.speed.current <= player.speed.sprint && (keys_down["w"] || keys_down["s"] || keys_down["a"] || keys_down["d"]) {
		player.speed.current += final_speed
	}
	if keys_down["ctrl"] && player.speed.current <= player.speed.sneak && (keys_down["w"] || keys_down["s"] || keys_down["a"] || keys_down["d"]) {
		player.speed.current += final_speed
	}
}

func (player *Player) lastKeyPressedPlayer() {
	if rl.IsKeyDown(rl.KeyW) {
		player.last_key_pressed = int32(rl.KeyW)
	}
	if rl.IsKeyDown(rl.KeyS) {
		player.last_key_pressed = int32(rl.KeyS)
	}
	if rl.IsKeyDown(rl.KeyA) {
		player.last_key_pressed = int32(rl.KeyA)
	}
	if rl.IsKeyDown(rl.KeyD) {
		player.last_key_pressed = int32(rl.KeyD)
	}
}
