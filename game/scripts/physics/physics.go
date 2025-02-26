package physics

import (
	"fmt"
	"hz/game/core"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type PhysicsBackwardPusher interface {
	PhysicsBackwardPush(shape Shape)
}

type PhysicsUpdate interface {
	PhysicsUpdate(dt float64)
}

type shape struct {
	id    int
	shape Shape
	owner any
}

type Physics struct {
	shapes   map[int]shape
	grid     map[int]map[int]map[int]struct{}
	cellSize int

	counter int
}

func NewPhysics() *Physics {
	return &Physics{
		shapes:   make(map[int]shape),
		grid:     make(map[int]map[int]map[int]struct{}),
		cellSize: 20,
		counter:  0,
	}
}

func (p *Physics) GetCell(x, y float64) (int, int) {
	return int(x) / p.cellSize, int(y) / p.cellSize
}

func (p *Physics) GetCellUpper(x, y float64) (int, int) {
	return int(math.Ceil(x / float64(p.cellSize))), int(math.Ceil(y / float64(p.cellSize)))
}

func (p *Physics) AddShape(sh Shape, owner any) int {
	id := p.counter
	p.counter++

	p.shapes[id] = shape{
		id:    id,
		shape: sh,
		owner: owner,
	}
	p.MarkShape(id, sh)
	return id
}

func (p *Physics) ForwardPush(id int, sh Shape) {
	p.UnmarkShape(id, p.shapes[id].shape)
	p.shapes[id] = shape{
		id:    p.shapes[id].id,
		shape: sh,
		owner: p.shapes[id].owner,
	}
	p.MarkShape(id, sh)
}

func (p *Physics) MarkShape(id int, shape Shape) {
	for cell := range shape.CellIter(p) {
		i, j := cell.i, cell.j
		if p.grid[i] == nil {
			p.grid[i] = make(map[int]map[int]struct{})
		}
		if p.grid[i][j] == nil {
			p.grid[i][j] = make(map[int]struct{})
		}

		p.grid[i][j][id] = struct{}{}
	}
}

func (p *Physics) UnmarkShape(id int, shape Shape) {
	for cell := range shape.CellIter(p) {
		i, j := cell.i, cell.j
		delete(p.grid[i][j], id)
	}
}

func (p *Physics) CheckLine(x1, y1, x2, y2 float64) bool {
	// for _, circ := range p.Circles {
	// 	if isIntersectingSegmentCircle(x1, y1, x2, y2, circ.X, circ.Y, circ.Radius) {
	// 		return true
	// 	}
	// }

	// for _, rect := range p.Rects {
	// 	if isIntersectingSegmentRectangle(x1, y1, x2, y2, rect.X, rect.Y, rect.X+rect.Width, rect.Y+rect.Height) {
	// 		return true
	// 	}
	// }

	return false
}

func (p *Physics) Draw(screen *ebiten.Image) {
	for _, shape := range p.shapes {
		shape.shape.Draw(screen)
	}
}

const substeps = 3

func (p *Physics) Update() {
	const subdt = core.Delta / substeps
	for i := 0; i < substeps; i++ {
		for id := range p.shapes {
			upd, ok := p.shapes[id].owner.(PhysicsUpdate)
			if ok {
				upd.PhysicsUpdate(subdt)
			}
		}

		fmt.Println("Collision check")
		p.UpdateIteration(0.5)
		fmt.Println("Collision done")

		for id := range p.shapes {
			upd, ok := p.shapes[id].owner.(PhysicsBackwardPusher)
			if ok {
				upd.PhysicsBackwardPush(p.shapes[id].shape)
			}
		}
	}
}

func (p *Physics) MoveAndCollide(id int, dx, dy float64) {
	lower, upper := p.shapes[id].shape.CellBounds(p)

	newShape := p.shapes[id].shape
	newShape
	for i := lower.i; i <= upper.i; i++ {
		if p.grid[i] == nil {
			continue
		}

		for j := lower.j; j <= lower.j; j++ {
			for otherID := range p.grid[i][j] {
				if id == otherID {
					continue
				}

				// if isIntersectingCircleCircle(p.Circles[id], p.Circles[otherID]) {
				// 	solveCircleCollision(&p.Circles[id], &p.Circles[otherID])
				// }
			}
		}
	}
}

func (p *Physics) UpdateIteration(portion float64) {
	for _, sh := range p.shapes {
		lower, upper := sh.shape.CellBounds(p)

		for i := lower.i; i <= upper.i; i++ {
			if p.grid[i] == nil {
				continue
			}

			for j := lower.j; j <= lower.j; j++ {
				for id := range p.grid[i][j] {
					if sh.id == id {
						continue
					}
					otherShape := p.shapes[id]

					switch s1 := sh.shape.(type) {
					case CircleShape:
						switch s2 := otherShape.shape.(type) {
						case CircleShape:
							if isCirclesCollide(s1, s2) {
								p.UnmarkShape(sh.id, sh.shape)
								p.UnmarkShape(otherShape.id, otherShape.shape)
								solveCircleCollision(&s1, &s2, portion)
								p.shapes[sh.id] = shape{
									id:    sh.id,
									shape: s1,
									owner: sh.owner,
								}
								p.shapes[otherShape.id] = shape{
									id:    otherShape.id,
									shape: s2,
									owner: otherShape.owner,
								}
								p.MarkShape(sh.id, sh.shape)
								p.MarkShape(otherShape.id, otherShape.shape)
							}
						case RectShape:
							if isCircleRectangleCollision(s1, s2) {
								p.UnmarkShape(sh.id, sh.shape)
								p.UnmarkShape(otherShape.id, otherShape.shape)
								solveCircleRectCollision(&s1, &s2, portion)
								p.shapes[sh.id] = shape{
									id:    sh.id,
									shape: s1,
									owner: sh.owner,
								}
								p.shapes[otherShape.id] = shape{
									id:    otherShape.id,
									shape: s2,
									owner: otherShape.owner,
								}
								p.MarkShape(sh.id, sh.shape)
								p.MarkShape(otherShape.id, otherShape.shape)
							}
						}
					case RectShape:
						switch s2 := otherShape.shape.(type) {
						case CircleShape:
							if isCircleRectangleCollision(s2, s1) {
								p.UnmarkShape(sh.id, sh.shape)
								p.UnmarkShape(otherShape.id, otherShape.shape)
								solveCircleRectCollision(&s2, &s1, portion)
								p.shapes[sh.id] = shape{
									id:    sh.id,
									shape: s1,
									owner: sh.owner,
								}
								p.shapes[otherShape.id] = shape{
									id:    otherShape.id,
									shape: s2,
									owner: otherShape.owner,
								}
								p.MarkShape(sh.id, sh.shape)
								p.MarkShape(otherShape.id, otherShape.shape)
							}
						}
					}
				}
			}
		}
	}

}

func isCirclesCollide(c1, c2 CircleShape) bool {
	dx := c1.X - c2.X
	dy := c1.Y - c2.Y
	r := c1.Radius + c2.Radius
	return dx*dx+dy*dy <= r*r
}

// Проверяет пересечение окружности и прямоугольника
func isCircleRectangleCollision(c1 CircleShape, r1 RectShape) bool {
	// Находим ближайшую точку прямоугольника к центру окружности
	nearestX := math.Max(r1.X, math.Min(c1.X, r1.X+r1.Width))
	nearestY := math.Max(r1.Y, math.Min(c1.Y, r1.Y+r1.Height))

	// Вычисляем расстояние от центра окружности до ближайшей точки
	dx := c1.X - nearestX
	dy := c1.Y - nearestY

	distanceSquared := dx*dx + dy*dy

	return distanceSquared < (c1.Radius * c1.Radius)
}

func solveCircleCollision(c1, c2 *CircleShape, portion float64) {
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

	c1.X += nx * delta * portion
	c1.Y += ny * delta * portion
	c2.X -= nx * delta * portion
	c2.Y -= ny * delta * portion

	// newDist := math.Sqrt((c1.X-c2.X)*(c1.X-c2.X) + (c1.Y-c2.Y)*(c1.Y-c2.Y))
	// if newDist+0.001 < minDist {
	// 	log.Panic("collision not resolved")
	// }
}

func solveCircleRectCollision(c1 *CircleShape, r1 *RectShape, portion float64) {
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
	c1.X += (dx / dist) * pushDist * portion
	c1.Y += (dy / dist) * pushDist * portion
}
