package images

import (
	"bytes"
	_ "embed"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type LayerOffset struct {
	OX, OY     int
	FrameCount int
}

type Image struct {
	Image  *ebiten.Image
	Width  int
	Height int
	Layers map[string]LayerOffset
}

//go:embed ChelWalk.png
var chelWalk []byte
var Chel Image

//go:embed room1.png
var room1 []byte
var Room1 Image

//go:embed EnemyZombie.png
var enemyZombie []byte
var EnemyZombie Image

//go:embed Bullet.png
var bullet []byte
var Bullet Image

func InitImage(name string, file []byte, width, height int, layers map[string]LayerOffset) Image {
	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		log.Println(name, err)
		return Image{}
	}

	return Image{
		Image:  ebiten.NewImageFromImage(img),
		Width:  width,
		Height: height,
		Layers: layers,
	}
}

func init() {

	Chel = InitImage("chel", chelWalk, 32, 32, map[string]LayerOffset{
		"body": {OX: 0, OY: 32, FrameCount: 16},
		"legs": {OX: 0, OY: 0, FrameCount: 16},
	})

	Room1 = InitImage("room1", room1, 288, 288, nil)

	EnemyZombie = InitImage("enemyZombie", enemyZombie, 12, 24, nil)
	Bullet = InitImage("bullet", bullet, 32, 32, nil)
}
