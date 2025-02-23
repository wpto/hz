package core

import "github.com/hajimehoshi/ebiten/v2"

const (
	ScreenWidth  = 480 / 2
	ScreenHeight = 360 / 2
)

type Camera struct {
	X, Y  float64
	Scale float64
}

var GlobalCamera Camera = Camera{0, 0, 1}

func (c *Camera) SetPosition(x, y float64) {
	c.X, c.Y = x, y
}

func (c *Camera) SetScale(scale float64) {
	c.Scale = scale
}

func (c *Camera) ScreenToWorld(x float64, y float64) (float64, float64) {
	op := ebiten.GeoM{}
	c.ApplyTransform(&op)
	op.Invert()

	x, y = op.Apply(x, y)
	return x, y
}

func (c *Camera) WorldToCamera(x float64, y float64) (float64, float64) {
	op := ebiten.GeoM{}
	c.ApplyTransform(&op)
	return op.Apply(x, y)
}

func (c *Camera) ApplyTransform(geom *ebiten.GeoM) {
	geom.Translate(-c.X, -c.Y)
	geom.Translate(ScreenWidth/2, ScreenHeight/2)
	geom.Scale(c.Scale, c.Scale)
}

func (c *Camera) GetMouseWorldCoordinates() (float64, float64) {
	x, y := ebiten.CursorPosition()
	// x, y = x-ScreenWidth/2, y-ScreenHeight/2
	return c.ScreenToWorld(float64(x), float64(y))
}
