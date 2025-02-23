package scripts

import (
	"hz/game/core"
	"hz/resources/images"

	"github.com/hajimehoshi/ebiten/v2"
)

type Level struct {
	Room  *core.Sprite
	Walls []PhysicsRect
}

func NewLevel(ph *Physics) *Level {
	walls := []PhysicsRect{
		{X: 0, Y: 0, Width: 288, Height: 8},
		{X: 0, Y: 0, Width: 8, Height: 288},
		{X: 288 - 8, Y: 0, Width: 8, Height: 288},
		{X: 0, Y: 288 - 8, Width: 288, Height: 8},
	}

	// for _, wall := range walls {
	// 	ph.AddRect(wall.X, wall.Y, wall.Width, wall.Height)
	// }

	room := core.NewSprite(images.Room1)
	room.SetPosition(core.ScreenWidth/2, core.ScreenHeight/2)
	return &Level{
		Room:  room,
		Walls: walls,
	}
}

func (l *Level) Update() {
	l.Room.Update()
}

func (l *Level) Draw(screen *ebiten.Image) {
	l.Room.Draw(screen)
}
