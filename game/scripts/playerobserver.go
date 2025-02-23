package scripts

import (
	"hz/game/util"
)

type PlayerObserver struct {
	Player       *Player
	EnemyManager *EnemyManager

	playerCoords []util.Vec2
	playerFirst  int
	playerLast   int
}

func NewPlayerObserver(player *Player, em *EnemyManager) *PlayerObserver {
	const delayFrames = 12
	return &PlayerObserver{
		Player:       player,
		EnemyManager: em,
		playerCoords: make([]util.Vec2, delayFrames),
		playerFirst:  0,
		playerLast:   delayFrames - 1,
	}
}

func (po *PlayerObserver) Update() {
	po.EnemyManager.SetTarget(EnemyTarget{
		Active: true,
		Walk:   po.playerCoords[po.playerFirst],
		Radius: 10,
		Look:   util.NewVec2(po.Player.X, po.Player.Y),
	})

	po.playerFirst++
	if po.playerFirst >= len(po.playerCoords) {
		po.playerFirst = 0
	}

	po.playerCoords[po.playerLast] = util.NewVec2(po.Player.X, po.Player.Y)
	po.playerLast++
	if po.playerLast >= len(po.playerCoords) {
		po.playerLast = 0
	}
}
