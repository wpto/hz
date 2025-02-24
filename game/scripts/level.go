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
		// {X: 288 - 8, Y: 0, Width: 8, Height: 288},
		// {X: 0, Y: 288 - 8, Width: 288, Height: 8},

		{X: 0, Y: 288 - 8, Width: 55, Height: 8},
		{X: 88, Y: 288 - 8, Width: 200, Height: 8},
		{X: 288 - 8, Y: 0, Width: 8, Height: 55},
		{X: 288 - 8, Y: 88, Width: 8, Height: 200},
	}

	for _, wall := range walls {
		ph.AddRect(wall.X, wall.Y, wall.Width, wall.Height)
	}

	room := core.NewSprite(images.Room1)
	room.SetPosition(288/2, 288/2)
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
