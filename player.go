package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	speed                  float32
	sprint_speed           float32
	sneak_speed            float32
	current_speed          float32
	acceleration           float32
	mouse_sensitivity      float32
	zoom_mouse_sensitivity float32
	fov                    float32
	zoom_fov               float32
	position               rl.Vector3
	rotation               rl.Vector2
	scale                  rl.Vector3
	crouch_scale           rl.Vector2
	is_crouching           bool
	y_velocity             float32
	gravity                float32
	jump_power             float32
	last_key_pressed       int32
	camera                 rl.Camera3D
}

func initPlayer() Player {
	player := Player{}
	player.speed = .1
	player.sprint_speed = .2
	player.sneak_speed = .05
	player.current_speed = 0.
	player.acceleration = .004
	player.mouse_sensitivity = .0025
	player.zoom_mouse_sensitivity = .0005
	player.fov = 70.
	player.zoom_fov = 20.
	player.rotation = rl.NewVector2(0., 0.)
	player.position = rl.NewVector3(4., .9, 4.)
	player.scale = rl.NewVector3(.8, 1.8, .8)
	player.crouch_scale = rl.NewVector2(.9, 1.8)
	player.is_crouching = false
	player.y_velocity = 0.
	player.gravity = .003
	player.jump_power = .065
	player.last_key_pressed = -1
	player.camera = initCamera(player)

	return player
}

func updatePlayer(player *Player, bounding_boxes []rl.BoundingBox) {
	lastKeyPressedPlayer(player)
	accelerationPlayer(player)
	movePlayer(player, bounding_boxes)
	rotatePlayer(player)
	applyGravityToPlayer(player, bounding_boxes)
	updateCameraFirstPerson(player)
}

func movePlayer(player *Player, bounding_boxes []rl.BoundingBox) {
	if rl.IsKeyDown(rl.KeyLeftControl) {
		player.scale.Y = player.crouch_scale.X
		if !player.is_crouching {
			player.position.Y -= player.crouch_scale.X / 2
		}
		player.is_crouching = true
	} else if checkPlayerUncrouch(*player, bounding_boxes) {
		player.scale.Y = player.crouch_scale.Y
		if player.is_crouching {
			player.position.Y += player.crouch_scale.X / 2
		}
		player.is_crouching = false
	}
	if rl.IsKeyDown(rl.KeySpace) && player.y_velocity == 0. && checkIfPlayerOnSurface(*player, bounding_boxes) {
		player.y_velocity = player.jump_power
	}

	if checkCollisionsForPlayer(*player, bounding_boxes) {
		return
	}

	player.position = getPlayerPositionAfterMoving(*player)
}

func rotatePlayer(player *Player) {
	if player.rotation.Y > 1.5 {
		player.rotation.Y = 1.5
	}
	if player.rotation.Y < -1.5 {
		player.rotation.Y = -1.5
	}

	if rl.IsKeyDown(rl.KeyC) {
		player.rotation.X += rl.GetMouseDelta().X * player.zoom_mouse_sensitivity
		player.rotation.Y -= rl.GetMouseDelta().Y * player.zoom_mouse_sensitivity
	} else {
		player.rotation.X += rl.GetMouseDelta().X * player.mouse_sensitivity
		player.rotation.Y -= rl.GetMouseDelta().Y * player.mouse_sensitivity
	}
}

func checkCollisionsForPlayer(player Player, bounding_boxes []rl.BoundingBox) bool {
	player.position = getPlayerPositionAfterMoving(player)

	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.position.X-player.scale.X/2, player.position.Y-player.scale.Y/2, player.position.Z-player.scale.Z/2),
			rl.NewVector3(player.position.X+player.scale.X/2, player.position.Y+player.scale.Y/2, player.position.Z+player.scale.Z/2)), box) {
			return true
		}
	}

	return false
}

func getPlayerPositionAfterMoving(player Player) rl.Vector3 {
	position := player.position
	position.Y += player.y_velocity

	current_speed := player.current_speed
	if rl.IsKeyDown(rl.KeyLeftShift) {
		current_speed = player.sprint_speed
	}
	if player.is_crouching {
		current_speed = player.sneak_speed
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

	speeds := rl.NewVector2(
		float32(math.Cos(float64(player.rotation.X)))*current_speed*(rl.GetFrameTime()*60),
		float32(math.Sin(float64(player.rotation.X)))*current_speed*(rl.GetFrameTime()*60),
	)

	if rl.IsKeyDown(rl.KeyW) || player.last_key_pressed == int32(rl.KeyW) {
		position.X -= speeds.X
		position.Z -= speeds.Y
	}
	if rl.IsKeyDown(rl.KeyS) || player.last_key_pressed == int32(rl.KeyS) {
		position.X += speeds.X
		position.Z += speeds.Y
	}
	if rl.IsKeyDown(rl.KeyA) || player.last_key_pressed == int32(rl.KeyA) {
		position.Z += speeds.X
		position.X -= speeds.Y
	}
	if rl.IsKeyDown(rl.KeyD) || player.last_key_pressed == int32(rl.KeyD) {
		position.Z -= speeds.X
		position.X += speeds.Y
	}

	return position
}

func checkPlayerUncrouch(player Player, bounding_boxes []rl.BoundingBox) bool {
	player.scale.Y = player.crouch_scale.Y
	player.position.Y += player.crouch_scale.Y / 2

	return !checkCollisionsForPlayer(player, bounding_boxes)
}

func applyGravityToPlayer(player *Player, bounding_boxes []rl.BoundingBox) {
	player.y_velocity -= player.gravity * (rl.GetFrameTime() * 60)

	if checkYCollisionsForPlayer(*player, bounding_boxes) || getPlayerPositionAfterMoving(*player).Y-(player.scale.Y/2) < 0. {
		player.y_velocity = 0.
		return
	}

	if checkCollisionsForPlayer(*player, bounding_boxes) {
		player.position.Y += player.y_velocity
	}
}

func checkYCollisionsForPlayer(player Player, bounding_boxes []rl.BoundingBox) bool {
	player.speed = 0
	player.sprint_speed = 0
	player.sneak_speed = 0

	return checkCollisionsForPlayer(player, bounding_boxes)
}

func checkIfPlayerOnSurface(player Player, bounding_boxes []rl.BoundingBox) bool {
	player.position.Y -= player.gravity * (rl.GetFrameTime() * 60)
	if checkCollisionsForPlayer(player, bounding_boxes) || player.position.Y-(player.scale.Y/2) < 0. {
		return true
	}
	return false
}

func accelerationPlayer(player *Player) {
	if player.current_speed <= player.speed && (rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyD)) {
		player.current_speed += player.acceleration
	}
	if rl.IsKeyDown(rl.KeyLeftShift) && player.current_speed <= player.sprint_speed && (rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyD)) {
		player.current_speed += player.acceleration
	}
	if rl.IsKeyDown(rl.KeyLeftControl) && player.current_speed <= player.sneak_speed && (rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyD)) {
		player.current_speed += player.acceleration
	}

	if !rl.IsKeyDown(rl.KeyW) && !rl.IsKeyDown(rl.KeyS) && !rl.IsKeyDown(rl.KeyA) && !rl.IsKeyDown(rl.KeyD) && player.current_speed > 0. {
		player.current_speed -= player.acceleration
	}
}

func lastKeyPressedPlayer(player *Player) {
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
