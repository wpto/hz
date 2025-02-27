package physics

import (
	"hz/game/core"
	"hz/resources/images"
	"iter"

	"github.com/hajimehoshi/ebiten/v2"
)

type RectShape struct {
	_                   Shape
	X, Y, Width, Height float64
}

func (RectShape) ShapeType() {}

type Cell struct {
	i, j int
}

func (m RectShape) CellIter(p *Physics) iter.Seq[Cell] {
	i1, j1 := p.GetCell(m.X, m.Y)
	i2, j2 := p.GetCellUpper(m.X+m.Width, m.Y+m.Height)
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

func (m RectShape) CellBounds(p *Physics) (Cell, Cell) {
	i1, j1 := p.GetCell(m.X, m.Y)
	i2, j2 := p.GetCellUpper(m.X+m.Width, m.Y+m.Height)
	return Cell{i1, j1}, Cell{i2, j2}
}

func (m RectShape) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(-m.Width/2, -m.Height/2)
	op.GeoM.Scale(m.Width/32, m.Height/32)
	op.GeoM.Translate(m.X, m.Y)

	core.GlobalCamera.ApplyTransform(&op.GeoM)

	screen.DrawImage(images.DebugRect.Image, op)
}

func (m RectShape) Move(dx, dy float64) Shape {
	m.X += dx
	m.Y += dy
	return m
}
