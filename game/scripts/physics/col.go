package physics

import (
	"hz/game/util"
	"math"
)

func isCirclesCollide(c1, c2 CircleShape) bool {
	return isCirclesCollideThreshold(c1, c2, 0.001)
}
func isCirclesCollideThreshold(c1, c2 CircleShape, threshold float64) bool {
	dx := c1.X - c2.X
	dy := c1.Y - c2.Y
	r := c1.Radius + c2.Radius
	return r*r-dx*dx-dy*dy > threshold
}

func isCircleRectangleCollision(c1 CircleShape, r1 RectShape) bool {
	return isCircleRectangleCollisionThreshold(c1, r1, 0.001)
}

// Проверяет пересечение окружности и прямоугольника
func isCircleRectangleCollisionThreshold(c1 CircleShape, r1 RectShape, threshold float64) bool {
	// Находим ближайшую точку прямоугольника к центру окружности
	nearestX := math.Max(r1.X, math.Min(c1.X, r1.X+r1.Width))
	nearestY := math.Max(r1.Y, math.Min(c1.Y, r1.Y+r1.Height))

	// Вычисляем расстояние от центра окружности до ближайшей точки
	dx := c1.X - nearestX
	dy := c1.Y - nearestY

	distanceSquared := dx*dx + dy*dy

	return (c1.Radius*c1.Radius)-distanceSquared > threshold
}

func solveCircleCollision(c1, c2 *CircleShape) {
	dx := c1.X - c2.X
	dy := c1.Y - c2.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	minDist := c1.Radius + c2.Radius

	if dist > minDist {
		return
	}

	if dist == 0 {
		dx, dy = 1, 0
		dist = minDist
	}

	delta := (minDist - dist) / 2
	nx, ny := dx/dist, dy/dist

	c1.X += nx * delta
	c1.Y += ny * delta
	c2.X -= nx * delta
	c2.Y -= ny * delta

	// newDist := math.Sqrt((c1.X-c2.X)*(c1.X-c2.X) + (c1.Y-c2.Y)*(c1.Y-c2.Y))
	// if newDist+0.001 < minDist {
	// 	log.Panic("collision not resolved")
	// }
}

func solveCircleRectCollision(c1 *CircleShape, r1 *RectShape) {
	// Находим ближайшую точку прямоугольника к центру окружности
	nearestX := math.Max(r1.X, math.Min(c1.X, r1.X+r1.Width))
	nearestY := math.Max(r1.Y, math.Min(c1.Y, r1.Y+r1.Height))

	// Вычисляем расстояние от центра окружности до ближайшей точки
	dx := c1.X - nearestX
	dy := c1.Y - nearestY

	dist := math.Sqrt(dx*dx + dy*dy)
	minDist := c1.Radius

	pushDist := minDist - dist
	// Если расстояние равно нулю, определяем направление толкания
	if dist == 0 {
		if c1.X < r1.X {
			dx = -1
		} else if c1.X > r1.X+r1.Width {
			dx = 1
		}
		if c1.Y < r1.Y {
			dy = -1
		} else if c1.Y > r1.Y+r1.Height {
			dy = 1
		}
		dist = 1 // Избегаем деления на ноль
	}

	// Перемещение окружности на границу столкновения
	c1.X += (dx / dist) * pushDist
	c1.Y += (dy / dist) * pushDist
}

func normalToRect(c1 util.Vec2, r1 RectShape) util.Vec2 {
	// Находим ближайшую точку прямоугольника к центру окружности
	nearestX := math.Max(r1.X, math.Min(c1.X, r1.X+r1.Width))
	nearestY := math.Max(r1.Y, math.Min(c1.Y, r1.Y+r1.Height))

	// Вычисляем расстояние от центра окружности до ближайшей точки
	dx := c1.X - nearestX
	dy := c1.Y - nearestY

	dist := math.Sqrt(dx*dx + dy*dy)

	// Если расстояние равно нулю, определяем направление толкания
	if dist == 0 {
		if c1.X < r1.X {
			dx = -1
		} else if c1.X > r1.X+r1.Width {
			dx = 1
		}
		if c1.Y < r1.Y {
			dy = -1
		} else if c1.Y > r1.Y+r1.Height {
			dy = 1
		}
		dist = 1 // Избегаем деления на ноль
	}

	// Перемещение окружности на границу столкновения
	return util.NewVec2((dx / dist), (dy / dist))
}

func normalToCircle(p util.Vec2, c1 CircleShape) util.Vec2 {
	normalX := p.X - c1.X
	normalY := p.Y - c1.Y
	dist := math.Sqrt(normalX*normalX + normalY*normalY)
	if dist != 0 {
		normalX /= dist
		normalY /= dist
	}
	return util.NewVec2(normalX, normalY)
}
