package scripts

import (
	"hz/game/core"
	"hz/game/util"
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

type Physics struct {
	Rects         map[int]PhysicsRect
	Circles       map[int]PhysicsCircle
	CircleUpdater map[int]PhysicsUpdater
	grid          map[int]map[int]map[int]struct{}
	cellSize      int
}

func NewPhysics() *Physics {
	return &Physics{
		Rects:         make(map[int]PhysicsRect),
		Circles:       make(map[int]PhysicsCircle),
		CircleUpdater: make(map[int]PhysicsUpdater),
		grid:          make(map[int]map[int]map[int]struct{}),
		cellSize:      32,
	}
}

func (p *Physics) GetCell(x, y float64) (int, int) {
	return int(x) / p.cellSize, int(y) / p.cellSize
}

func (p *Physics) MarkCellRect(x, y, w, h float64, id int) {
	i1, j1 := p.GetCell(x, y)
	i2, j2 := p.GetCell(x+w, y+h)
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
	i2, j2 := p.GetCell(x+w, y+h)
	for i := i1; i <= i2; i++ {
		for j := j1; j <= j2; j++ {
			delete(p.grid[i][j], id)
		}
	}
}

func (p *Physics) AddRect(x, y, w, h float64) int {
	id := len(p.Rects)
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
	id := len(p.Circles)
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
}

func (p *Physics) Update() {
	for i := 0; i < 8; i++ {
		p.UpdateIteration()
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

					other := p.Circles[id]
					if isCirclesCollide(circ, other) {
						solveCircleCollision(&circ, &other)
						p.CircleUpdater[circ.ID].SetPosition(circ.X, circ.Y)
						p.CircleUpdater[other.ID].SetPosition(other.X, other.Y)
					}
				}
			}
		}
	}
}

func isCirclesCollide(c1, c2 PhysicsCircle) bool {
	dx := c1.X - c2.X
	dy := c1.Y - c2.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	return dist < c1.Radius+c2.Radius
}

func solveCircleCollision(c1, c2 *PhysicsCircle) {
	dx := c1.X - c2.X
	dy := c1.Y - c2.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist <= .1 {

	}

	// Нормализация вектора
	dx /= dist
	dy /= dist

	// Перемещение кругов

	cdx, cdy := util.NewVec2(dx, dy).Mul(c1.Radius - dist/2).Values()
	c1.X += cdx
	c1.Y += cdy

	cdx, cdy = util.NewVec2(dx, dy).Mul(c2.Radius - dist/2).Values()
	c2.X -= cdx
	c2.Y -= cdy
}
