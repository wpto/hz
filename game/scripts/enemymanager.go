package scripts

import (
	"hz/game/util"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyTarget struct {
	Active bool
	Walk   util.Vec2
	Radius float64
	Look   util.Vec2
}

type EnemyManager struct {
	Enemies     []*Enemy
	SpawnTicker *time.Ticker

	target EnemyTarget

	physics *Physics
}

func NewEnemyManager(p *Physics) *EnemyManager {
	return &EnemyManager{
		Enemies:     []*Enemy{},
		SpawnTicker: time.NewTicker(2 * time.Second),
		physics:     p,
	}
}

func (em *EnemyManager) SetTarget(et EnemyTarget) {
	em.target = et
}

func (em *EnemyManager) Update() {
	select {
	case <-em.SpawnTicker.C:
		pos := util.NewVec2(
			float64(100*(rand.Int31n(5)+1)),
			float64(100*(rand.Int31n(5)+1)),
		)
		enemy := NewEnemy(em.physics, pos)
		em.Enemies = append(em.Enemies, enemy)

	default:
	}

	for _, e := range em.Enemies {
		if em.target.Active {
			EnemyFollowTarget(e, em.target)
		}
		e.Update()
	}
}

func (em *EnemyManager) Draw(screen *ebiten.Image) {
	for _, e := range em.Enemies {
		e.Draw(screen)
	}
}
