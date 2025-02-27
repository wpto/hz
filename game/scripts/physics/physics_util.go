package physics

import "iter"

func (p *Physics) iterateCells(lower Cell, upper Cell) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := lower.i; i <= upper.i; i++ {
			if p.grid[i] == nil {
				continue
			}

			for j := lower.j; j <= lower.j; j++ {
				for id := range p.grid[i][j] {
					if !yield(id) {
						return
					}
				}
			}
		}
	}
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
