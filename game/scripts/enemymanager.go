package scripts

import (
	"hz/game/scripts/physics"
	"hz/game/util"
	"math/rand"
	"sync"
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

	physics *physics.Physics
}

func NewEnemyManager(p *physics.Physics) *EnemyManager {
	return &EnemyManager{
		Enemies:     []*Enemy{},
		SpawnTicker: time.NewTicker(2 * time.Second),
		physics:     p,
	}
}

func (em *EnemyManager) SetTarget(et EnemyTarget) {
	em.target = et
}

var spawnOnce sync.Once

func (em *EnemyManager) Update() {
	spawnOnce.Do(func() {
		pos := util.NewVec2(
			float64(100*(rand.Int31n(5)+1)),
			float64(100*(rand.Int31n(5)+1)),
		)
		enemy := NewEnemy(em.physics, pos)
		em.Enemies = append(em.Enemies, enemy)

	})

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
		e.target = em.target
	}

	for _, e := range em.Enemies {
		e.Update()
	}
}

func (em *EnemyManager) PhysicsUpdate(dt float64) {
	for _, e := range em.Enemies {
		if em.target.Active {
			EnemyFollowTarget(dt, e, em.target)
		}
	}
}

func (em *EnemyManager) Draw(screen *ebiten.Image) {
	for _, e := range em.Enemies {
		e.Draw(screen)
	}
}
