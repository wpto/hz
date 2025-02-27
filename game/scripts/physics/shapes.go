package physics

import (
	"iter"

	"github.com/hajimehoshi/ebiten/v2"
)

type Shape interface {
	ShapeType()
	CellIter(p *Physics) iter.Seq[Cell]
	CellBounds(p *Physics) (Cell, Cell)
	Move(dx, dy float64) Shape

	Draw(screen *ebiten.Image)
}
