package physics

import (
	"hz/game/core"
	"hz/resources/images"
	"iter"

	"github.com/hajimehoshi/ebiten/v2"
)

type CircleShape struct {
	X, Y, Radius float64
}

func (CircleShape) ShapeType() {}

func (m CircleShape) CellIter(p *Physics) iter.Seq[Cell] {
	i1, j1 := p.GetCell(m.X-m.Radius, m.Y-m.Radius)
	i2, j2 := p.GetCellUpper(m.X+m.Radius, m.Y+m.Radius)
	return func(yield func(Cell) bool) {
		for i := i1; i <= i2; i++ {
			for j := j1; j <= j2; j++ {
				if !yield(Cell{i, j}) {
					return
				}
			}
		}
	}
}

func (m CircleShape) CellBounds(p *Physics) (Cell, Cell) {
	i1, j1 := p.GetCell(m.X-m.Radius, m.Y-m.Radius)
	i2, j2 := p.GetCellUpper(m.X+m.Radius, m.Y+m.Radius)
	return Cell{i1, j1}, Cell{i2, j2}
}

func (m CircleShape) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-16, -16)
	op.GeoM.Scale(m.Radius/16.0, m.Radius/16.0)
	op.GeoM.Translate(m.X, m.Y)
	// op.GeoM.Scale(m.Radius/32.0, m.Radius/32.0)
	core.GlobalCamera.ApplyTransform(&op.GeoM)

	screen.DrawImage(images.DebugCircle.Image, op)
}

func (m CircleShape) Move(dx, dy float64) Shape {
	m.X += dx
	m.Y += dy
	return m
}
