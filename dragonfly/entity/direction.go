package entity

import (
	"git.jetbrains.space/dragonfly/dragonfly.git/dragonfly/block"
	"git.jetbrains.space/dragonfly/dragonfly.git/dragonfly/world"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

// Facing returns the horizontal direction that an entity is facing.
func Facing(e world.Entity) world.Face {
	yaw := math.Mod(float64(e.Yaw())-90, 360)
	if yaw < 0 {
		yaw += 360
	}
	switch {
	case (yaw > 0 && yaw < 45) || (yaw > 315 && yaw < 360):
		return world.West
	case yaw > 45 && yaw < 135:
		return world.North
	case yaw > 135 && yaw < 225:
		return world.East
	case yaw > 225 && yaw < 315:
		return world.South
	}
	return 0
}

// DirectionVector returns a vector that describes the direction of the entity passed. The length of the Vec3
// returned is always 1.
func DirectionVector(e world.Entity) mgl32.Vec3 {
	yaw, pitch := float64(mgl32.DegToRad(e.Yaw())), float64(mgl32.DegToRad(e.Pitch()))
	m := math.Cos(pitch)

	return mgl32.Vec3{
		float32(-m * math.Sin(yaw)),
		float32(-math.Sin(pitch)),
		float32(m * math.Cos(yaw)),
	}.Normalize()
}

// TargetBlock finds the target block of the entity passed. The block position returned will be at most
// maxDistance away from the entity. If no block can be found there, the block position returned will be
// that of an air block.
func TargetBlock(e world.Entity, maxDistance float64) world.BlockPos {
	// TODO: Implement accurate ray tracing for this.
	directionVector := DirectionVector(e)
	current := e.Position()
	if eyed, ok := e.(Eyed); ok {
		current = current.Add(mgl32.Vec3{0, eyed.EyeHeight()})
	}

	step := 0.5
	for i := 0.0; i < maxDistance; i += step {
		current = current.Add(directionVector.Mul(float32(step)))
		pos := vec3ToPos(current)

		b := e.World().Block(pos)
		if _, ok := b.(block.Air); !ok {
			// We hit a block that isn't air.
			return pos
		}
	}
	return vec3ToPos(current)
}

// Eyed represents an entity that has eyes.
type Eyed interface {
	// EyeHeight returns the offset from their base position that the eyes of an entity are found at.
	EyeHeight() float32
}

// vec3ToPos converts a Vec3 to a world.BlockPos.
func vec3ToPos(vec mgl32.Vec3) world.BlockPos {
	return world.BlockPos{int(math.Floor(float64(vec[0]))), int(math.Floor(float64(vec[1]))), int(math.Floor(float64(vec[2])))}
}