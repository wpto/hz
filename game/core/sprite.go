package core

import (
	"hz/resources/images"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	img images.Image
	op  *ebiten.DrawImageOptions

	px, py   float64
	rotation float64
}

func NewSprite(img images.Image) *Sprite {
	return &Sprite{
		img: img,
	}
}

func (s *Sprite) SetPosition(x, y float64) {
	s.px, s.py = x, y
}

func (s *Sprite) SetRotation(r float64) {
	s.rotation = r
}

func (s *Sprite) Draw(screen *ebiten.Image) {

	screen.DrawImage(s.img.Image, s.op)
}

func (s *Sprite) Update() {
	s.op = &ebiten.DrawImageOptions{}
	s.op.GeoM.Translate(-float64(s.img.Width)/2, -float64(s.img.Height)/2)
	s.op.GeoM.Rotate(s.rotation)
	s.op.GeoM.Translate(s.px, s.py)
	GlobalCamera.ApplyTransform(&s.op.GeoM)
}
