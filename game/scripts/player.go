package scripts

import (
	"hz/game/core"
	"hz/game/util"
	"hz/resources/images"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	speed              = 200
	baseAnimationSpeed = speed * core.Delta
)

type Player struct {
	prevx, prevy float64
	X, Y         float64

	ChelBody, ChelLegs *core.AnimatedSprite

	crosshairX, crosshairY float64
	crosshairDirection     float64

	movingVelocity util.Vec2

	physics *Physics
	id      int
}

func NewPlayer(p *Physics) *Player {
	player := &Player{
		X:        0,
		Y:        0,
		ChelBody: core.NewAnimatedSprite(images.Chel, "body", 1.0),
		ChelLegs: core.NewAnimatedSprite(images.Chel, "legs", 0.5),
		physics:  p,
	}
	id := p.AddCircle(0, 0, 10, player)
	player.id = id

	return player
}

func (p *Player) Input() {
	var dx, dy float64
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy++
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx--
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx++
	}

	if dx != 0 && dy != 0 {
		dx *= math.Sqrt(2) / 2
		dy *= math.Sqrt(2) / 2
	}

	p.movingVelocity = util.NewVec2(dx, dy).Mul(speed)
}

func (p *Player) UpdateCrosshair() {
	mx, my := core.GlobalCamera.GetMouseWorldCoordinates()
	angle := math.Atan2(my-p.Y, mx-p.X)
	p.crosshairX, p.crosshairY = mx, my
	p.crosshairDirection = angle
}

func (p *Player) SetPosition(x, y float64) {
	p.prevx, p.prevy = p.X, p.Y
	p.X, p.Y = x, y
}

func (p *Player) GetPosition() (float64, float64) {
	return p.X, p.Y
}

func (p *Player) GetLookDirection() float64 {
	return p.crosshairDirection
}

func (p *Player) GetVelocitySqr() float64 {
	return (p.X-p.prevx)*(p.X-p.prevx) + (p.Y-p.prevy)*(p.Y-p.prevy)
}

func (p *Player) GetAnimationSpeed() float64 {
	return p.GetVelocitySqr() / baseAnimationSpeed / baseAnimationSpeed
}

func (p *Player) PhysicsUpdate(dt float64) {
	p.prevx, p.prevy = p.X, p.Y
	p.X += p.movingVelocity.X * dt
	p.Y += p.movingVelocity.Y * dt
}

func (p *Player) Update() error {
	p.UpdateCrosshair()
	p.physics.UpdateCircle(p.id, p.X, p.Y, 10)

	animationSpeed := p.GetAnimationSpeed()

	core.GlobalCamera.SetPosition(p.X, p.Y)
	p.ChelBody.SetPosition(p.X, p.Y)
	p.ChelBody.SetRotation(p.crosshairDirection)
	p.ChelBody.SetAnimationSpeed(animationSpeed)
	p.ChelBody.Update()

	movingDirection := math.Atan2(p.Y-p.prevy, p.X-p.prevx)
	p.ChelLegs.SetPosition(p.X, p.Y)
	p.ChelLegs.SetRotation(movingDirection)
	p.ChelLegs.SetAnimationSpeed(animationSpeed)
	p.ChelLegs.Update()
	return nil
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.ChelLegs.Draw(screen)
	p.ChelBody.Draw(screen)

	x, y := ebiten.CursorPosition()
	screen.Set(x, y, color.RGBA{255, 0, 0, 255})

	mx, my := core.GlobalCamera.GetMouseWorldCoordinates()
	back := ebiten.GeoM{}
	core.GlobalCamera.ApplyTransform(&back)
	mx, my = back.Apply(mx, my)
	screen.Set(int(mx), int(my), color.RGBA{0, 255, 0, 255})
}
