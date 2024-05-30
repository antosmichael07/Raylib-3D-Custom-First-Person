package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	speed                  float32
	sprint_speed           float32
	sneak_speed            float32
	mouse_sensitivity      float32
	zoom_mouse_sensitivity float32
	fov                    float32
	zoom_fov               float32
	position               rl.Vector3
	rotation               rl.Vector2
	scale                  rl.Vector3
	crouch_scale           rl.Vector2
	is_crouching           bool
	camera                 rl.Camera3D
}

func initPlayer() Player {
	player := Player{}
	player.speed = .1
	player.sprint_speed = .2
	player.sneak_speed = .05
	player.mouse_sensitivity = .0025
	player.zoom_mouse_sensitivity = .0005
	player.fov = 70.
	player.zoom_fov = 20.
	player.rotation = rl.NewVector2(0., 0.)
	player.position = rl.NewVector3(4., .9, 4.)
	player.scale = rl.NewVector3(.8, 1.8, .8)
	player.crouch_scale = rl.NewVector2(1.,2.)
	player.camera = initCamera(player)

	return player
}

func updatePlayer(player *Player, bounding_box rl.BoundingBox) {
	movePlayer(player, bounding_box)
	rotatePlayer(player)
	updateCameraFirstPerson(player)
}

func movePlayer(player *Player, bounding_box rl.BoundingBox) {
	if rl.IsKeyDown(rl.KeyLeftControl) {
		player.scale.Y = player.crouch_scale.X
		player.is_crouching = true
	} else if checkPlayerUncrouch(*player, bounding_box) {
		player.scale.Y = player.crouch_scale.Y
		player.is_crouching = false
	}

	if checkCollisionsForPlayer(*player, bounding_box) {
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

func checkCollisionsForPlayer(player Player, bounding_box rl.BoundingBox) bool {
	player.position = getPlayerPositionAfterMoving(player)
	
	if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.position.X - player.scale.X / 2, player.position.Y - player.scale.Y / 2, player.position.Z - player.scale.Z / 2), 
	rl.NewVector3(player.position.X + player.scale.X / 2, player.position.Y + player.scale.Y / 2, player.position.Z + player.scale.Z / 2)), bounding_box) {
		return true
	}
	return false
}

func getPlayerPositionAfterMoving(player Player) rl.Vector3 {
	position := player.position

	current_speed := player.speed
	if rl.IsKeyDown(rl.KeyLeftShift) {
		current_speed = player.sprint_speed
	}
	if player.is_crouching {
		current_speed = player.sneak_speed
	}
	if rl.IsKeyDown(rl.KeyW) && rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyW) && rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyS) && rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyS) && rl.IsKeyDown(rl.KeyD) {
		current_speed = current_speed * .707
	}

	speeds := rl.NewVector2(
		float32(math.Cos(float64(player.rotation.X)))*current_speed*(rl.GetFrameTime()*60),
		float32(math.Sin(float64(player.rotation.X)))*current_speed*(rl.GetFrameTime()*60),
	)

	if rl.IsKeyDown(rl.KeyW) {
		position.X -= speeds.X
		position.Z -= speeds.Y
	}
	if rl.IsKeyDown(rl.KeyS) {
		position.X += speeds.X
		position.Z += speeds.Y
	}
	if rl.IsKeyDown(rl.KeyA) {
		position.Z += speeds.X
		position.X -= speeds.Y
	}
	if rl.IsKeyDown(rl.KeyD) {
		position.Z -= speeds.X
		position.X += speeds.Y
	}

	return position
}

func checkPlayerUncrouch(player Player, bounding_box rl.BoundingBox) bool {
	player.scale.Y = player.crouch_scale.Y

	if checkCollisionsForPlayer(player, bounding_box) {
		return false
	}

	return true
}
