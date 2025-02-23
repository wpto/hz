package core

import (
	"hz/resources/images"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type AnimatedSprite struct {
	layer    string
	img      images.Image
	duration float64

	// update
	currentFrame image.Rectangle
	op           *ebiten.DrawImageOptions
	t            float64

	// setter
	px, py   float64
	rotation float64
	speed    float64
}

func NewAnimatedSprite(img images.Image, layer string, duration float64) *AnimatedSprite {
	return &AnimatedSprite{
		layer:    layer,
		img:      img,
		duration: duration,
		speed:    1,
	}
}

func (a *AnimatedSprite) SetPosition(x, y float64) {
	a.px, a.py = x, y
}

func (a *AnimatedSprite) SetRotation(r float64) {
	a.rotation = r
}

func (a *AnimatedSprite) SetAnimationSpeed(speed float64) {
	a.speed = speed
	if speed < 0.001 {
		a.speed = 0
		a.t = 0
	}
}

func (a *AnimatedSprite) Update() {
	if a.speed > 0 {
		a.t += Delta * a.speed
		a.t = math.Mod(a.t, a.duration)
	}

	layer := a.img.Layers[a.layer]
	i := int((a.t / a.duration) * float64(layer.FrameCount))
	a.currentFrame = image.Rect(
		layer.OX+i*a.img.Width, layer.OY,
		layer.OX+(i+1)*a.img.Width, layer.OY+a.img.Height)

	a.op = &ebiten.DrawImageOptions{}
	a.op.GeoM.Translate(-float64(a.img.Width)/2, -float64(a.img.Height)/2)
	a.op.GeoM.Rotate(a.rotation)
	a.op.GeoM.Translate(float64(a.px), float64(a.py))
	GlobalCamera.ApplyTransform(&a.op.GeoM)

}

func (a *AnimatedSprite) Draw(screen *ebiten.Image) {
	screen.DrawImage(
		a.img.Image.SubImage(a.currentFrame).(*ebiten.Image), a.op)
}
