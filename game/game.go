package game

import (
	"fmt"
	"hz/game/core"
	"hz/game/scripts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 240
	screenHeight = 240
)

type Game struct {
	Level          *scripts.Level
	Player         *scripts.Player
	EnemyManager   *scripts.EnemyManager
	PlayerObserver *scripts.PlayerObserver
	BulletMananger *scripts.BulletManager
	Weapon         *scripts.Weapon

	Physics *scripts.Physics
}

func NewGame() *Game {
	ph := scripts.NewPhysics()
	p := scripts.NewPlayer(ph)
	em := scripts.NewEnemyManager(ph)
	bm := scripts.NewBulletManager()
	po := scripts.NewPlayerObserver(p, em)
	return &Game{
		Level:          scripts.NewLevel(ph),
		Player:         p,
		EnemyManager:   em,
		PlayerObserver: po,
		BulletMananger: bm,
		Weapon:         scripts.NewWeapon(p, bm),
		Physics:        ph,
	}
}

func (g *Game) Update() error {
	g.Level.Update()
	g.Weapon.Update()
	g.Player.Update()
	g.PlayerObserver.Update()
	g.EnemyManager.Update()
	g.BulletMananger.Update()
	g.Physics.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Level.Draw(screen)
	g.Player.Draw(screen)
	g.EnemyManager.Draw(screen)
	g.BulletMananger.Draw(screen)
	// g.Physics.Draw(screen)

	x, y := ebiten.CursorPosition()
	mx, my := core.GlobalCamera.ScreenToWorld(float64(x), float64(y))
	_, _ = mx, my
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f : %0.2f %0.2f : %.2f %.2f", ebiten.ActualTPS(), g.Player.X, g.Player.Y, mx, my))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
