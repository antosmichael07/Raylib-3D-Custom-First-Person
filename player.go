package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	speed             float32
	sprint_speed      float32
	sneak_speed       float32
	mouse_sensitivity float32
	fov               float32
	position          rl.Vector3
	rotation          rl.Vector2
	scale             rl.Vector3
	camera            rl.Camera3D
}

func initPlayer() Player {
	player := Player{}
	player.speed = .1
	player.sprint_speed = .2
	player.sneak_speed = .05
	player.mouse_sensitivity = .0025
	player.fov = 70.
	player.rotation = rl.NewVector2(0., 0.)
	player.position = rl.NewVector3(4., .9, 4.)
	player.scale = rl.NewVector3(.8, 1.8, .8)
	player.camera = initCamera(player)

	return player
}

func updatePlayer(player *Player, object rl.BoundingBox) {
	movePlayer(player, object)
	rotatePlayer(player)
	updateCameraFirstPerson(player)
}

func movePlayer(player *Player, object rl.BoundingBox) {
	if rl.IsKeyDown(rl.KeyLeftControl) {
		player.scale.Y = 1.
	} else if !checkCollisionsForPlayer(*player, object) {
		player.scale.Y = 2.
	}

	if checkCollisionsForPlayer(*player, object) {
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

	player.rotation.X += rl.GetMouseDelta().X * player.mouse_sensitivity
	player.rotation.Y -= rl.GetMouseDelta().Y * player.mouse_sensitivity
}

func checkCollisionsForPlayer(player Player, object rl.BoundingBox) bool {
	player.position = getPlayerPositionAfterMoving(player)
	
	if rl.CheckCollisionBoxes(rl.NewBoundingBox(rl.NewVector3(player.position.X - player.scale.X / 2, player.position.Y - player.scale.Y / 2, player.position.Z - player.scale.Z / 2), 
	rl.NewVector3(player.position.X + player.scale.X / 2, player.position.Y + player.scale.Y / 2, player.position.Z + player.scale.Z / 2)), object) {
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
	if rl.IsKeyDown(rl.KeyLeftControl) {
		current_speed = player.sneak_speed
		player.scale.Y = 1.
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
