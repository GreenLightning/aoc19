package main

type Vector2 struct {
	x, y int
}

func (v Vector2) Plus(other Vector2) Vector2 {
	return Vector2{
		x: v.x + other.x,
		y: v.y + other.y,
	}
}

func (v Vector2) Minus(other Vector2) Vector2 {
	return Vector2{
		x: v.x - other.x,
		y: v.y - other.y,
	}
}

func (v Vector2) Times(factor int) Vector2 {
	return Vector2{
		x: factor * v.x,
		y: factor * v.y,
	}
}

func (v Vector2) Min(other Vector2) Vector2 {
	return Vector2{
		x: min(v.x, other.x),
		y: min(v.y, other.y),
	}
}

func (v Vector2) Max(other Vector2) Vector2 {
	return Vector2{
		x: max(v.x, other.x),
		y: max(v.y, other.y),
	}
}

func (v Vector2) LengthSquared() int {
	return v.x*v.x + v.y*v.y
}

func (v Vector2) ManhattenLength() int {
	return abs(v.x) + abs(v.y)
}

func (v Vector2) DistanceSquared(o Vector2) int {
	return v.Minus(o).LengthSquared()
}

func (v Vector2) ManhattenDistance(o Vector2) int {
	return v.Minus(o).ManhattenLength()
}

type Vector3 struct {
	x, y, z int
}

func (v Vector3) Plus(other Vector3) Vector3 {
	return Vector3{
		x: v.x + other.x,
		y: v.y + other.y,
		z: v.z + other.z,
	}
}

func (v Vector3) Minus(other Vector3) Vector3 {
	return Vector3{
		x: v.x - other.x,
		y: v.y - other.y,
		z: v.z - other.z,
	}
}

func (v Vector3) Times(factor int) Vector3 {
	return Vector3{
		x: factor * v.x,
		y: factor * v.y,
		z: factor * v.z,
	}
}

func (v Vector3) LengthSquared() int {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

func (v Vector3) ManhattenLength() int {
	return abs(v.x) + abs(v.y) + abs(v.z)
}

func (v Vector3) DistanceSquared(o Vector3) int {
	return v.Minus(o).LengthSquared()
}

func (v Vector3) ManhattenDistance(o Vector3) int {
	return v.Minus(o).ManhattenLength()
}
