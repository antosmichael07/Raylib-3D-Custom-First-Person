package main

import (
	"fmt"
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
	player.sprint_speed = .15
	player.sneak_speed = .05
	player.current_speed = 0.
	player.acceleration = .01
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
	player.gravity = .0065
	player.jump_power = .15
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

	player.position.Y += player.y_velocity * (rl.GetFrameTime() * 60)

	collisions_x, collisions_z := checkCollisionsXZForPlayer(*player, bounding_boxes)
	if collisions_x && collisions_z {
		return
	} else if collisions_x {
		player.position.Z = getPlayerPositionAfterMoving(*player).Z
		return
	} else if collisions_z {
		player.position.X = getPlayerPositionAfterMoving(*player).X
		return
	}

	player.position = getPlayerPositionAfterMoving(*player)
}

func rotatePlayer(player *Player) {
	if rl.IsKeyDown(rl.KeyC) {
		player.rotation.X += rl.GetMouseDelta().X * player.zoom_mouse_sensitivity
		player.rotation.Y -= rl.GetMouseDelta().Y * player.zoom_mouse_sensitivity
	} else {
		player.rotation.X += rl.GetMouseDelta().X * player.mouse_sensitivity
		player.rotation.Y -= rl.GetMouseDelta().Y * player.mouse_sensitivity
	}

	if player.rotation.Y > 1.5 {
		player.rotation.Y = 1.5
	}
	if player.rotation.Y < -1.5 {
		player.rotation.Y = -1.5
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

func checkCollisionsXZForPlayer(player Player, bounding_boxes []rl.BoundingBox) (bool, bool) {
	collision_x, collision_z := false, false

	player_position_x := getPlayerPositionAfterMoving(player).X

	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player_position_x-player.scale.X/2, player.position.Y-player.scale.Y/2, player.position.Z-player.scale.Z/2),
			rl.NewVector3(player_position_x+player.scale.X/2, player.position.Y+player.scale.Y/2, player.position.Z+player.scale.Z/2)), box) {
			collision_x = true
		}
	}

	player_position_z := getPlayerPositionAfterMoving(player).Z

	for _, box := range bounding_boxes {
		if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.position.X-player.scale.X/2, player.position.Y-player.scale.Y/2, player_position_z-player.scale.Z/2),
			rl.NewVector3(player.position.X+player.scale.X/2, player.position.Y+player.scale.Y/2, player_position_z+player.scale.Z/2)), box) {
			collision_z = true
		}
	}

	fmt.Printf("player_y: %f\n", player.position.Y)
	return collision_x, collision_z
}

func getPlayerPositionAfterMoving(player Player) rl.Vector3 {
	current_speed := player.current_speed
	if player.speed == 0. {
		player.position.Y += player.y_velocity * (rl.GetFrameTime() * 60)
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

func checkPlayerUncrouch(player Player, bounding_boxes []rl.BoundingBox) bool {
	player.scale.Y = player.crouch_scale.Y
	player.position.Y += player.crouch_scale.Y / 2

	return !checkCollisionsForPlayer(player, bounding_boxes)
}

func applyGravityToPlayer(player *Player, bounding_boxes []rl.BoundingBox) {
	player.y_velocity -= player.gravity * (rl.GetFrameTime() * 60)

	if checkCollisionsYForPlayer(*player, bounding_boxes) || getPlayerPositionAfterMoving(*player).Y-(player.scale.Y/2) < 0. {
		player.y_velocity = 0.
		return
	}
}

func checkCollisionsYForPlayer(player Player, bounding_boxes []rl.BoundingBox) bool {
	player.speed = 0
	player.sprint_speed = 0
	player.sneak_speed = 0
	player.current_speed = 0

	return checkCollisionsForPlayer(player, bounding_boxes)
}

func checkIfPlayerOnSurface(player Player, bounding_boxes []rl.BoundingBox) bool {
	player.position.Y -= player.gravity * (rl.GetFrameTime() * 60)
	if checkCollisionsYForPlayer(player, bounding_boxes) || player.position.Y-(player.scale.Y/2) < 0. {
		return true
	}
	return false
}

func accelerationPlayer(player *Player) {
	if !rl.IsKeyDown(rl.KeyW) && !rl.IsKeyDown(rl.KeyS) && !rl.IsKeyDown(rl.KeyA) && !rl.IsKeyDown(rl.KeyD) {
		if player.current_speed > 0. {
			player.current_speed -= player.acceleration * (rl.GetFrameTime() * 60)
		} else {
			player.current_speed = 0.
		}
	} else if !rl.IsKeyDown(rl.KeyLeftShift) && player.current_speed > player.speed {
		player.current_speed -= player.acceleration * (rl.GetFrameTime() * 60)
	}
	if rl.IsKeyDown(rl.KeyLeftControl) && player.current_speed > player.sneak_speed {
		player.current_speed -= player.acceleration * (rl.GetFrameTime() * 60)
	}

	if player.current_speed <= player.speed && (rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyD)) && !rl.IsKeyDown(rl.KeyLeftShift) && !rl.IsKeyDown(rl.KeyLeftControl) {
		player.current_speed += player.acceleration * (rl.GetFrameTime() * 60)
	}
	if rl.IsKeyDown(rl.KeyLeftShift) && player.current_speed <= player.sprint_speed && (rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyD)) {
		player.current_speed += player.acceleration * (rl.GetFrameTime() * 60)
	}
	if rl.IsKeyDown(rl.KeyLeftControl) && player.current_speed <= player.sneak_speed && (rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyD)) {
		player.current_speed += player.acceleration * (rl.GetFrameTime() * 60)
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
