package scripts

import (
	"hz/game/core"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type PhysicsUpdater interface {
	SetPosition(x, y float64)
}

type PhysicsRect struct {
	ID                  int
	X, Y, Width, Height float64
}

type PhysicsCircle struct {
	ID           int
	X, Y, Radius float64
}

type PhysicsShape interface {
}

type PhysicsUpdate interface {
	PhysicsUpdate(dt float64)
}

type Physics struct {
	Rects         map[int]PhysicsRect
	Circles       map[int]PhysicsCircle
	CircleUpdater map[int]PhysicsUpdater
	grid          map[int]map[int]map[int]struct{}
	cellSize      int

	Updaters []PhysicsUpdate

	counter int
}

func NewPhysics() *Physics {
	return &Physics{
		Rects:         make(map[int]PhysicsRect),
		Circles:       make(map[int]PhysicsCircle),
		CircleUpdater: make(map[int]PhysicsUpdater),
		grid:          make(map[int]map[int]map[int]struct{}),
		cellSize:      20,
		counter:       0,
	}
}

func (p *Physics) AddPhysicsUpdater(updater PhysicsUpdate) {
	p.Updaters = append(p.Updaters, updater)
}

func (p *Physics) GetCell(x, y float64) (int, int) {
	return int(x) / p.cellSize, int(y) / p.cellSize
}

func (p *Physics) GetCellUpper(x, y float64) (int, int) {
	return int(math.Ceil(x / float64(p.cellSize))), int(math.Ceil(y / float64(p.cellSize)))
}

func (p *Physics) MarkCellRect(x, y, w, h float64, id int) {
	i1, j1 := p.GetCell(x, y)
	i2, j2 := p.GetCellUpper(x+w, y+h)
	for i := i1; i <= i2; i++ {
		for j := j1; j <= j2; j++ {
			if p.grid[i] == nil {
				p.grid[i] = make(map[int]map[int]struct{})
			}
			if p.grid[i][j] == nil {
				p.grid[i][j] = make(map[int]struct{})
			}

			p.grid[i][j][id] = struct{}{}
		}
	}
}

func (p *Physics) UnmarkCellRect(x, y, w, h float64, id int) {
	i1, j1 := p.GetCell(x, y)
	i2, j2 := p.GetCellUpper(x+w, y+h)
	for i := i1; i <= i2; i++ {
		for j := j1; j <= j2; j++ {
			delete(p.grid[i][j], id)
		}
	}
}

func (p *Physics) AddRect(x, y, w, h float64) int {
	id := p.counter
	p.counter++
	p.Rects[id] = PhysicsRect{ID: id, X: x, Y: y, Width: w, Height: h}
	p.MarkCellRect(x, y, w, h, id)
	return id
}

func (p *Physics) UpdateRect(id int, x, y, w, h float64) {
	p.UnmarkCellRect(p.Rects[id].X, p.Rects[id].Y, p.Rects[id].Width, p.Rects[id].Height, id)
	p.Rects[id] = PhysicsRect{ID: p.Rects[id].ID, X: x, Y: y, Width: w, Height: h}
	p.MarkCellRect(x, y, w, h, id)
}

func (p *Physics) AddCircle(x, y, r float64, updater PhysicsUpdater) int {
	id := p.counter
	p.counter++
	p.Circles[id] = PhysicsCircle{ID: id, X: x, Y: y, Radius: r}
	p.CircleUpdater[id] = updater
	p.MarkCellRect(x-r, y-r, x+r, y+r, id)
	return id
}

func (p *Physics) UpdateCircle(id int, x, y, r float64) {
	p.UnmarkCellRect(p.Circles[id].X-p.Circles[id].Radius, p.Circles[id].Y-p.Circles[id].Radius, p.Circles[id].Radius*2, p.Circles[id].Radius*2, id)
	p.Circles[id] = PhysicsCircle{ID: id, X: x, Y: y, Radius: r}
	p.MarkCellRect(x-r, y-r, x+r, y+r, id)
}

func (p *Physics) CheckLine(x1, y1, x2, y2 float64) bool {
	for _, circ := range p.Circles {
		if isIntersectingSegmentCircle(x1, y1, x2, y2, circ.X, circ.Y, circ.Radius) {
			return true
		}
	}

	// for _, rect := range p.Rects {
	// 	if isIntersectingSegmentRectangle(x1, y1, x2, y2, rect.X, rect.Y, rect.X+rect.Width, rect.Y+rect.Height) {
	// 		return true
	// 	}
	// }

	return false
}

func (p *Physics) Draw(screen *ebiten.Image) {
	for _, circ := range p.Circles {
		for i := -circ.Radius; i < circ.Radius; i += 1 {
			for j := -circ.Radius; j < +circ.Radius; j += 1 {
				if i*i+j*j < circ.Radius*circ.Radius {
					x, y := core.GlobalCamera.WorldToCamera(i+circ.X, circ.Y+j)
					screen.Set(int(x), int(y), color.RGBA{128, 0, 0, 128})
				}
			}
		}
	}

	for _, rect := range p.Rects {
		for i := 0; i < int(rect.Width); i++ {
			for j := 0; j < int(rect.Height); j++ {
				x, y := core.GlobalCamera.WorldToCamera(rect.X+float64(i), rect.Y+float64(j))
				screen.Set(int(x), int(y), color.RGBA{0, 128, 0, 128})
			}
		}
	}
}

const substeps = 8

func (p *Physics) Update() {
	const subdt = core.Delta / substeps
	for i := 0; i < 8; i++ {
		p.UpdateIteration()
		for _, updater := range p.Updaters {
			updater.PhysicsUpdate(subdt)
		}
	}
}

func (p *Physics) UpdateIteration() {
	for _, circ := range p.Circles {
		i, j := p.GetCell(circ.X-circ.Radius, circ.Y-circ.Radius)
		i2, j2 := p.GetCell(circ.X+circ.Radius, circ.Y+circ.Radius)

		for i := i; i <= i2; i++ {
			for j := j; j <= j2; j++ {
				if p.grid[i] == nil {
					continue
				}
				for id := range p.grid[i][j] {
					if id == circ.ID {
						continue
					}

					useCircle := true
					otherCircle, ok := p.Circles[id]
					if !ok {
						useCircle = false

					}

					if useCircle {
						if isCirclesCollide(circ, otherCircle) {
							solveCircleCollision(&circ, &otherCircle)
							p.CircleUpdater[circ.ID].SetPosition(circ.X, circ.Y)
							p.CircleUpdater[otherCircle.ID].SetPosition(otherCircle.X, otherCircle.Y)
						}
					}

					useRect := true
					otherRect, ok := p.Rects[id]
					if !ok {
						useRect = false
					}

					if useRect {
						if isCircleRectangleCollision(circ, otherRect) {
							solveCircleRectCollision(&circ, &otherRect)
							p.CircleUpdater[circ.ID].SetPosition(circ.X, circ.Y)
						}
					}
				}
			}
		}
	}
}

func isCirclesCollide(c1, c2 PhysicsCircle) bool {
	dx := c1.X - c2.X
	dy := c1.Y - c2.Y
	r := c1.Radius + c2.Radius
	return dx*dx+dy*dy <= r*r
}

// Проверяет пересечение окружности и прямоугольника
func isCircleRectangleCollision(c1 PhysicsCircle, r1 PhysicsRect) bool {
	// Находим ближайшую точку прямоугольника к центру окружности
	nearestX := math.Max(r1.X, math.Min(c1.X, r1.X+r1.Width))
	nearestY := math.Max(r1.Y, math.Min(c1.Y, r1.Y+r1.Height))

	// Вычисляем расстояние от центра окружности до ближайшей точки
	dx := c1.X - nearestX
	dy := c1.Y - nearestY

	distanceSquared := dx*dx + dy*dy

	return distanceSquared < (c1.Radius * c1.Radius)
}

func solveCircleCollision(c1, c2 *PhysicsCircle) {
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

func solveCircleRectCollision(c1 *PhysicsCircle, r1 *PhysicsRect) {
	// Находим ближайшую точку прямоугольника к центру окружности
	nearestX := math.Max(r1.X, math.Min(c1.X, r1.X+r1.Width))
	nearestY := math.Max(r1.Y, math.Min(c1.Y, r1.Y+r1.Height))
	// Вычисляем расстояние от центра окружности до ближайшей точки
	dx := c1.X - nearestX
	dy := c1.Y - nearestY

	distanceSquared := math.Sqrt(dx*dx + dy*dy)
	minDist := c1.Radius

	delta := (minDist - distanceSquared)

	c1.X += dx * delta
	c1.Y += dy * delta
}
