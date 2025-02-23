package scripts

import (
	"hz/game/core"
	"hz/game/util"

	"github.com/hajimehoshi/ebiten/v2"
)

type WeaponPlayer interface {
	GetPosition() (float64, float64)
	GetLookDirection() float64
}

type Weapon struct {
	cooldown float64
	timer    float64

	Player        WeaponPlayer
	BulletManager *BulletManager
}

func NewWeapon(p WeaponPlayer, BulletManager *BulletManager) *Weapon {
	return &Weapon{
		cooldown:      0.5,
		timer:         0,
		Player:        p,
		BulletManager: BulletManager,
	}
}

func (w *Weapon) Update() {
	w.timer += core.Delta

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if w.timer > w.cooldown {
			w.timer = 0
			pos := util.NewVec2(w.Player.GetPosition())
			dir := util.Vec2Right.Rotate(w.Player.GetLookDirection())
			spawnPos := pos.Add(dir.Mul(3))

			b := NewBullet(spawnPos.X, spawnPos.Y, dir.X, -dir.Y)
			w.BulletManager.AddBullet(b)
		}
	}
}
