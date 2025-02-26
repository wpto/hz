package physics

import (
	"hz/game/core"
	"hz/resources/images"
	"iter"

	"github.com/hajimehoshi/ebiten/v2"
)

type LineShape struct {
	_              Shape
	X1, Y1, X2, Y2 float64
}

func (LineShape) ShapeType() {}

func (m LineShape) CellIter(p *Physics) iter.Seq[Cell] {
	i1, j1 := p.GetCell(m.X1, m.Y1)
	i2, j2 := p.GetCell(m.X2, m.Y2)
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
func (m LineShape) CellBounds(p *Physics) (Cell, Cell) {
	i1, j1 := p.GetCell(m.X1, m.Y1)
	i2, j2 := p.GetCellUpper(m.X2, m.Y2)
	return Cell{i1, j1}, Cell{i2, j2}
}

func (m LineShape) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(m.X1, m.Y1)
	op.GeoM.Scale((m.X2-m.X1)/32, (m.Y2-m.Y1)/32)

	core.GlobalCamera.ApplyTransform(&op.GeoM)

	screen.DrawImage(images.DebugLine.Image, op)
}

func (m LineShape) Move(dx, dy float64) {
	m.X1 += dx
	m.Y1 += dy
	m.X2 += dx
	m.Y2 += dy
}
