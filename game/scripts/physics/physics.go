package physics

import (
	"fmt"
	"hz/game/core"
	"hz/game/util"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type PhysicsBackwardPusher interface {
	PhysicsBackwardPush(shape Shape)
}

type PhysicsUpdate interface {
	PhysicsUpdate(dt float64)
}

type bodyObj struct {
	id    int
	shape Shape
	owner any
}

type Physics struct {
	shapes   map[int]bodyObj
	grid     map[int]map[int]map[int]struct{}
	cellSize int

	counter int
}

func NewPhysics() *Physics {
	return &Physics{
		shapes:   make(map[int]bodyObj),
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
	p.shapes[id] = bodyObj{
		id:    id,
		shape: sh,
		owner: owner,
	}
	p.MarkShape(id, sh)
	return id
}

func (p *Physics) ForwardPush(id int, sh Shape) {
	p.UnmarkShape(id, p.shapes[id].shape)
	p.shapes[id] = bodyObj{
		id:    p.shapes[id].id,
		shape: sh,
		owner: p.shapes[id].owner,
	}
	p.MarkShape(id, sh)
}

func (p *Physics) Draw(screen *ebiten.Image) {
	for _, shape := range p.shapes {
		shape.shape.Draw(screen)
	}
}

const substeps = 3

func (p *Physics) Update() {
	const subdt = core.Delta / substeps
	for range substeps {
		for id := range p.shapes {
			upd, ok := p.shapes[id].owner.(PhysicsUpdate)
			if ok {
				upd.PhysicsUpdate(subdt)
			}
		}

		p.UpdateIteration()

		for id := range p.shapes {
			upd, ok := p.shapes[id].owner.(PhysicsBackwardPusher)
			if ok {
				upd.PhysicsBackwardPush(p.shapes[id].shape)
			}
		}
	}
}

func (p *Physics) MoveAndCollide(id int, dx, dy float64, substeps int) {
	if substeps == 0 {
		substeps = 4
	}
	subd := 1.0 / float64(substeps)
	origShape := p.shapes[id].shape
	tempShape := p.shapes[id].shape
	i := 0

	for ; i < substeps; i++ {
		newTempShape := p.shapes[id].shape.Move(dx*subd, dy*subd)
		otherObjID, found := p.CheckShapeCollision(bodyObj{
			id:    id,
			shape: tempShape,
		})
		if !found {
			tempShape = newTempShape
		} else {
			fmt.Println("move and collide: id", id, "collided with", otherObjID)
			break
		}

		if i > 0 {
			fmt.Println("move and collide update :)")
			p.UnmarkShape(id, origShape)
			p.shapes[id] = bodyObj{
				id:    id,
				shape: tempShape,
				owner: p.shapes[id].owner,
			}
			p.MarkShape(id, tempShape)
		}
	}
}

func (p *Physics) MoveAndSlide(id int, dx, dy float64, substeps int) {
	if substeps == 0 {
		substeps = 4
	}
	subd := 1.0 / float64(substeps)
	origShape := p.shapes[id].shape
	tempShape := p.shapes[id].shape
	offsetX, offsetY := 0.0, 0.0
	i := 0

	fmt.Println("move and slide:", id, util.NewVec2(dx, dy).Len())
	for ; i < substeps; i++ {
		newOffsetX, newOffsetY := offsetX+dx*subd, offsetY+dy*subd
		newTempShape := p.shapes[id].shape.Move(newOffsetX, newOffsetY)
		otherObjID, found := p.CheckShapeCollision(bodyObj{
			id:    id,
			shape: tempShape,
		})
		if !found {
			tempShape = newTempShape
			offsetX, offsetY = newOffsetX, newOffsetY
			fmt.Println("move and slide: ", id, "offset", util.NewVec2(dx*subd, dy*subd).Len())
		} else {
			x, y := float64(0), float64(0)
			switch tempShape := tempShape.(type) {
			case CircleShape:
				x, y = tempShape.X, tempShape.Y
			case RectShape:
				x, y = tempShape.X, tempShape.Y
			}

			normal := util.Vec2{X: 0, Y: 0}

			switch otherShape := p.shapes[otherObjID].shape.(type) {
			case CircleShape:
				normal = normalToCircle(util.Vec2{X: x, Y: y}, otherShape)
				fmt.Println("normal to circle", normal)
			case RectShape:
				normal = normalToRect(util.Vec2{X: x, Y: y}, otherShape)
				fmt.Println("normal to rect")
			}
			dot := normal.Dot(util.Vec2{X: dx, Y: dy})
			oldDelta := util.NewVec2(dx, dy)
			dx -= dot * normal.X
			dy -= dot * normal.Y
			fmt.Println("move and slide: id", id, "collide: old dx", oldDelta, "new dx", util.NewVec2(dx, dy))
		}
	}

	if i > 0 {
		fmt.Println("move and slide update :)")
		p.UnmarkShape(id, origShape)
		p.shapes[id] = bodyObj{
			id:    id,
			shape: tempShape,
			owner: p.shapes[id].owner,
		}
		p.MarkShape(id, tempShape)
	}
}

func (p *Physics) CheckShapeCollision(thisObj bodyObj) (int, bool) {
	for otherObjID := range p.iterateCells(thisObj.shape.CellBounds(p)) {
		if thisObj.id == otherObjID {
			continue
		}

		thisShape := thisObj.shape
		otherShape := p.shapes[otherObjID].shape

		switch s1 := thisShape.(type) {
		case CircleShape:
			switch s2 := otherShape.(type) {
			case CircleShape:
				if isCirclesCollideThreshold(s1, s2, 0.01) {
					return otherObjID, true
				}
			case RectShape:
				if isCircleRectangleCollisionThreshold(s1, s2, 0.01) {
					return otherObjID, true
				}
			}
		case RectShape:
			switch s2 := otherShape.(type) {
			case CircleShape:
				if isCircleRectangleCollisionThreshold(s2, s1, 0.01) {
					return otherObjID, true
				}
			}
		}
	}
	return 0, false
}

func (p *Physics) UpdateIteration() {
	for _, thisObj := range p.shapes {
		for otherObjID := range p.iterateCells(thisObj.shape.CellBounds(p)) {

			if thisObj.id == otherObjID {
				continue
			}
			otherShape := p.shapes[otherObjID]

			switch s1 := thisObj.shape.(type) {
			case CircleShape:
				switch s2 := otherShape.shape.(type) {
				case CircleShape:
					if isCirclesCollide(s1, s2) {
						p.UnmarkShape(thisObj.id, thisObj.shape)
						p.UnmarkShape(otherShape.id, otherShape.shape)
						solveCircleCollision(&s1, &s2)
						p.shapes[thisObj.id] = bodyObj{
							id:    thisObj.id,
							shape: s1,
							owner: thisObj.owner,
						}
						p.shapes[otherShape.id] = bodyObj{
							id:    otherShape.id,
							shape: s2,
							owner: otherShape.owner,
						}
						p.MarkShape(thisObj.id, thisObj.shape)
						p.MarkShape(otherShape.id, otherShape.shape)
					}
				case RectShape:
					if isCircleRectangleCollision(s1, s2) {
						p.UnmarkShape(thisObj.id, thisObj.shape)
						p.UnmarkShape(otherShape.id, otherShape.shape)
						solveCircleRectCollision(&s1, &s2)
						p.shapes[thisObj.id] = bodyObj{
							id:    thisObj.id,
							shape: s1,
							owner: thisObj.owner,
						}
						p.shapes[otherShape.id] = bodyObj{
							id:    otherShape.id,
							shape: s2,
							owner: otherShape.owner,
						}
						p.MarkShape(thisObj.id, thisObj.shape)
						p.MarkShape(otherShape.id, otherShape.shape)
					}
				}
			case RectShape:
				switch s2 := otherShape.shape.(type) {
				case CircleShape:
					if isCircleRectangleCollision(s2, s1) {
						p.UnmarkShape(thisObj.id, thisObj.shape)
						p.UnmarkShape(otherShape.id, otherShape.shape)
						solveCircleRectCollision(&s2, &s1)
						p.shapes[thisObj.id] = bodyObj{
							id:    thisObj.id,
							shape: s1,
							owner: thisObj.owner,
						}
						p.shapes[otherShape.id] = bodyObj{
							id:    otherShape.id,
							shape: s2,
							owner: otherShape.owner,
						}
						p.MarkShape(thisObj.id, thisObj.shape)
						p.MarkShape(otherShape.id, otherShape.shape)
					}
				}
			}
		}
	}
}
